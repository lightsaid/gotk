package main

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

/*
使用 Trie 结构实现路由动态路由注册和匹配

// 定义路由规则：
// 假设需要匹配以下路由

// ---------------------------------------------------------------------
/v1/api/users
/v2/api/users
/v1/api/users/:id
/v2/api/users/:id/mine
/v1/api/users/sys/:id|[0-9]+
/v2/api/users/sys/:id|[0-9]+/profile
/v1/api/users/sys/:status|^[A-Z]+$
/v2/api/users/sys/:status|^[A-Z]+$
/v1/api/users/:gender/:age|[0-9]{1,3}$
/v2/api/users/:gender/:age|[0-9]{1,3}$/:name

/api/products  		   -> Get 获取列表
/api/products/:id      -> Get 获取单个
/api/products/:id      -> Post 添加
/api/products/:id      -> Put 更新
/api/products/:id      -> Delete 删除

/api/category
/api/category/list
/api/category/get
/api/category/post
/api/category/put
/api/category/del

/api/products/:category/:id

/static/...   -> 通配符，匹配下面剩余部分

// ---------------------------------------------------------------------

*/

// RouteNode 定义一个Trie的路由结构，用来描述和存储路由节点
type RouteNode struct {
	Children map[string]*RouteNode
	IsLeaf   bool
}

// RouteTrie 一个路由树
type RouteTrie struct {
	Root  *RouteNode
	RxMap map[string]*regexp.Regexp
	Size  int
}

// NewRouteTrie 实例化一个RouteTrie
func NewRouteTrie() *RouteTrie {
	return &RouteTrie{
		Root:  &RouteNode{make(map[string]*RouteNode), false},
		RxMap: make(map[string]*regexp.Regexp),
	}
}

// Insert 插入一个字符串组
func (t *RouteTrie) Insert(parts []string) bool {
	exists := true
	curNode := t.Root
	for _, part := range parts {
		node, ok := curNode.Children[part]
		if !ok {
			exists = false
			node = &RouteNode{make(map[string]*RouteNode), false}
			curNode.Children[part] = node
		}

		// 存在正则表达式
		if _, _, rx, found := t.Parse(part); found {
			t.addRX(rx)
		}

		curNode = node
	}
	curNode.IsLeaf = true
	if !exists {
		t.Size++
	}
	return exists
}

// Search 查找 parts 字符串组是否存在
func (t *RouteTrie) Search(parts []string) bool {
	curNode := t.Root
	for _, part := range parts {
		node, ok := curNode.Children[part]
		if !ok {
			return false
		}
		curNode = node
	}
	return curNode.IsLeaf
}

// Parse
func (t *RouteTrie) Parse(part string) (dyRoute bool, key string, rx string, found bool) {
	if strings.HasPrefix(part, ":") {
		dyRoute = true
		key, rx, found = strings.Cut(part, "|")
		if found {
			key = strings.Split(key, ":")[1]
		} else {
			key = strings.Split(part, ":")[1]
		}
	}
	return
}

func (t *RouteTrie) addRX(rx string) {
	t.RxMap[rx] = regexp.MustCompile(rx)
}

// Match 查找 parts 字符串组是否存在，如果有正则 则匹配正则
func (t *RouteTrie) Match(parts []string) bool {
	curNode := t.Root
	for _, part := range parts {
		node, ok := curNode.Children[part]
		if !ok {
			var flag bool
			// 遍历 curNode.Children 获取 key，检查是否有正则表达式，
			// 有则正则匹配（匹配成功true，反之flase），无则直接返回false
			for key := range curNode.Children {
				fmt.Println(">> key: ", key)
				dyRoute, _, rx, found := t.Parse(key)
				fmt.Println(">> parse: ", dyRoute, rx, found)
				// 动态路由，但是没有正则表达式; part 匹配成功
				if dyRoute && !found {
					flag = true
					node = curNode.Children[key]
				}

				// 动态路由，有正则表达式
				if dyRoute && found {
					flag = t.RxMap[rx].MatchString(part)
					if flag {
						node = curNode.Children[key]
					}
				}
			}

			if !flag {
				return flag
			}
		}
		curNode = node
	}
	return curNode.IsLeaf
}

// Size 存在多少个分支（有多少条路径）
// func (t RouteTrie) Size() int {
// 	return t.size
// }

func main() {
	var paths = []string{
		"/v1/api/users",
		"/v2/api/users/sys/:id|[0-9]+/profile",
		"/v2/api/users/:gender/:age|[0-9]{1,3}$/:name",
		// "/api/category",
		// "/api/category",
		// "/api/category/list",
	}

	var trie = NewRouteTrie()

	for _, path := range paths {
		t := time.Now()
		parts := strings.Split(path, "/")[1:]
		trie.Insert(parts)
		fmt.Println(time.Since(t))
	}

	// s, _ := json.MarshalIndent(trie, "", " ")
	// fmt.Println(string(s))

	fmt.Println("==========================")

	// 查找
	var matchPaths = []string{
		"/v1/api/users",
		"/v2/api/users/sys/100/profile",
		"/v2/api/users/sys/aaa/profile",
		"/v2/api/users/male/22/mario",
		"/v2/api/users/male/aabb/mario",

		// "/api/category",
		// "/api/category",
		// "/api/category/list",
	}
	for _, path := range matchPaths {
		t := time.Now()
		parts := strings.Split(path, "/")[1:]
		exists := trie.Match(parts)
		fmt.Println(time.Since(t))
		fmt.Printf("[%s]=%t\n", path, exists)
	}
}
