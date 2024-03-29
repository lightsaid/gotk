package mux

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"sync"
)

var ErrConflict = errors.New("pattern conflict")

// Node 用来描述 Trie 结构的结点
type Node struct {
	children    map[string]*Node // 孩子节点
	pattern     string           // 匹配模式
	fullpattern string
	handler     http.Handler     // http handler
	middlewares []MiddlewareFunc // 中间件
	isLeaf      bool             // 是否叶子节点
}

func (n *Node) addMiddleware(mws ...MiddlewareFunc) {
	n.middlewares = append(n.middlewares, mws...)
}

// Trie 树结构路由
type Trie struct {
	root  *Node                     // 节点
	rxMap map[string]*regexp.Regexp // 保存路由所有正则
	size  int                       // 总共有多少个路由
	mutex sync.Mutex
}

// NewTrie 实例化一个Trie路由
func NewTrie() *Trie {
	trie := &Trie{
		root: &Node{
			children: make(map[string]*Node),
		},
		rxMap: make(map[string]*regexp.Regexp),
	}

	return trie
}

func (t *Trie) Insert(pattern string, handler http.Handler, method string) (*Node, error) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	var segments = []string{method}
	if pattern == "/" {
		segments = append(segments, "/")
	} else {
		segments = append(segments, strings.Split(pattern, "/")[1:]...)
	}

	var exists = true // 假设路由已存在
	var curNode = t.root
	var height = len(segments) - 1

	for index, seg := range segments {
		segment := seg
		// 解析 seg 判断是否动态路由, 将类似 :id、:name、:title 动态参数 node key 设置为 :*
		paramKey, rx := t.parse(seg)
		if paramKey != "" && rx != "" {
			t.rxMap[rx] = regexp.MustCompile(rx)
		}
		// NOTE: filepath 是静态资源服务的参数，因此它没有正则表达式
		if paramKey != "" && rx == "" && paramKey != "filepath" {
			seg = ":*"
		}
		matchNode, ok := curNode.children[seg]
		if ok {
			if matchNode.isLeaf && index == height {
				return nil, fmt.Errorf("%w: %s %q", ErrConflict, method, pattern)
			}
		} else {
			exists = false
			matchNode = &Node{children: make(map[string]*Node), pattern: segment, isLeaf: false}
			// 追加到当前节点元素孩子里
			curNode.children[seg] = matchNode
		}
		curNode = matchNode
	}

	if !exists {
		t.size++
	}

	curNode.isLeaf = true
	curNode.fullpattern = pattern
	curNode.handler = handler

	return curNode, nil
}

// Match 查找 segments 路由段是否在树中，匹配成功返回true
func (t *Trie) Match(r *http.Request, notAllowedMethods ...string) (*Node, bool) {
	path := r.URL.Path
	method := r.Method
	if len(notAllowedMethods) > 0 {
		method = notAllowedMethods[0]
	}
	var segments = []string{method}
	if path == "/" {
		segments = append(segments, "/")
	} else {
		segments = append(segments, strings.Split(path, "/")[1:]...)
	}
	// fmt.Println("segments >>> ", len(segments), segments)
	var curNode = t.root
	for _, seg := range segments {
		matchNode, exists := curNode.children[seg]
		// fmt.Println("exists >>> ", exists, seg)
		if !exists {
			var segMatch bool // 标记该段是否匹配成功
			// 遍历 children 获取 key 查找是否有符合条件的动态param
			for key, childNode := range curNode.children {
				paramKey, regx := t.parse(childNode.pattern)
				// fmt.Println("parse >>> ", paramKey, regx, childNode.pattern)
				// 存在动态 parma，但不存在 正则匹配
				if len(paramKey) > 0 && len(regx) == 0 {
					segMatch = true
					// NOTE: 判断是否静态资源
					if paramKey == "filepath" && childNode.isLeaf {
						curNode = childNode
						goto FilePath
					}
				}

				//  存在动态正则 parma
				if len(paramKey) > 0 && len(regx) > 0 {
					rx, ok := t.rxMap[regx]
					if ok {
						segMatch = rx.MatchString(seg)
						// log.Println("segMatch >>>> ", segMatch, seg, regx)
					}
				}

				// 匹配成功，设置 parma
				if segMatch {
					// 此处没必要设置，有模糊(:id、/:name ... => :*)参数，设置了也不准,交由ServeMux匹配
					// ctx := context.WithValue(r.Context(), contextKey(paramKey), seg)
					// r = r.WithContext(ctx)
					// logDebug("%s=%s", paramKey, seg)
					matchNode = curNode.children[key]
					break
				}
			}

			// 如果匹配不成功，退出、返回
			if !segMatch {
				return matchNode, false
			}
		}
		curNode = matchNode
	}

FilePath:
	return curNode, curNode.isLeaf
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
		"children":    n.children,
		"isLeaf":      n.isLeaf,
		"pattern":     n.pattern,
		"fullpattern": n.fullpattern,
	}, "", "  ")
}

func (t *Trie) MarshalJSON() ([]byte, error) {
	return json.MarshalIndent(map[string]interface{}{
		"root":  t.root,
		"rxMap": t.rxMap,
		"size":  t.size,
	}, "", "  ")
}
