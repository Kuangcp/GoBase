package ctool

import (
	"fmt"
	"testing"
	"time"
)

func watchNs(title string, act func()) {
	start := time.Now().UnixMilli()
	act()
	end := time.Now().UnixMilli()
	fmt.Printf("%s %d ms\n", title, end-start)
}

func TestStr(t *testing.T) {
	println(RandomAlpha(10))

	time.Sleep(time.Nanosecond * 20)
	l := 7
	loop := 1000000
	watchNs("one", func() {
		langs := NewSet[string]()
		for i := 0; i < loop; i++ {
			langs.Add(RandomAlpha(l))
		}
		println(langs.Len())
	})
}

func TestNum(t *testing.T) {
	RandomAlNum(8)

	langs := NewSet[string]()
	time.Sleep(time.Nanosecond * 20)
	for i := 0; i < 10000; i++ {
		langs.Add(RandomAlNum(4))
	}
	println(langs.Len())
}

func TestNameValid(t *testing.T) {
	langs := NewSet[string]()
	for i := 0; i < 10000; i++ {
		valid := RandomAlNumValid(5)
		langs.Add(valid)
	}
	println(langs.Len())
}
