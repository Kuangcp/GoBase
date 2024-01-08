package main

import (
	"encoding/json"
	"fmt"
	"github.com/kuangcp/gobase/pkg/ctool"
	"io/fs"
	"log"
	"os"
	"path/filepath"
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

	writer, _ := ctool.NewWriter("dict.log", true)
	for _, r := range sli {
		s := r["word"]
		writer.WriteLine(fmt.Sprint(s))
	}
}

func TestInit(t *testing.T) {
	tokenizer := InitTrieTokenizer("zh.log")
	fmt.Println(FmtTokens(tokenizer.Tokenize("我们的世界和平，编程的速度才会得到发展")))
	fmt.Println(FmtTokens(tokenizer.Tokenize("AI和AIDS是不一样的东西,AID是瞄准")))
}

func TestFile(t *testing.T) {
	tokenizer := InitTrieTokenizer("dict.log")
	tokens := tokenizer.TokenizeFile("/home/zk/Note/Note/Skills/Vcs/GitBase.md")
	//println(FmtTokens(tokens))

	var result = make(map[string]int)
	statistics(tokens, result)

	for k, v := range result {
		if v < 60 {
			continue
		}
		fmt.Println(k, v)
	}
}

func TestDir(t *testing.T) {
	tokenizer := InitTrieTokenizer("dict.log")

	tokenizer.Append("code.dict")

	var result = make(map[string]int)
	err := filepath.WalkDir("/home/zk/Note/Note/", func(path string, d fs.DirEntry, err error) error {
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
	writer, _ := ctool.NewWriter("result-"+format+".log", true)
	defer writer.Close()
	for k, v := range result {
		writer.WriteLine(fmt.Sprint(v, " ", k))
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
