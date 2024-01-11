package main

import (
	"fmt"
	"testing"
)

func TestHannSvg(t *testing.T) {
	list := []any{"*", "*", "*", "/", 5, "*", 3.14, 1, 3, nil, nil, 6, 6}
	tree := build(list...)
	tree.buildTreeInfo()
	tree.info2SVG()
	tree.showSVG()

	fmt.Println(tree.Info.Data)
	fmt.Println(tree.Info.DataLevel)
}

func TestHannSvg2(t *testing.T) {
	bst := InitBsBalanceTree[int](103, 51, 52, 53, 54, 81, 76, 75, 74, 80, 100, 102, 77, 56, 39, 42, 4, 86, 9, 7, 22, 83, 24, 25, 40, 10, 18, 8, 5, 30, 87, 19, 28, 29, 10, 11, 15, 99, 6,
		24, 23, 88, 1, 27, 55, 3, 12, 13, 21, 14, 45, 48, 49, 50, 66, 81, 82, 16, 17, 2, 26, 20, 41, 5, 6)
	tree := buildByTree(bst.Root)
	tree.buildTreeInfo()
	tree.info2SVG()
	tree.showSVG()

	fmt.Println(tree.Info.Data)
	fmt.Println(tree.Info.DataLevel)
}
