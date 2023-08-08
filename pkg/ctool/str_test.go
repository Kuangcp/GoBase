package ctool

import (
	"testing"
	"time"
)

func TestStr(t *testing.T) {
	println(RandomAlpha(10))

	langs := NewSet[string]()

	time.Sleep(time.Nanosecond * 20)
	for i := 0; i < 10000; i++ {
		langs.Add(RandomAlpha(4))
	}
	println(langs.Len())
}

func TestNum(t *testing.T) {
	for i := 0; i < 10; i++ {
		println(RandomAlNum(8))
	}
}

func TestNameValid(t *testing.T) {
	for i := 0; i < 10000; i++ {
		valid := RandomAlNumValid(5)
		println(valid)
	}
}
