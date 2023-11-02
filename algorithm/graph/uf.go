package graph

type UF struct {
	// 连通分量个数
	count int
	// 存储每个节点的父节点
	parent []int
}

// n 为图中节点的个数
func NewUF(n int) *UF {
	parent := make([]int, n)
	for i := 0; i < n; i++ {
		parent[i] = i
	}
	return &UF{
		count:  n,
		parent: parent,
	}
}

// 将节点 p 和节点 q 连通
func (u *UF) Union(p, q int) {
	rootP := u.Find(p)
	rootQ := u.Find(q)

	if rootP == rootQ {
		return
	}

	// 一个树的根指向另一个树的根， 即挂接成为一个子树导致树高度为3，但是随着下次路径压缩的Find调用 会让树的高度压缩回到2
	u.parent[rootQ] = rootP

	// 两个连通分量合并成一个连通分量
	u.count--
}

// 判断节点 p 和节点 q 是否连通
func (u *UF) Connected(p, q int) bool {
	rootP := u.Find(p)
	rootQ := u.Find(q)
	return rootP == rootQ
}

// 路径压缩 查找
func (u *UF) Find(x int) int {
	if u.parent[x] != x {
		// 压缩层级
		u.parent[x] = u.Find(u.parent[x])
	}
	return u.parent[x]
}

// 返回图中的连通分量个数
func (u *UF) Count() int {
	return u.count
}
