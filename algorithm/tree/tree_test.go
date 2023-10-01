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
