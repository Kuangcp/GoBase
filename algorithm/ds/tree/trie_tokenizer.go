package main

import (
	"github.com/kuangcp/gobase/pkg/ctool"
	"strings"
)

type (
	TrieTokenizer struct {
		trie *TrieTree
	}
)

func InitTrieTokenizer(path string) *TrieTokenizer {
	trie := NewTrie()

	lines := ctool.ReadStrLinesNoFilter(path)
	for _, r := range lines {
		trie.Insert(r)
	}

	return &TrieTokenizer{
		trie: trie,
	}
}
func (t *TrieTokenizer) mine(trie *TrieTree, sentence []rune, cursor int) int {
	if cursor <= len(sentence)-1 {
		cur := sentence[cursor]
		if trie.Match(cur) {
			cursor = t.mine(trie.Childes[cur], sentence, cursor+1)
		}
	}
	return cursor
}

func (t *TrieTokenizer) Tokenize(sentence string) []string {
	if len(sentence) <= 0 {
		return nil
	}
	var tokens []string
	runes := []rune(sentence)
	chars := len(runes)
	for chars != 0 {
		idx := t.mine(t.trie, runes, 0)
		if idx == 0 {
			//fmt.Println(idx, string(runes[idx]))
			tokens = append(tokens, string(runes[0]))
			runes = runes[1:]
			chars = len(runes)
		} else {
			//fmt.Println(idx, string(runes[idx-1]))
			tokens = append(tokens, string(runes[0:idx]))
			if idx == chars {
				return tokens
			}
			runes = runes[idx:]
			chars = len(runes)
		}
	}
	return tokens
}

func FmtTokens(tokens []string) string {
	if len(tokens) == 0 {
		return ""
	}

	return strings.Join(tokens, "_")
}
