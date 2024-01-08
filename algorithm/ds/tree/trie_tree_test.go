package main

import (
	"fmt"
	"log"
	"testing"
)

func TestCnIndex(t *testing.T) {
	trie := NewTrie()
	trie.Inserts("我", "我们", "世界")

	fmt.Println(trie)
	fmt.Println(trie.SearchKey("我的"))
	fmt.Println(trie.SearchKey("我们"))

}

func TestMerge(t *testing.T) {
	trie := NewTrie()
	trie.Inserts("我", "我们", "常量")
	trie2 := NewTrie()
	trie2.Inserts("世界", "我的", "我们的", "常量值")
	trie.Merge(trie2)

	log.Println(trie.SearchKey("我们"))
	log.Println(trie.SearchKey("我的"))
	log.Println(trie.SearchKey("我们的"))
	log.Println(trie.SearchKey("世"))
	log.Println(trie.SearchKey("世界"))
	log.Println(trie.SearchKey("常量值"))
}
