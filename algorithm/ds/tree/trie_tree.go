package main

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

func (t *TrieTree) Match(r rune) bool {
	if len(t.Childes) == 0 {
		return false
	}
	_, ok := t.Childes[r]
	return ok
}

// TODO Merge https://oi-wiki.org/string/trie/
