package mux

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

type contextKey string

var patternKey contextKey = "pattern"

// Param 获取param
func Param(r *http.Request, key string) string {
	val, ok := r.Context().Value(contextKey(key)).(string)
	// logDebug("Param(%q) is %t value=%q", key, ok, val)
	if !ok {
		return ""
	}
	return val
}

// Node 用来描述 Trie 结构的结点
type Node struct {
	children map[string]*Node // 孩子节点
	isLeaf   bool             // 是否叶子节点
}

// Trie 树结构路由
type Trie struct {
	root  *Node                     // 节点
	rxMap map[string]*regexp.Regexp // 保存路由所有正则
	size  int                       // 总共有多少个路由
}

// NewTrie 实例化一个Trie路由
func NewTrie() *Trie {
	return &Trie{
		root:  &Node{children: make(map[string]*Node), isLeaf: false},
		rxMap: make(map[string]*regexp.Regexp),
	}
}

// Insert 添加一个路由, segments 每段路由集合, 返回 false 则路由已存在
//
// 如："/v1/api/product/:id" => ["v1","api","product",":id"], 因为每一个路由都是以“/”开始，因此不存储“/”
//
// 动态路由规则: 其中动态参数正则表达式以 "|" 分割，如：
// "/v1/api/:id|[0-9]/:name"
func (t *Trie) Insert(segments []string) bool {
	var exists = true // 假设路由已存在
	var curNode = t.root
	for _, seg := range segments {
		node, ok := curNode.children[seg]
		if !ok {
			exists = false
			node = &Node{make(map[string]*Node), false}
			// 追加到当前节点元素孩子里
			curNode.children[seg] = node
		}

		// 解析 seg 判断是否存在正则表达式
		_, rx := t.parse(seg)
		if len(rx) > 0 {
			t.rxMap[rx] = regexp.MustCompile(rx)
		}

		curNode = node
	}
	curNode.isLeaf = true
	if !exists {
		t.size++
	}
	return !exists
}

// Match 查找 segments 路由段是否在树中，匹配成功返回true
func (t *Trie) Match(ctx context.Context, segments []string) (context.Context, bool) {
	fmt.Println("=============", segments, "===============")
	// 这里实现的方式：由访问路由的 segments，逆推出注册路由时的Path
	// 如：注册 /v1/api/id:|[0-9]/:name => ["v1","api","id:|[0-9]",":name"]
	// 访问 Api URL: /v1/api/100/tom, 由这个URL逆推得到 /v1/api/id:|[0-9]/:name，然后存在context里
	var curNode = t.root
	var url strings.Builder
	for _, seg := range segments {
		node, exists := curNode.children[seg]
		// logDebug("seg: %v", seg)
		if !exists {
			var segMatch bool // 标记该段是否匹配成功

			// 遍历 children 获取 key 查找是否有符合条件的动态param
			for key := range curNode.children {
				paramKey, regx := t.parse(key)
				// 存在动态 parma，但不存在 正则匹配
				if len(paramKey) > 0 && len(regx) == 0 {
					segMatch = true
				}

				//  存在动态正则 parma
				if len(paramKey) > 0 && len(regx) > 0 {
					rx, ok := t.rxMap[regx]
					if ok {
						segMatch = rx.MatchString(seg)
					}
				}

				// 匹配成功，设置 parma
				if segMatch {
					logDebug("set param key= %s, value= %s", paramKey, seg)

					ctx = context.WithValue(ctx, contextKey(paramKey), seg)
					url.WriteString("/" + key)
					node = curNode.children[key]
					break
				}
			}

			logDebug("seg: %s is match %t", seg, segMatch)

			// 如果匹配不成功，退出、返回
			if !segMatch {
				fmt.Println(">>>>>>>>>> end >>>>>>>>>>>>>>")
				return ctx, false
			}
		} else {
			url.WriteString("/" + seg)
		}
		curNode = node
	}
	logDebug("pattern: %s\n", url.String())
	fmt.Println(">>>>>>>>>> end >>>>>>>>>>>>>>")

	// 如果最后一个是叶子节点，则说明匹配成功，反之匹配不成功
	ctx = context.WithValue(ctx, patternKey, url.String())
	return ctx, curNode.isLeaf

}

// Size 返回路由总数
func (t *Trie) Size() int {
	return t.size
}

// parse 解析提取 segment 里面是否包含动态参数、正则表达式
//
// 如："/v1/api/:id|[0-9]/:name"
func (t *Trie) parse(segment string) (paramKey string, regx string) {
	if strings.HasPrefix(segment, ":") {
		var found bool
		paramKey, regx, found = strings.Cut(segment, "|")
		if found {
			paramKey = strings.Split(paramKey, ":")[1]
		} else {
			paramKey = strings.Split(segment, ":")[1]
		}
	}
	return
}

func (n *Node) MarshalJSON() ([]byte, error) {
	return json.MarshalIndent(map[string]interface{}{
		"children": n.children,
		"isLeaf":   n.isLeaf,
	}, "", "  ")
}

func (t *Trie) MarshalJSON() ([]byte, error) {
	return json.MarshalIndent(map[string]interface{}{
		"root":  t.root,
		"rxMap": t.rxMap,
		"size":  t.size,
	}, "", "  ")
}
