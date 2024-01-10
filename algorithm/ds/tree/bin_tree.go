package main

import (
	"fmt"
	"github.com/kuangcp/gobase/pkg/ctool"
	"github.com/kuangcp/gobase/pkg/ctool/algo"
)

type (
	// BinTree 二叉树
	BinTree[T any] struct {
		Parent *BinTree[T]
		Left   *BinTree[T]
		Right  *BinTree[T]
		Data   T
	}
)

func (t *BinTree[T]) GetLeft() algo.IBinTree {
	return t.Left
}

func (t *BinTree[T]) GetRight() algo.IBinTree {
	return t.Right
}

func (t *BinTree[T]) ToString() string {
	return fmt.Sprint(t.Data)
}

func ArrayToTree[T any](data []T) *BinTree[T] {
	length := len(data)
	if length == 0 {
		return nil
	}

	cache := make(map[int]*BinTree[T])
	var tree = &BinTree[T]{Data: data[0]}
	cache[0] = tree
	for i := 0; i < length/2; i++ {
		// left
		parent := cache[i]
		leftIdx := 2*i + 1
		if leftIdx < length {
			val := data[leftIdx]
			if !ctool.IsNil(val) {
				node := &BinTree[T]{Data: val}
				node.Parent = parent
				cache[leftIdx] = node
				parent.Left = node
			}
		}

		// right
		rightIdx := 2*i + 2
		if rightIdx < length {
			val := data[rightIdx]
			if !ctool.IsNil(val) {
				node := &BinTree[T]{Data: val}
				node.Parent = parent
				cache[rightIdx] = node
				parent.Right = node
			}
		}

	}
	return tree
}

// DfsPre 前序遍历
func DfsPre[T any](t *BinTree[T], handler func(node *BinTree[T])) {
	if t == nil || handler == nil {
		return
	}
	handler(t)
	DfsPre(t.Left, handler)
	DfsPre(t.Right, handler)
}

// DfsIn 中序遍历
func DfsIn[T any](t *BinTree[T], handler func(node *BinTree[T])) {
	if t == nil || handler == nil {
		return
	}
	DfsIn(t.Left, handler)
	handler(t)
	DfsIn(t.Right, handler)
}

// DfsPost 后序遍历
func DfsPost[T any](t *BinTree[T], handler func(node *BinTree[T])) {
	if t == nil || handler == nil {
		return
	}
	DfsPost(t.Left, handler)
	DfsPost(t.Right, handler)
	handler(t)
}

// Bfs 广度优先遍历 层次遍历
func Bfs[T any](t *BinTree[T], handler func(node *BinTree[T])) {
	if t == nil || handler == nil {
		return
	}

	handler(t)

	var cur []*BinTree[T]
	cur = appendIfAbsent(cur, t.Left)
	cur = appendIfAbsent(cur, t.Right)
	for {
		var nextLayer []*BinTree[T]
		for _, node := range cur {
			handler(node)
			nextLayer = appendIfAbsent(nextLayer, node.Left)
			nextLayer = appendIfAbsent(nextLayer, node.Right)
		}

		if len(nextLayer) == 0 {
			break
		}

		cur = nextLayer
	}
}

func appendIfAbsent[T any](layer []*BinTree[T], node *BinTree[T]) []*BinTree[T] {
	if node == nil {
		return layer
	}
	return append(layer, node)
}

func PrintNode[T any](node *BinTree[T]) {
	fmt.Print(node.Data, " ")
}
