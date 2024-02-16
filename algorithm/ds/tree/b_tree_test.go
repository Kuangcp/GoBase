package main

import (
	"github.com/kuangcp/gobase/pkg/ctool"
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

// 模拟数据库中存储行数据的索引
func TestInitRecordBayerTree(t *testing.T) {
	type User struct {
		id     int
		name   string
		addr   string
		job    string
		school string
	}

	tree := CreateBayerTree[int, User](5)

	tree.Insert(3, User{id: 3, name: ctool.RandomAlpha(4)})
	tree.Insert(3, User{id: 3, name: ctool.RandomAlpha(4)})

	println(algo.PrintNMindMap(tree.root))
}
