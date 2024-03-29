package main

import (
	"fmt"
	"github.com/kuangcp/gobase/pkg/ctool/algo"
	"testing"
)

type (
	User struct {
		Id   int
		Name string
	}
)

func (u User) String() string {
	return fmt.Sprint(u.Id, u.Name)
}

func buildSimpleTree(l int) *BinTree[int] {
	var data []int
	for i := 0; i < l; i++ {
		data = append(data, i)
	}

	return ArrayToTree[int](data)
}

func buildSlotTree() *BinTree[*User] {
	// 注意构造树时不要出现无父节点
	return ArrayToTree([]*User{{Id: 0, Name: "a"}, nil, {Id: 2, Name: "b"},
		nil, nil, {Id: 3, Name: "c"}, nil,
		nil, nil, nil, nil, {Id: 4, Name: "d"}, {Id: 5, Name: "e"}})
}

func TestMindMap(t *testing.T) {
	tree := buildSimpleTree(20)
	algo.WriteBiMindMap(tree, "user.pu")
}

func TestBuildTree(t *testing.T) {
	tree := buildSimpleTree(121)
	fmt.Println(tree)
}

func BenchmarkDfsPre(b *testing.B) {
	tree := buildSimpleTree(10000)
	for i := 0; i < b.N; i++ {
		sum := 0
		DfsPre(tree, func(node *BinTree[int]) {
			sum += node.Data
		})
		fmt.Print(sum, " ")
	}
	fmt.Println()
}

func BenchmarkBfs(b *testing.B) {
	tree := buildSimpleTree(10000)
	for i := 0; i < b.N; i++ {
		sum := 0
		Bfs(tree, func(node *BinTree[int]) {
			sum += node.Data
		})
		fmt.Print(sum, " ")
	}
	fmt.Println()
}

func TestDfsPreBench(t *testing.T) {
	tree := buildSimpleTree(14)
	DfsPre(tree, PrintNode[int])

	slotTree := buildSlotTree()
	DfsPre(slotTree, PrintNode[*User])
}

func TestDfsPre(t *testing.T) {
	tree := buildSimpleTree(14)
	DfsPre(tree, PrintNode[int])

	slotTree := buildSlotTree()
	DfsPre(slotTree, PrintNode[*User])
}

func TestDfsIn(t *testing.T) {
	tree := buildSimpleTree(14)
	DfsIn(tree, PrintNode[int])
}

func TestDfsPost(t *testing.T) {
	tree := buildSimpleTree(14)
	DfsPost(tree, PrintNode[int])

	slotTree := buildSlotTree()
	DfsPost(slotTree, PrintNode[*User])
}

func TestBfs(t *testing.T) {
	tree := buildSimpleTree(14)
	Bfs(tree, PrintNode[int])

	slotTree := buildSlotTree()
	Bfs(slotTree, PrintNode[*User])
}

// 比较两个树是否相同
func TestSameTree(t *testing.T) {
	a := ArrayToTree[int]([]int{1, 3, 5, 2, 0, 5})
	b := ArrayToTree[int]([]int{1, 3, 5, 2, 0, 5})
	println(isDiffVal(a, b))
}

func isDiffVal[T comparable](a, b *BinTree[T]) bool {
	if a == nil && b == nil {
		return false
	}
	if a == nil || b == nil {
		return true
	}

	if a.Data != b.Data {
		return true
	} else {
		return isDiffVal(a.Left, b.Left) || isDiffVal(a.Right, b.Right)
	}
}

func DfsPreInvert[T any](t *BinTree[T]) {
	if t == nil {
		return
	}

	tmp := t.Left
	t.Left = t.Right
	t.Right = tmp

	DfsPreInvert(t.Left)
	DfsPreInvert(t.Right)
}

// 镜像二叉树
func TestInvertTree(t *testing.T) {
	// 前序遍历, 交换子节点
	tree := buildSimpleTree(7)

	DfsPre(tree, PrintNode[int])
	fmt.Println()
	DfsPreInvert(tree)
	DfsPre(tree, PrintNode[int])
}
