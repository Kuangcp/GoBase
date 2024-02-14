package main

import (
	"github.com/kuangcp/gobase/pkg/ctool/algo"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInitMerkleByHash(t *testing.T) {
	merkle := CreateMerkle([]string{"a", "b", "c", "d", "e", "f"})

	println(algo.PrintBiMindMap(merkle.tree))
}

func TestMerkle_Same(t *testing.T) {
	a := CreateMerkle([]string{"a", "b", "c", "d", "e", "f"})
	println(algo.PrintBiMindMap(a.tree))

	b := CreateMerkle([]string{"a", "b", "c", "d", "e", "f"})
	println(algo.PrintBiMindMap(b.tree))

	c := CreateMerkle([]string{"a", "b", "c", "d", "e"})
	println(algo.PrintBiMindMap(c.tree))

	assert.Equal(t, a.Same(b), true)
	assert.Equal(t, a.Same(c), false)
}
