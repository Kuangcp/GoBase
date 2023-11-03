package huffman

import (
	"fmt"
	"sort"
)

type (
	Node struct {
		Char  uint8
		Code  uint8
		Left  *Node
		Right *Node
	}

	HuffmanTree struct {
		Root *Node
		Len  int
	}
)

func NewHuffmanTree(formatStr string) *HuffmanTree {
	type cc struct {
		char  uint8
		count int
	}

	cache := make(map[uint8]*cc)
	for i := range formatStr {
		char := formatStr[i]
		_, ok := cache[char]
		if !ok {
			cache[char] = &cc{char: char, count: 1}
		} else {
			cache[char].count++
		}
	}

	var l []*cc
	for _, v := range cache {
		l = append(l, v)
	}

	sort.Slice(l, func(i, j int) bool {
		return l[i].count < l[j].count
	})

	tree := &HuffmanTree{Len: len(l)}
	rootPtr := &Node{
		Code: 0,
	}
	rootPtr = nil

	for _, i := range l {
		char := i.char
		if rootPtr == nil {
			tree.Root = &Node{
				Char: char,
				Code: 0,
			}
			rootPtr = tree.Root
		} else if rootPtr.Left == nil {
			rootPtr.Left = &Node{
				Char: char,
				Code: 1,
			}
		} else if rootPtr.Right == nil {
			rootPtr.Right = &Node{
				Char: char,
				Code: 0,
			}
		} else {
			rootPtr = rootPtr.Left
			rootPtr.Left = &Node{
				Char: char,
				Code: 1,
			}
		}
	}
	return tree
}

func (t *HuffmanTree) print() {
	if t.Root == nil {
		return
	}
	printNode(t.Root, "")
}
func printNode(n *Node, path string) {
	if n == nil {
		return
	}
	curPath := fmt.Sprintf("%v%v", path, n.Code)
	fmt.Println(string(n.Char), curPath)
	if n.Left != nil {
		printNode(n.Left, curPath)
	}
	if n.Right != nil {
		printNode(n.Right, curPath)
	}
}

func Build() {

}
