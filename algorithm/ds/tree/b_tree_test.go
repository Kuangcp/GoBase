package main

import (
	"math"
	"testing"
)

func TestCeil(t *testing.T) {
	println(math.Floor(7 / 3))
	println(math.Ceil(7 / 3))
}

func TestInitBayerTree(t *testing.T) {
	tree := CreateBayerTree[int, string](5)

	tree.Insert(3, "3key")
}
