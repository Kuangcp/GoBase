package huffman

import (
	"github.com/kuangcp/gobase/pkg/ctool/algo"
	"testing"
)

func TestNewHuffmanTree(t *testing.T) {
	tree := NewHuffmanTree("dataSyncTaskMapper.insert(dataSyncTask);   sitttttttt")
	//tree.MarkMindMap()
	//tree.Print()

	algo.WriteBiMindMap(tree.Root, "m1.pu")
}
