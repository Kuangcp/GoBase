package main

import "github.com/kuangcp/gobase/pkg/ctool"

type (
	// BsTree 二叉搜索树 binary search tree https://oi-wiki.org/ds/bst/
	BsTree[T ctool.Numberic] struct {
		Root *BsNode[T]
	}
	BsNode[T ctool.Numberic] struct {
		Left  *BsNode[T]
		Right *BsNode[T]
		Data  T
		Count int
		Size  int
	}
)

func InitBsTree[T ctool.Numberic]() *BsTree[T] {
	return &BsTree[T]{}
}

func (b *BsTree[T]) Insert(val T) {
	if b.Root == nil {
		b.Root = newNode(val)
		return
	}

	insert(b.Root, val)
}
func (b *BsTree[T]) Delete(val T) {
	delete(b.Root, val)

}
func (b *BsTree[T]) Search(val T) *BsNode[T] {
	return search(b.Root, val)
}

func (b *BsTree[T]) Min() T {
	if b.Root == nil {
		return 0
	}
	var val T
	p := b.Root
	for p != nil {
		val = p.Data
		p = p.Left
	}
	return val
}

func (b *BsTree[T]) Max() T {
	if b.Root == nil {
		return 0
	}
	var val T
	p := b.Root
	for p != nil {
		val = p.Data
		p = p.Right
	}
	return val
}

func newNode[T ctool.Numberic](val T) *BsNode[T] {
	return &BsNode[T]{Data: val, Count: 1}
}

func delete[T ctool.Numberic](b *BsNode[T], val T) *BsNode[T] {
	if b == nil {
		return nil
	}
	if b.Data == val {
		if b.Count > 1 {
			b.Count -= 1
			return nil
		} else {
			return b
		}
	}
	if val < b.Data {
		return delete(b.Left, val)
	} else {
		return delete(b.Right, val)
	}
}

func search[T ctool.Numberic](b *BsNode[T], val T) *BsNode[T] {
	if b == nil {
		return nil
	}
	if b.Data == val {
		return b
	}
	if val < b.Data {
		return search(b.Left, val)
	} else {
		return search(b.Right, val)
	}
}

func insert[T ctool.Numberic](b *BsNode[T], val T) {
	if b.Data == val {
		b.Count += 1
		return
	}

	if val < b.Data {
		if b.Left == nil {
			b.Left = newNode(val)
		} else {
			insert(b.Left, val)
		}
	} else {
		if b.Right == nil {
			b.Right = newNode(val)
		} else {
			insert(b.Right, val)
		}
	}
}
