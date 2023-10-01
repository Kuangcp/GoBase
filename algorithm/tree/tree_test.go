package main

import (
	"fmt"
	"testing"
)

func buildSimpleTree(l int) *Tree[int] {
	var data []int
	for i := 0; i < l; i++ {
		data = append(data, i)
	}

	return ArrayToTree[int](data)
}

func TestBuildTree(t *testing.T) {
	tree := buildSimpleTree(121)
	fmt.Println(tree)
}

func TestDfsPre(t *testing.T) {
	tree := buildSimpleTree(14)
	DfsPre(tree, PrintNode)
}

func TestDfsIn(t *testing.T) {
	tree := buildSimpleTree(14)
	DfsIn(tree, PrintNode)
}

func TestDfsPost(t *testing.T) {
	tree := buildSimpleTree(14)
	DfsPost(tree, PrintNode)
}

func TestBfs(t *testing.T) {
	tree := buildSimpleTree(14)
	Bfs(tree, PrintNode)
}
