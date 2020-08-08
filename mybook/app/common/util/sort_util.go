package util

import (
	"sort"
)

type Comparable interface {
	CompareLess(b Comparable) bool
}

type SortWrapper struct {
	Data            []interface{}
	CompareLessFunc func(a interface{}, b interface{}) bool
	Reverse         bool
}

func (pw SortWrapper) Len() int {
	return len(pw.Data)
}

func (pw SortWrapper) Swap(i, j int) {
	pw.Data[i], pw.Data[j] = pw.Data[j], pw.Data[i]
}

func (pw SortWrapper) Less(i, j int) bool {
	result := pw.CompareLessFunc(pw.Data[i], pw.Data[j])
	if pw.Reverse {
		return !result
	}
	return result
}

func Sort(wrapper SortWrapper) {
	sort.Sort(wrapper)
}
