package main

import (
	"fmt"
	"testing"
)

func TestInitBsTree(t *testing.T) {
	var tree *Tree[int] = nil
	for i := 0; i < 7; i++ {
		tree = BsTreeInsert(tree, i)
	}

	fmt.Print(tree)
}
