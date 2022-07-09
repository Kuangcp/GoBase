package ctool

import (
	"fmt"
	"testing"
)

func TestLoop(t *testing.T) {
	v := make(map[string]string)
	v["xx"] = "rr"

	for s := range v {
		fmt.Println(s)
	}
}

func TestNew(t *testing.T) {
	set := NewSet[int](4, 5)
	set.Loop(func(i int) {
		fmt.Println(i)
	})

	set.Add(5)

	set.Loop(func(i int) { fmt.Println(i) })
}
