package main

import (
	"encoding/json"
	"fmt"
	"github.com/kuangcp/gobase/pkg/ctool"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
	"time"
	"unicode"
)

func TestSelfDict(t *testing.T) {
	// https://github.com/mapull/chinese-dictionary
	bts, _ := os.ReadFile("/home/zk/Code/go/chinese-dictionary/word/word.json")
	var sli []map[string]interface{}
	json.Unmarshal(bts, &sli)

	writer, _ := ctool.NewWriter("dict/dict.log", true)
	for _, r := range sli {
		s := r["word"]
		writer.WriteLine(fmt.Sprint(s))
	}
}

func TestInit(t *testing.T) {
	tokenizer := InitTrieTokenizer("dict/zh.log")
	fmt.Println(FmtTokens(tokenizer.Tokenize("我们的世界和平，编程的速度才会得到发展")))
	fmt.Println(FmtTokens(tokenizer.Tokenize("AI和AIDS是不一样的东西,AID是瞄准")))
}

func TestFile(t *testing.T) {
	tokenizer := InitTrieTokenizer("dict/dict.log")
	tokenizer.Append("dict/code.dict")
	tokenizer.Append("dict/zk.dict.log")
	tokens := tokenizer.TokenizeFile("meeting_record.md")
	//println(FmtTokens(tokens))

	// Error
	var result = make(map[string]int)
	statisticsJudge(tokens, result, func(runes []rune) bool {
		return len(runes) == 1
	})
	printSort(result, func(s string, i int) bool {
		return i > 15
	})
	println("==============")

	// word
	result = make(map[string]int)
	statisticsJudge(tokens, result, func(runes []rune) bool {
		return len(runes) > 1
	})
	printSort(result, func(s string, i int) bool {
		return i > 10
	})
}

func TestDir(t *testing.T) {
	tokenizer := InitTrieTokenizer("dict/dict.log")

	tokenizer.Append("dict/code.dict")

	var result = make(map[string]int)
	err := filepath.WalkDir("/home/zk/Note/WorkLog/", func(path string, d fs.DirEntry, err error) error {
		if strings.HasSuffix(path, "md") {
			tokens := tokenizer.TokenizeFile(path)
			statistics(tokens, result)
		}
		return nil
	})
	if err != nil {
		log.Println(err)
		return
	}

	format := time.Now().Format(ctool.HH_MM_SS_MS)
	writer, _ := ctool.NewWriter("log/result-"+format+".log", true)
	defer writer.Close()
	for k, v := range result {
		writer.WriteLine(fmt.Sprint(v, " ", k))
	}
}

func statisticsError(tokens []string, result map[string]int) {
	for _, t := range tokens {
		runes := []rune(t)
		if unicode.Is(unicode.Scripts["Han"], runes[0]) {
			if len(runes) != 1 {
				continue
			}
			n, ok := result[t]
			if !ok {
				result[t] = 1
			} else {
				result[t] = n + 1
			}
		}
	}
}

func statistics(tokens []string, result map[string]int) {
	for _, t := range tokens {
		runes := []rune(t)
		if unicode.Is(unicode.Scripts["Han"], runes[0]) {
			n, ok := result[t]
			if !ok {
				result[t] = 1
			} else {
				result[t] = n + 1
			}
		}
	}
}

func printSort(data map[string]int, filter func(string, int) bool) {
	type KV struct {
		k string
		v int
	}
	var result []KV
	for k, v := range data {
		if filter != nil && !filter(k, v) {
			continue
		}
		result = append(result, KV{k: k, v: v})
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].v < result[j].v
	})
	for _, kv := range result {
		fmt.Println(kv.k, kv.v)
	}

}
func statisticsJudge(tokens []string, result map[string]int, filter func([]rune) bool) {
	for _, t := range tokens {
		runes := []rune(t)
		if unicode.Is(unicode.Scripts["Han"], runes[0]) {
			if filter != nil && !filter(runes) {
				continue
			}
			n, ok := result[t]
			if !ok {
				result[t] = 1
			} else {
				result[t] = n + 1
			}
		}
	}
}
