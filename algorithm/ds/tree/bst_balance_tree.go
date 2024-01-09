package main

import "github.com/kuangcp/gobase/pkg/ctool"

type (
	// BsBalanceTree 平衡二叉搜索树  https://en.wikipedia.org/wiki/Self-balancing_binary_search_tree
	// 平衡性的定义是指：以 T 为根节点的树，每一个结点的左子树和右子树高度差最多为 1。
	// https://www.javatpoint.com/binary-search-tree-vs-avl-tree
	BsBalanceTree[T ctool.Numberic] struct {
		Root *BsNode[T]
	}
)

func InitBsBalanceTree[T ctool.Numberic]() *BsBalanceTree[T] {
	return &BsBalanceTree[T]{}
}

func (b *BsBalanceTree[T]) Inserts(val ...T) {
	if len(val) == 0 {
		return
	}
	for _, v := range val {
		b.Insert(v)
	}
}
func (b *BsBalanceTree[T]) Insert(val T) {
	if b.Root == nil {
		b.Root = newNode(val)
		return
	}

	insert(b.Root, val)
}
