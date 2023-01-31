package main

import "fmt"

const AlphabetSize = 26

// Node
type Node struct {
	childrens [AlphabetSize]*Node
	isLeaf    bool
}

// Trie
type Trie struct {
	root *Node
}

func InitTrie() *Trie {
	return &Trie{root: &Node{}}
}

// Insert
func (t *Trie) Insert(w string) {
	wl := len(w)
	currentNode := t.root
	for i := 0; i < wl; i++ {
		// 总共26个字母且都是小写的, 因此  w[i] - 'a' 得到每个字母存放下标
		charIndex := w[i] - 'a'
		if currentNode.childrens[charIndex] == nil {
			currentNode.childrens[charIndex] = &Node{}
		}
		currentNode = currentNode.childrens[charIndex]
	}
	currentNode.isLeaf = true
}

// Search
func (t *Trie) Search(w string) bool {
	wl := len(w)
	currentNode := t.root
	for i := 0; i < wl; i++ {
		// 总共26个字母且都是小写的, 因此  w[i] - 'a' 得到每个字母存放下标
		charIndex := w[i] - 'a'
		if currentNode.childrens[charIndex] == nil {
			return false
		}
		currentNode = currentNode.childrens[charIndex]
	}
	return currentNode.isLeaf
}

func main2() {
	// testTrie := InitTrie()
	// fmt.Println(testTrie.root)
	// fmt.Println('B' - 'a')
	// fmt.Println(rune('a'), 'b', 'B')

	trie := InitTrie()
	trie.Insert("apple")
	trie.Insert("huawei")
	trie.Insert("oppo")
	trie.Insert("haha")
	trie.Insert("art")
	trie.Insert("appo")
	trie.Insert("applp")

	fmt.Println(trie.Search("apple"))
	fmt.Println(trie.Search("art"))
	fmt.Println(trie.Search("haha"))
	fmt.Println(trie.Search("appo"))
	fmt.Println(trie.Search("applp"))

	fmt.Println(trie.root)

}
