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
	tree.Inserts(4, 9, 7, 22, 24, 25, 40, 10, 18, 8, 5, 30, 10, 11, 15, 99, 6, 24, 23, 1, 27, 55, 3, 12, 13, 21, 14, 16, 17, 2, 26, 20)

	fmt.Println(algo.PrintBiMindMap(tree.Root))
	// 调试方式 安装plantuml插件后分屏到右侧，修改上诉序列可实现实时查看树调整情况
	// 由于是脑图渲染，左右不分，需要观察确认下才行
	algo.WriteBiMindMap(tree.Root, "init-balance.puml")
}

func TestBstBalanceDfsPre(t *testing.T) {

}
