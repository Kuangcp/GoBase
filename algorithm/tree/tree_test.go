package main

import (
	"fmt"
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

func buildSimpleTree(l int) *Tree[int] {
	var data []int
	for i := 0; i < l; i++ {
		data = append(data, i)
	}

	return ArrayToTree[int](data)
}

func buildSlotTree() *Tree[*User] {
	// 注意构造树时不要出现无父节点
	return ArrayToTree([]*User{{Id: 0, Name: "a"}, nil, {Id: 2, Name: "b"},
		nil, nil, {Id: 3, Name: "c"}, nil,
		nil, nil, nil, nil, {Id: 4, Name: "d"}, {Id: 5, Name: "e"}})
}

func TestBuildTree(t *testing.T) {
	tree := buildSimpleTree(121)
	fmt.Println(tree)
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
	println(isDiff(a, b))
}

func isDiff[T comparable](a, b *Tree[T]) bool {
	if a == nil && b == nil {
		return false
	}
	if a == nil || b == nil {
		return true
	}

	if isDiffVal(a, b) || isDiffVal(a.Left, b.Left) || isDiffVal(a.Right, b.Right) {
		return true
	}

	if isDiff(a.Left, b.Left) {
		return true
	}
	if isDiff(a.Right, b.Right) {
		return true
	}
	return false
}

func isDiffVal[T comparable](a, b *Tree[T]) bool {
	if a == nil && b == nil {
		return false
	}
	if a == nil || b == nil {
		return true
	}

	return a.Data != b.Data
}

// 镜像二叉树
func TestInvertTree(t *testing.T) {

}
