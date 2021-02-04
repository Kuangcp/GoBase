package main

import (
	"fmt"
	"testing"
)

func TestReplace(t *testing.T) {
	title := "pre（）测试 【使用】"

	result := normalizeForTitle(title)
	fmt.Println(result)

	if result != "pre测试-使用" {
		t.Fail()
	}
}
