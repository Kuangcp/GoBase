package main

import (
	"fmt"
	"os"
	"testing"
)

func init() {
	os.MkdirAll("img", 0777)
}

func TestHannSvg(t *testing.T) {
	list := []any{"*", "*", "*", "/", 5, "*", 3.14, 1, 3, nil, nil, 6, 6}
	tree := build(list...)
	tree.buildTreeInfo()
	tree.info2SVG()
	tree.showSVG("img/expression")

	fmt.Println(tree.Info.Data)
	fmt.Println(tree.Info.DataLevel)
}

func TestHannSvg2(t *testing.T) {
	var ds = []int{1, 0, 5, 8, 3, 7, 2, 9, 11, 5, 3}
	balance := InitAvlTree[int](ds...)
	tree := buildByTree(balance.Root)
	tree.buildTreeInfo()
	tree.info2SVG()
	tree.showSVG("img/balance")

	bst := InitBsTree[int](ds...)
	tree = buildByTree(bst.Root)
	tree.buildTreeInfo()
	tree.info2SVG()
	tree.showSVG("img/bst")
}
