package main

import (
	"github.com/kuangcp/gobase/pkg/ctool/algo"
	"math"
	"testing"
)

func TestCeil(t *testing.T) {
	println(math.Floor(7 / 3))
	println(math.Ceil(7 / 3))
}

func TestInitBayerTree(t *testing.T) {
	tree := CreateBayerTree[int, string](5)

	tree.Insert(3, "3key")
	// 验证左右等值三种情况 DONE
	//tree.Insert(5, "5key")
	//tree.Insert(2, "2key")
	tree.Insert(3, "3-2key")

	println(algo.PrintNMindMap(tree.root))
}
