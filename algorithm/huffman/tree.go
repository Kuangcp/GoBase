package huffman

import (
	"fmt"
	"github.com/kuangcp/gobase/pkg/ctool/algo"
	"sort"
	"strings"
)

type (
	Node struct {
		Char  uint8
		Code  uint8
		Sum   int
		Left  *Node
		Right *Node
	}

	HuffmanTree struct {
		Root        *Node
		Len         int
		hiddenEmpty bool
		printMind   bool
	}
)

func (n *Node) GetLeft() algo.IBinTree {
	return n.Left
}

func (n *Node) GetRight() algo.IBinTree {
	return n.Right
}

func (n *Node) ToString() string {
	c := "🔹"
	if n.Char != 0 {
		c = string(n.Char)
	}
	return fmt.Sprintf("%v %v", n.Sum, c)
}

func makeSort(list []*Node) {
	sort.Slice(list, func(i, j int) bool {
		return list[i].Sum < list[j].Sum
	})
}

func NewHuffmanTree(formatStr string) *HuffmanTree {
	cache := make(map[uint8]*Node)
	for i := range formatStr {
		char := formatStr[i]
		_, ok := cache[char]
		if !ok {
			cache[char] = &Node{Char: char, Sum: 1}
		} else {
			cache[char].Sum++
		}
	}

	var l []*Node
	for _, v := range cache {
		l = append(l, v)
	}
	makeSort(l)

	tree := &HuffmanTree{Len: len(l)}
	if len(l) == 1 {
		tree.Root = &Node{Char: l[0].Char}
		return tree
	}

	for len(l) > 1 {
		left := l[0]
		right := l[1]
		left.Code = 1
		right.Code = 0

		l = l[2:]
		l = append(l, &Node{Char: 0, Sum: left.Sum + right.Sum, Left: left, Right: right})
		makeSort(l)
	}
	tree.Root = l[0]
	return tree
}

func (h *HuffmanTree) MarkHiddenEmpty() {
	h.hiddenEmpty = true
}
func (h *HuffmanTree) MarkMindMap() {
	h.printMind = true
}

func (t *HuffmanTree) Print() {
	if t.Root == nil {
		return
	}
	t.printNode(t.Root, "")
}
func (t *HuffmanTree) printNode(n *Node, path string) {
	if n == nil {
		return
	}
	curPath := fmt.Sprintf("%v%v", path, n.Code)
	if !(t.hiddenEmpty && n.Char == 0) {
		if t.printMind {
			c := "🚫"
			if n.Char != 0 {
				c = string(n.Char)
			}
			fmt.Println(strings.Repeat("*", len(curPath)), fmt.Sprintf("%v[%v]", n.Sum, c))
		} else {
			fmt.Println(curPath, string(n.Char))
		}
	}
	if n.Left != nil {
		t.printNode(n.Left, curPath)
	}
	if n.Right != nil {
		t.printNode(n.Right, curPath)
	}
}
