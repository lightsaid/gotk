package gox

import (
	"context"
	"regexp"
	"strings"
)

var trieNodeKey contextKey = "trieNodeKey"

type contextKey string

func ContextKey(key string) contextKey {
	return contextKey(key)
}

// RouteTrie 路由接口
type RouteTrie interface {
	// 注册路由
	Register(parts []string) bool
	// 匹配路由
	Match(ctx context.Context, parts []string) (context.Context, bool)
	// 路由总数
	Size() int
}

// routeNode 定义一个Trie的路由结构，用来描述和存储路由节点
type routeNode struct {
	children map[string]*routeNode
	isLeaf   bool
}

// routeTrie 一个路由树
type routeTrie struct {
	root  *routeNode
	rxMap map[string]*regexp.Regexp
	size  int
}

// NewRouteTrie 实例化一个 RouteTrie
func NewRouteTrie() RouteTrie {
	return &routeTrie{
		root:  &routeNode{children: make(map[string]*routeNode), isLeaf: false},
		rxMap: make(map[string]*regexp.Regexp),
		size:  0,
	}
}

// Register 插入（注册路由）一个字符串组（path的一每段）,如果已存在返回false，反之true
func (t *routeTrie) Register(parts []string) bool {
	exists := true
	curNode := t.root
	for _, part := range parts {
		node, ok := curNode.children[part]
		// fmt.Println(">>>> ok: ", ok)
		if !ok {
			exists = false
			node = &routeNode{make(map[string]*routeNode), false}
			curNode.children[part] = node
		}

		// 存在正则表达式
		if _, _, rx, found := t.parse(part); found {
			t.addRX(rx)
		}

		curNode = node
	}
	curNode.isLeaf = true
	if !exists {
		t.size++
	}
	return !exists
}

// Match 查找 parts 字符串组是否存在，如果有正则 则匹配正则
func (t *routeTrie) Match(ctx context.Context, parts []string) (context.Context, bool) {
	curNode := t.root
	var matchKey strings.Builder
	for _, part := range parts {
		node, ok := curNode.children[part]
		if ok {
			matchKey.WriteString("/" + part)
		} else {
			// flag 标记是否匹配成功
			var flag bool
			// 遍历 curNode.children 获取 key，检查是否有正则表达式，
			// 有则正则匹配（匹配成功true，反之flase），无则直接返回false
			for key := range curNode.children {
				dyRoute, paramKey, rx, found := t.parse(key)
				// fmt.Println(curNode.children)
				// fmt.Printf(">> Match: part=%s key=%s, dyRoute=%t, paramKey=%s, rx=%s, found=%t\n", part, key, dyRoute, paramKey, rx, found)
				// 动态路由，但是没有正则表达式; part 匹配成功
				if dyRoute && !found {
					flag = true
				}

				// 动态路由，有正则表达式
				if dyRoute && found {
					rxp, isTrue := t.rxMap[rx]
					if isTrue {
						flag = rxp.MatchString(part)
						// fmt.Println(">>> isTrue: ", rx, "  ", flag)
					}
				}

				// 匹配part成功
				if flag {
					// fmt.Println("param: ", paramKey, part)
					// 存储 param 参数
					ctx = context.WithValue(ctx, contextKey(paramKey), part)
					// 设置 node
					node = curNode.children[key]
					matchKey.WriteString("/" + key)
					break
				}
			}
			// 匹配不成功
			if !flag {
				return ctx, flag
			}
		}
		curNode = node
	}
	ctx = context.WithValue(ctx, trieNodeKey, matchKey.String())
	return ctx, curNode.isLeaf
}

// Size 返回路由总数
func (t *routeTrie) Size() int {
	return t.size
}

// parse 解析每段 path
func (t *routeTrie) parse(orgPart string) (dyRoute bool, paramKey string, rx string, found bool) {
	if strings.HasPrefix(orgPart, ":") {
		dyRoute = true
		paramKey, rx, found = strings.Cut(orgPart, "|")
		if found {
			paramKey = strings.Split(paramKey, ":")[1]
		} else {
			paramKey = strings.Split(orgPart, ":")[1]
		}
	}
	return
}

// addRX 添加正则匹配
func (t *routeTrie) addRX(rx string) {
	t.rxMap[rx] = regexp.MustCompile(rx)
}
