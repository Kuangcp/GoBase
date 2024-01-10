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

func TestBalanceUml(t *testing.T) {
	tree := InitBsBalanceTree[int]()
	tree.Inserts(77, 56, 39, 42, 4, 86, 9, 7, 22, 83, 24, 25, 40, 10, 18, 8, 5, 30, 87, 19, 28, 29, 10, 11, 15, 99,
		6, 24, 23, 88, 1, 27, 55, 3, 12, 13, 21, 14, 45, 48, 49, 50, 66, 81, 82, 16, 17, 2, 26, 20, 41)

	fmt.Println(algo.PrintBiMindMap(tree.Root))
	// 调试方式 安装plantuml插件后分屏到右侧，修改上诉序列可实现实时查看树调整情况
	// 由于是脑图渲染，左右不分
	algo.WriteBiMindMap(tree.Root, "init-balance.puml")
}

// 对比同样的数据，BST的高度更高
func TestCompareWithBst(t *testing.T) {
	bst := InitBsTree[int]()
	balance := InitBsBalanceTree[int]()
	ns := []int{77, 56, 39, 42, 4, 86, 9, 7, 22, 83, 24, 25, 40, 10, 18, 8, 5, 30, 87, 19, 28, 29, 10, 11, 15, 99, 6,
		24, 23, 88, 1, 27, 55, 3, 12, 13, 21, 14, 45, 48, 49, 50, 66, 81, 82, 16, 17, 2, 26, 20, 41}
	bst.Inserts(ns...)
	balance.Inserts(ns...)

	algo.WriteBiMindMap(bst.Root, "img-bst1.puml")
	algo.WriteBiMindMap(balance.Root, "img-balance1.puml")
}

func TestBstBalanceDfsIn(t *testing.T) {
	tree := InitBsBalanceTree[int](1, 7, 6, 1, 2, 8, 10, 12, 99, 3, 4, 30, 20, 23)
	fmt.Println(tree.Lists())
}
