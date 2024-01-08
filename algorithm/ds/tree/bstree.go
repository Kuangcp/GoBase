package main

// 二叉搜索树 binary search tree https://oi-wiki.org/ds/bst/

// BsTreeInsert 入参t为nil时则创建树
func BsTreeInsert[T any](t *Tree[T], val T) *Tree[T] {
	if t == nil {
		return &Tree[T]{Data: val}
	}

	return t
}

func BsTreeDelete[T any](t *Tree[T], val T) {
	if t == nil {
		return
	}
}

func BsTreeFind[T any](t *Tree[T], val T) {
	if t == nil {
		return
	}
}
