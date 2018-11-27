package util

import (
	"testing"
)

func TestGenarateGrid(t *testing.T) {

	grid := new(GenerateGrid)
	data := grid.ReadConfig()
	for _, line := range data {
		line.toString()
	}
}
