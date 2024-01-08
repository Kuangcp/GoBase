package main

import (
	"fmt"
	"testing"
)

func TestInit(t *testing.T) {
	tokenizer := InitTrieTokenizer("zh.log")
	result := tokenizer.Tokenize("我们的世界和平，编程的速度才会得到发展")
	fmt.Println(result)
	fmt.Println(FmtTokens(result))
}
