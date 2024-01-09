package main

import (
	"fmt"
	"log"
	"testing"
)

func TestInitBsTree(t *testing.T) {
	var tree = InitBsTree[int]()
	ns := []int{1, 3, 5, 2, 9, 3, 3, 5, -1, -19}
	for _, n := range ns {
		tree.Insert(n)
	}

	fmt.Print(tree)
	node := tree.Search(9)
	fmt.Println(node)
	log.Println("max", tree.Max(), "min", tree.Min())
}
