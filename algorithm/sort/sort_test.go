package sort

import (
	"fmt"
	"github.com/kuangcp/gobase/pkg/stopwatch"
	"math/rand"
	"testing"
)

const (
	length = 20
	max    = 20000
)

func TestSort(t *testing.T) {
	var data []int
	rand.Seed(777)
	for i := 0; i < length; i++ {
		data = append(data, rand.Intn(max))
	}
	fmt.Println(data)
	watch := stopwatch.NewWithName("sort")
	watch.Start("merge")
	result := Merge(data)
	watch.Stop()

	fmt.Println(result)
	fmt.Println(watch.PrettyPrint())
}
