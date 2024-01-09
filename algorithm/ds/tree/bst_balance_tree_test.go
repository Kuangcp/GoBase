package main

import (
	"fmt"
	"github.com/kuangcp/gobase/pkg/ctool/algo"
	"testing"
)

func TestInitBsBalance(t *testing.T) {
	tree := InitBsBalanceTree[int]()
	tree.Inserts(3, 4, 7, 3, 6, 6, 1, 0, -1, 5, 6)

	fmt.Println(algo.PrintBiMindMap(tree.Root))
	//algo.WriteBiMindMap(tree.Root, "init-balance.pu")
}

func TestBalance1(t *testing.T) {
	tree := InitBsBalanceTree[int]()
	tree.Inserts(36, 24, 10, 8, 5, 10, 11, 15, 6, 24, 23, 27)

	fmt.Println(algo.PrintBiMindMap(tree.Root))
	algo.WriteBiMindMap(tree.Root, "init-balance.puml")
}
