package stream

import (
	"fmt"
	"testing"
)

func TestToSet2(t *testing.T) {
	s := Just(1, 3, 4, 3, 2, 1)
	x := Collect(s, Set[int]())
	fmt.Println(x)
}
