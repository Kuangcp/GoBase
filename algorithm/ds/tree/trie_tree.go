package main

import "fmt"

// https://en.wikipedia.org/wiki/Trie
// https://oi-wiki.org/string/trie/
// https://blog.csdn.net/qq_32445015/article/details/82711249

type (
	TrieTree struct {
		Childes map[rune]*TrieTree
		Data    rune
		End     bool
	}
)

func NewTrie() *TrieTree {
	return &TrieTree{Childes: make(map[rune]*TrieTree)}
}

func (t *TrieTree) Inserts(str ...string) {
	if len(str) > 0 {
		for _, s := range str {
			t.Insert(s)
		}
	}
}

// github.com/dghubble/trie
func (t *TrieTree) Insert(str string) {
	runes := []rune(str)

	var p = t
	for _, r := range runes {
		next, ok := p.Childes[r]
		if ok {
			p = next
		} else {
			next := &TrieTree{Data: r, Childes: make(map[rune]*TrieTree)}
			p.Childes[r] = next
			p = next
		}
	}
	p.End = true
}
func (t *TrieTree) SearchKey(key string) string {
	node := t.Search(key)
	if node == nil {
		return "NIL"
	}
	return fmt.Sprint(string(node.Data), " ", node.End)
}

func (t *TrieTree) Search(key string) *TrieTree {
	runes := []rune(key)
	var p = t
	for _, r := range runes {
		next, ok := p.Childes[r]
		if !ok {
			return nil
		}
		p = next
	}
	return p
}

func (t *TrieTree) IsChild(r rune) bool {
	if len(t.Childes) == 0 {
		return false
	}
	_, ok := t.Childes[r]
	return ok
}

func (t *TrieTree) Merge(tree *TrieTree) {
	if tree == nil {
		return
	}

	for _, v := range tree.Childes {
		merge(t, v)
	}
}

// https://oi-wiki.org/string/trie/
func merge(target, origin *TrieTree) {
	if !target.IsChild(origin.Data) {
		target.Childes[origin.Data] = &TrieTree{Data: origin.Data, End: origin.End, Childes: make(map[rune]*TrieTree)}
	}
	for _, v1 := range origin.Childes {
		merge(target.Childes[origin.Data], v1)
	}
}
