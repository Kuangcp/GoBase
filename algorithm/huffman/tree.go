package huffman

import "fmt"

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

func NewHuffmanTree(formatStr string) HuffmanTree {
	tree := HuffmanTree{Len: len(formatStr)}
	rootPtr := &Node{
		Code: 0,
	}

	for i := range formatStr {
		char := formatStr[i]
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
	var ptr = t.Root
	fmt.Println(string(t.Root.Char), t.Root.Code)
	for {
		if ptr.Left == nil {
			break
		} else {
			fmt.Println(string(ptr.Left.Char), ptr.Left.Code)
		}
		if ptr.Right != nil {
			fmt.Println(string(ptr.Right.Char), ptr.Right.Code)
			ptr = ptr.Left
		}
	}
}

func Build(){

}