package sort

import (
	"fmt"
	"github.com/kuangcp/gobase/pkg/ctool"
	"math/rand"
	"testing"
)

const (
	length = 20000
	maxVal = 2000000
)

var data []int
var sorts = []Algo{
	{name: "Merge", fun: Merge},
	{name: "Patience", fun: Patience},
}

// https://github.com/Kuangcp/JavaBase/tree/master/algorithms/src/test/java/com/github/kuangcp/sort
func init() {
	for i := 0; i < length; i++ {
		data = append(data, rand.Intn(maxVal))
	}
}
func TestCorrect(t *testing.T) {
	watch := ctool.NewStopWatchWithName("sort")
	for _, s := range sorts {
		watch.Start(s.name)
		validate(s.name, s.fun(data))
		watch.Stop()
	}
	fmt.Println(watch.PrettyPrint())
}

func validate(name string, result []int) {
	for i := 1; i < len(result); i++ {
		if result[i-1] > result[i] {
			panic(fmt.Sprintf("%v: index %v, %v>%v", name, i, result[i-1], result[i]))
		}
	}
}

func TestSort(t *testing.T) {
	watch := ctool.NewStopWatchWithName("sort")
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
