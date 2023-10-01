package sort

import (
	"fmt"
	"github.com/kuangcp/gobase/pkg/ctool"
	"math/rand"
	"testing"
)

const (
	length = 200000
	maxVal = 2000000
)

func TestSort(t *testing.T) {
	var data []int

	watch := ctool.NewStopWatchWithName("sort")
	watch.Start("init")
	rand.NewSource(7799)
	for i := 0; i < length; i++ {
		data = append(data, rand.Intn(maxVal))
	}
	watch.Stop()
	//fmt.Println(data)

	watch.Start("merge")
	result := Merge(data)
	watch.Stop()

	for i := 1; i < len(result); i++ {
		if result[i-1] > result[i] {
			t.Failed()
		}
	}
	//fmt.Println(result)
	fmt.Println(watch.PrettyPrint())
}
