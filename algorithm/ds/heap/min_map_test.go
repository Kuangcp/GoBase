package heap

import (
	"fmt"
	"testing"
)

type Score struct {
	score int
	char  string
}

func (t *Score) Value() int {
	return t.score
}

func TestNewMinHeap(t *testing.T) {
	heap := MinHeap{}
	heap.insertValue(NewIntItem(0))
	for i := 6; i > 0; i-- {
		score := Score{score: i, char: "test"}
		heap.insertValue(&score)
	}

	for _, value := range heap.Values {
		fmt.Println(value)
	}
}

