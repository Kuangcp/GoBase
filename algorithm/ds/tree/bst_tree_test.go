package main

import (
	"fmt"
	"github.com/kuangcp/gobase/pkg/ctool"
	"github.com/kuangcp/gobase/pkg/ctool/algo"
	"github.com/stretchr/testify/assert"
	"log"
	"math"
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

func TestDelete(t *testing.T) {
	var tree = InitBsTree[int]()
	tree.Inserts(1, 7, 3, 1)
	tree.Delete(1)
	fmt.Println(tree)
}

func TestRank(t *testing.T) {
	var tree = InitBsTree[int]()
	tree.Inserts(3, 4, 7, 3, 6, 6, 1, 0, -1, 5, 6)
	assert.Equal(t, 1, tree.Rank(-1))
	assert.Equal(t, 2, tree.Rank(0))
	assert.Equal(t, 3, tree.Rank(1))
	assert.Equal(t, 4, tree.Rank(3))
	assert.Equal(t, 6, tree.Rank(4))
	assert.Equal(t, 7, tree.Rank(5))
	assert.Equal(t, 8, tree.Rank(6))
	assert.Equal(t, 11, tree.Rank(7))

	assert.Equal(t, 5, algo.Height(tree.Root))
	assert.Equal(t, 4, algo.Height(tree.Root.Right))
	assert.Equal(t, 3, algo.Height(tree.Root.Right.Right))
	assert.Equal(t, 2, algo.Height(tree.Root.Right.Right.Left))
	assert.Equal(t, 3, algo.Height(tree.Root.Left))
}

func TestBstUml(t *testing.T) {
	tree := InitBsTree[int](7, 6, 1, 2, 8, 10, 12, 99, 3, 4, 30, 20, 23)
	fmt.Println(algo.PrintBiMindMap(tree.Root))
	algo.WriteBiMindMap(tree.Root, "init-bst.puml")
}

func TestBstDfsIn(t *testing.T) {
	tree := InitBsTree[int](1, 7, 6, 1, 2, 8, 10, 12, 99, 3, 4, 30, 20, 23)
	fmt.Println(tree.Lists())
}

func TestSvg(t *testing.T) {
	writer, _ := ctool.NewWriter("bst.svg", true)
	defer writer.Close()

	tree := InitBsTree[int](1, 7, 6, 1, 2, 8, 10, 12, 99, 3, 4, 30, 20, 23)
	svg := algo.BuildSvg(tree.Root)
	writer.WriteString(svg)
}

func TestLog2(t *testing.T) {
	println(math.Log2(3))
	println(int(math.Pow(2, 4)))
}
