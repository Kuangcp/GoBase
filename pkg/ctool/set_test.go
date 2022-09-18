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

func TestSet_Intersect(t *testing.T) {
	set := NewSet(3, 4)
	set.Add(3, 5)

	b := NewSet(1, 2, 3, 4)
	intersection := set.Intersect(b)
	fmt.Println(set, b, intersection)
}

func TestSet_Difference(t *testing.T) {
	a := NewSet(1, 2, 3)
	b := NewSet(3, 5, 4)

	result := a.Difference(b)
	fmt.Println(a, b, result)
}

func TestSet_Union(t *testing.T) {
	a := NewSet(1, 2, 3)
	b := NewSet(4, 5)
	union := a.Union(b)
	fmt.Println(a, b, union)
}
