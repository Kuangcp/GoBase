package main

import (
	"fmt"
	"testing"
)

func TestReplace(t *testing.T) {
	title := "pre（）测试 【使用】"

	prepareContext()
	result := normalizeForTitle(title)
	fmt.Println(result)

	if result != "pre测试-使用" {
		t.Fail()
	}
}

func TestRefreshTitle(t *testing.T) {
	prepareContext()
	//RefreshTagAndCatalog("test.md")
}
