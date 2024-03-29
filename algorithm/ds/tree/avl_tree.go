package main

import (
	"fmt"
	"github.com/kuangcp/gobase/pkg/ctool"
	"github.com/kuangcp/gobase/pkg/ctool/algo"
	"math"
)

type (
	// https://en.wikipedia.org/wiki/AVL_tree
	// https://yufengbiji.com/posts/data-structures-09

	// AvlTree 自平衡二叉搜索树  https://en.wikipedia.org/wiki/Self-balancing_binary_search_tree
	// 平衡性的定义是指：以 T 为根节点的树，每一个节点的左子树和右子树高度差最多为 1。
	// 显然地，增加了树维护成本，但是使得查询的成本变得均匀
	AvlTree[T ctool.Ordered] struct {
		Root *BsNode[T]
	}
)

func InitAvlTree[T ctool.Ordered](val ...T) *AvlTree[T] {
	tree := &AvlTree[T]{}
	if len(val) > 0 {
		tree.Inserts(val...)
	}
	return tree
}

func (t *BsNode[T]) GetLeft() algo.IBinTree {
	return t.Left
}

func (t *BsNode[T]) GetRight() algo.IBinTree {
	return t.Right
}

func (t *BsNode[T]) ToString() string {
	return fmt.Sprint(t.Data)
}

func (b *AvlTree[T]) Lists() []T {
	var data []T
	return dfsIn(b.Root, data)
}

func (b *AvlTree[T]) Inserts(val ...T) {
	if len(val) == 0 {
		return
	}
	for _, v := range val {
		b.Insert(v)
	}
}

func (b *AvlTree[T]) Insert(val T) {
	if b.Root == nil {
		b.Root = newNode(val)
		return
	}

	// bst 插入
	insert(b.Root, val)
	// 旋转调整
	b.Root = reBalance(b.Root)
}

// reBalance 二叉平衡树在调整时需要保证中序遍历序列不变
// 依据左右子树的差异，可分为四种情况：LL LR RR RL。
// 又因为不平衡的情况仅由插入删除导致，每次不平衡后都处理为平衡，所以只会出现这四种不平衡情况（局部）
func reBalance[T ctool.Ordered](root *BsNode[T]) *BsNode[T] {
	if root == nil {
		return nil
	}
	// 递归完成对每个节点的整理
	root.Left = reBalance(root.Left)
	root.Right = reBalance(root.Right)

	lm := algo.Height(root.Left)
	rm := algo.Height(root.Right)
	diff := math.Abs(float64(lm - rm))
	if diff <= 1 {
		return root
	}

	// 注意此时高度差只会是2，不会大于2
	if lm > rm {
		llm := algo.Height(root.Left.Left)
		lrm := algo.Height(root.Left.Right)

		if llm < lrm {
			// LR 先左旋左节点使其转换为LL
			root.Left = rotateLeft(root.Left)
		}

		// LL
		root = rotateRight(root)
	} else {
		rlm := algo.Height(root.Right.Left)
		rrm := algo.Height(root.Right.Right)

		if rrm < rlm {
			// RL 先右旋右节点 使其转换为RR
			root.Right = rotateRight(root.Right)
		}
		// RR
		root = rotateLeft(root)
	}
	return root
}

func rotateLeft[T ctool.Ordered](root *BsNode[T]) *BsNode[T] {
	// 旧根右节点成为新根
	newRoot := root.Right
	// 新根左子树成为旧根右子树
	root.Right = newRoot.Left
	// 旧根成为新根的左子树
	newRoot.Left = root
	return newRoot
}

func rotateRight[T ctool.Ordered](root *BsNode[T]) *BsNode[T] {
	// 旧根左节点成为新根
	newRoot := root.Left
	// 新根右子树成为旧根左子树
	root.Left = newRoot.Right
	// 旧根成为新根的右子树
	newRoot.Right = root
	return newRoot
}
