package main

import (
	"fmt"
	"github.com/kuangcp/gobase/pkg/ctool/algo"
	"golang.org/x/exp/constraints"
	"math"
)

// https://en.wikipedia.org/wiki/B-tree
// https://zhuanlan.zhihu.com/p/340721689 相比于AVL树放硬盘降低树高度更实用
type (
	Order           constraints.Ordered
	BValNode[V any] struct {
		val  V
		pre  *BValNode[V]
		next *BValNode[V]
	}
	BayerEntry[K Order, V any] struct {
		key  K            // 索引字段
		val  V            // 记录值
		head *BValNode[V] // key值重复时溢出值链表 或者溢出块的链表(单个node有多个val) 提高缓存命中, 注意头插法读取时需要逆序
	}
	BayerNode[K Order, V any] struct {
		keyN    int // 有效key数量 1 2 3 ...
		Keys    []*BayerEntry[K, V]
		Childes []*BayerNode[K, V] // 例如m为5时 单个节点上 3个key 4个child时的逻辑结构: c0 k0 c1 k1 c2 k3 c4
		Parent  *BayerNode[K, V]
	}
	BayerTree[K Order, V any] struct {
		m        int // 树阶数
		minChild int // 每个节点最少子节点数: (m/2 向上取整)
		minKey   int // 最少key: minChild -1 注意: 根节点最少为1个
		maxKey   int // 最多key: m -1
		root     *BayerNode[K, V]
	}
)

func (b *BayerNode[K, V]) IsLeafNode() bool {
	if len(b.Childes) == 0 {
		return true
	}

	for i := range b.Childes {
		if b.Childes[i] != nil {
			return false
		}
	}
	return true
}

func (b *BayerNode[K, V]) ToString() string {
	re := ""
	for i := range b.Keys {
		t := b.Keys[i]
		if t == nil {
			continue
		}
		re += fmt.Sprint(b.Keys[i].key) + " "
	}
	return re
}

func (b *BayerNode[K, V]) GetChild() []algo.INTree {
	var cs []algo.INTree
	for _, childe := range b.Childes {
		cs = append(cs, childe)
	}
	return cs
}

func CreateBayerTree[K Order, V any](m int) *BayerTree[K, V] {
	minChild := int(math.Ceil(float64(m)/2) + 1)
	return &BayerTree[K, V]{m: m, minChild: minChild, minKey: minChild - 1, maxKey: m - 1}
}

func (t *BayerTree[K, V]) Insert(key K, val V) {
	if t.root == nil {
		keys := make([]*BayerEntry[K, V], t.m-1)
		childes := make([]*BayerNode[K, V], t.m)
		keys[0] = &BayerEntry[K, V]{key: key, val: val}
		t.root = &BayerNode[K, V]{Keys: keys, Childes: childes}
		t.root.keyN = 1
		return
	}
	node, idx, match := searchBayerNode(t.root, key)
	if match {
		entry := node.Keys[idx]
		var exist *BValNode[V]
		if entry.head == nil {
			exist = &BValNode[V]{val: entry.val}
		} else {
			exist = entry.head
		}
		entry.head = &BValNode[V]{val: val, next: exist}
		exist.pre = entry.head
	} else {
		// TODO  key的插入以及拆分节点
		if node.keyN >= t.maxKey {
			mid := t.maxKey / 2

			rc := -1
			var tmpKey *BayerEntry[K, V]
			for i := range node.Keys {
				if i < idx {
					rc++
					continue
				} else if i == idx {
					newKey := &BayerEntry[K, V]{key: key, val: val}
					rc++
					if rc == mid {
						// 提升中间节点到父节点去， 判断父节点空的情况
						node.Parent.Keys[node.Parent.keyN-1] = newKey
					} else {
						tmpKey = node.Keys[i]
						node.Keys[i] = newKey
					}
				} else {
					rc++
					if rc == mid {
						// 提升中间节点到父节点去
						node.Parent.Keys[node.Parent.keyN-1] = tmpKey
					} else {
						roundKey := node.Keys[i]
						node.Keys[i] = tmpKey
						tmpKey = roundKey
					}
				}

			}
			// TODO 新key补到右节点上 右节点从上面拆分来
			if idx == node.keyN {

			}

		} else {
			// 单纯插入key的情况 左边 中间 最右边
			var tmp *BayerEntry[K, V]
			for i, b := range node.Keys {
				if i < idx {
					continue
				}

				if i == idx {
					node.Keys[idx] = &BayerEntry[K, V]{key: key, val: val}
					node.keyN += 1
				} else {
					node.Keys[i] = tmp
				}
				tmp = b
				if tmp == nil {
					break
				}
			}
		}

	}
}

func (t *BayerTree[K, V]) Search(key K) *BayerEntry[K, V] {
	if t.root == nil {
		return nil
	}
	node, i, match := searchBayerNode(t.root, key)
	if !match {
		return nil
	}
	return node.Keys[i]
}

// searchBayerNode 搜索key值: 设计思路为了命中key, 为了插入使用
// 返回值:
// 1 key所在节点
// 2 key值命中则返回所在下标, 否则返回插入位的右下标, 如果值大于已有的key 需要调用方调整树结构层级
// 3 key如果已存在则返回true
func searchBayerNode[K Order, V any](node *BayerNode[K, V], key K) (*BayerNode[K, V], int, bool) {
	i := 0
	for i = range node.Keys {
		cur := node.Keys[i]
		if cur == nil {
			break
		}
		curK := cur.key
		if key == curK {
			return node, i, true
		}
		if key < curK {
			if node.IsLeafNode() {
				return node, i, false
			}
			return searchBayerNode(node.Childes[i], key)
		}
	}
	if node.IsLeafNode() {
		return node, node.keyN, false
	} else {
		return searchBayerNode(node.Childes[i+1], key)
	}
}
