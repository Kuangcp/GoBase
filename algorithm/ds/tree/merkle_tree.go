package main

import (
	"crypto/md5"
	"fmt"
)

// 基于二叉树实现 哈希树
// Git 用于比对文件系统差异时提升效率
// P2P下载中BT种子的构成是hash列表(按块拆分文件后做hash得到)
// https://en.wikipedia.org/wiki/Merkle_tree
// https://yeasy.gitbook.io/blockchain_guide/05_crypto/merkle_trie

type (
	Merkle struct {
		tree *BinTree[string]
	}
)

func CreateMerkleByFile(path string) *Merkle {
	// 对文件分块计算Hash构造树并存储
	return nil
}

func CreateMerkle(hash []string) *Merkle {
	if len(hash) == 0 {
		return nil
	}

	var tmp []*BinTree[string]
	for _, s := range hash {
		tmp = append(tmp, &BinTree[string]{Data: s})
	}
	for {
		if len(tmp) == 1 {
			break
		}
		var par []*BinTree[string]
		var node *BinTree[string]
		for i := range tmp {
			if node == nil {
				node = &BinTree[string]{Left: tmp[i]}
				par = append(par, node)
			} else {
				node.Right = tmp[i]
				node = nil
			}
		}
		for _, b := range par {
			fillData(b)
		}
		tmp = par
	}
	return &Merkle{tree: tmp[0]}
}
func (m *Merkle) Same(merkle *Merkle) bool {
	if merkle == nil {
		return false
	}
	return m.tree.Data == merkle.tree.Data
}
func fillData(node *BinTree[string]) {
	if node == nil {
		return
	}
	c := node.Left.Data
	if node.Right != nil {
		c += " " + node.Right.Data
	}

	node.Data = hash(c)
}
func hash(val string) string {
	sum := md5.Sum([]byte(val))
	return fmt.Sprintf("%x", sum)
}
