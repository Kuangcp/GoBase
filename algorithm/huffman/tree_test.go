package huffman

import (
	"testing"
)

func TestNewHuffmanTree(t *testing.T) {
	tree := NewHuffmanTree("aaaaaabbccccjjjjj")
	tree.print()
}
