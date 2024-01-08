package main

import (
	"fmt"
	"testing"
)

func TestCnIndex(t *testing.T) {
	trie := NewTrie()
	trie.Inserts("我", "我们", "世界")

	fmt.Println(trie)
	fmt.Println(trie.Search("我的"))
	fmt.Println(trie.Search("我们"))

}
