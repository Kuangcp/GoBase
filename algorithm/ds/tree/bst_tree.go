package main

import "github.com/kuangcp/gobase/pkg/ctool"

type (
	// BsTree 二叉搜索树 binary search tree https://oi-wiki.org/ds/bst/
	// 实现了有序的搜索树，但是还有一个问题是容易出现数据倾斜，因此有了 BsBalanceTree
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
func (b *BsTree[T]) Inserts(val ...T) {
	if len(val) == 0 {
		return
	}
	for _, v := range val {
		b.Insert(v)
	}
}
func (b *BsTree[T]) Insert(val T) {
	if b.Root == nil {
		b.Root = newNode(val)
		return
	}

	insert(b.Root, val)
}

func (b *BsTree[T]) Delete(val T) {
	b.Root = removeNode(b.Root, val)
}

func removeNode[T ctool.Numberic](b *BsNode[T], val T) *BsNode[T] {
	if b == nil {
		return nil
	}

	if val < b.Data {
		b.Left = removeNode(b.Left, val)
	} else if val > b.Data {
		b.Right = removeNode(b.Right, val)
	} else {
		if b.Count > 1 {
			b.Count -= 1
			return b
		} else {
			if b.Left == nil {
				return b.Right
			} else if b.Right == nil {
				return b.Left
			} else {
				// 重新选 左树最大值或者右树最小值作为中间值
				p := b.Left
				for p != nil {
					temp := p.Right
					if temp == nil {
						break
					}
					p = p.Right
				}
				b.Data = p.Data
				b.Count = p.Count
				b.Left = removeNode(b.Left, p.Data)
			}
		}
	}
	return b
}

func (b *BsTree[T]) Search(val T) *BsNode[T] {
	return search(b.Root, val)
}

func (b *BsTree[T]) Rank(val T) int {
	return rank(b.Root, val) + 1
}

func rank[T ctool.Numberic](node *BsNode[T], val T) int {
	if node == nil {
		return 0
	}
	if node.Data == val {
		return sumCount(node.Left)
	} else if val < node.Data {
		return rank(node.Left, val)
	} else {
		return sumCount(node.Left) + rank(node.Right, val) + node.Count
	}
}

func sumCount[T ctool.Numberic](node *BsNode[T]) int {
	if node == nil {
		return 0
	}
	return sumCount(node.Left) + sumCount(node.Right) + node.Count
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
