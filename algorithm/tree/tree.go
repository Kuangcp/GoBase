package main

import (
	"fmt"
	"github.com/kuangcp/gobase/pkg/ctool"
)

type (
	Tree[T any] struct {
		Parent *Tree[T]
		Left   *Tree[T]
		Right  *Tree[T]
		Data   T
	}
)

func ArrayToTree[T any](data []T) *Tree[T] {
	length := len(data)
	if length == 0 {
		return nil
	}

	cache := make(map[int]*Tree[T])
	var tree = &Tree[T]{Data: data[0]}
	cache[0] = tree
	for i := 0; i < length/2; i++ {
		// left
		parent := cache[i]
		leftIdx := 2*i + 1
		if leftIdx < length {
			val := data[leftIdx]
			if !ctool.IsNil(val) {
				node := &Tree[T]{Data: val}
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
				node := &Tree[T]{Data: val}
				node.Parent = parent
				cache[rightIdx] = node
				parent.Right = node
			}
		}

	}
	return tree
}

// DfsPre 前序遍历
func DfsPre[T any](t *Tree[T], handler func(node *Tree[T])) {
	if t == nil || handler == nil {
		return
	}
	handler(t)
	DfsPre(t.Left, handler)
	DfsPre(t.Right, handler)
}

// DfsIn 中序遍历
func DfsIn[T any](t *Tree[T], handler func(node *Tree[T])) {
	if t == nil || handler == nil {
		return
	}
	DfsIn(t.Left, handler)
	handler(t)
	DfsIn(t.Right, handler)
}

// DfsPost 后序遍历
func DfsPost[T any](t *Tree[T], handler func(node *Tree[T])) {
	if t == nil || handler == nil {
		return
	}
	DfsPost(t.Left, handler)
	DfsPost(t.Right, handler)
	handler(t)
}

// Bfs 广度优先遍历
func Bfs[T any](t *Tree[T], handler func(node *Tree[T])) {
	if t == nil || handler == nil {
		return
	}

	handler(t)

	var cur []*Tree[T]
	cur = appendIfAbsent(cur, t.Left)
	cur = appendIfAbsent(cur, t.Right)
	for {
		var nextLayer []*Tree[T]
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

func appendIfAbsent[T any](layer []*Tree[T], node *Tree[T]) []*Tree[T] {
	if node == nil {
		return layer
	}
	return append(layer, node)
}

func PrintNode[T any](node *Tree[T]) {
	fmt.Println(node.Data)
}