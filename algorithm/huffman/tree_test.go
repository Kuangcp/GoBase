package huffman

import (
	"fmt"
	"github.com/theodesp/go-heaps"
	pairingHeap "github.com/theodesp/go-heaps/pairing"
	"testing"
)

func TestNewHuffmanTree(t *testing.T) {
	tree := NewHuffmanTree("abgts")
	tree.print()
	heap := pairingHeap.New()
	heap.Insert(go_heaps.Integer(4))
	heap.Insert(go_heaps.Integer(2))

	fmt.Printf("%v", heap.DeleteMin())
	fmt.Printf("%v", heap.DeleteMin())
}
