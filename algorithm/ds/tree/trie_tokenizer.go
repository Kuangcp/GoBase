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
		trie.Insert(strings.TrimSpace(r))
	}

	return &TrieTokenizer{
		trie: trie,
	}
}
func (t *TrieTokenizer) Append(path string) {
	lines := ctool.ReadStrLinesNoFilter(path)
	for _, r := range lines {
		t.trie.Insert(strings.TrimSpace(r))
	}
}

func (t *TrieTokenizer) mine(trie *TrieTree, sentence []rune, cursor int) int {
	result := cursor
	if cursor <= len(sentence)-1 {
		cur := sentence[cursor]
		if trie.Match(cur) {
			result = t.mine(trie.Childes[cur], sentence, cursor+1)
		}
	}
	return result
}

func (t *TrieTokenizer) TokenizeFile(file string) []string {
	lines := ctool.ReadStrLinesNoFilter(file)
	all := ""
	for _, r := range lines {
		all += strings.TrimSpace(r)
	}
	return t.Tokenize(all)
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
		if idx != 0 {
			// 贪心匹配到了非结束字符，需要回退
			token := string(runes[0:idx])
			node := t.trie.Search(token)
			//fmt.Println(token, node)
			for !node.End {
				if idx == 0 {
					break
				}
				idx--
				token := string(runes[0:idx])
				node = t.trie.Search(token)
			}
		}

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
