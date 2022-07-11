package ctool

import (
	"fmt"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMaxSize(t *testing.T) {
	a := assert.New(t)

	cache := NewLRUCache[string](2)
	cache.Save("1", "")
	a.Equal(cache.Get("1"), "")
	cache.Save("2", "2")
	cache.Save("3", "3")

	a.Equal(cache.Get("1"), "")

	a.Equal(cache.Size(), 2)
	a.Equal(cache.Get("2"), "2")
	a.Equal(cache.Get("3"), "3")
}

type vo struct {
	id   string
	name string
}

func TestMaxSize2(t *testing.T) {
	a := assert.New(t)

	cache := NewLRUCache[vo](2)
	cache.Save("1", vo{id: "???ss"})
	a.Equal(cache.Get("1"), vo{id: "???ss"})

	cache.Save("2", vo{})
	cache.Save("3", vo{})

	a.Equal(cache.Get("1"), vo{})

	a.Equal(cache.Size(), 2)
	a.Equal(cache.Get("2"), "2")
	a.Equal(cache.Get("3"), "3")
}

func TestArray(t *testing.T) {
	// 要求实现 comparable 但是结构体数组无法实现 comparable， 只能 reflect.DeepEqual(v1, v2)来比较
	cache := NewLRUCache[[]vo](2)
	cache.Save("x", []vo{{id: "xx"}})

	vos := cache.Get("x")
	fmt.Println(vos)
}

func TestLRU(t *testing.T) {
	a := assert.New(t)

	cache := NewLRUCache[string](2)
	cache.Save("1", "1")
	cache.Save("2", "2")
	cache.Get("1")
	cache.Save("3", "3")

	a.Equal(cache.Size(), 2)
	a.Equal(cache.Get("1"), "1")
	a.Nil(cache.Get("2"))
	a.Equal(cache.Get("3"), "3")
}

func TestConcurrency(t *testing.T) {
	a := assert.New(t)

	maxSize := 5
	cache := NewLRUCache[string](maxSize)
	for i := 0; i < 50; i++ {
		go func() {
			for i := 0; i < 200000; i++ {
				cache.Save(fmt.Sprint(time.Now().UnixNano()), "")
				a.LessOrEqual(cache.Size(), maxSize)
			}
			printMemStats()
		}()
		time.Sleep(time.Second * 1)
	}
	time.Sleep(time.Second * 120)
	fmt.Println("max:", cache.Size())
}

// Alloc： 当前堆上对象占用的内存大小;
// TotalAlloc：堆上总共分配出的内存大小;
// Sys： 程序从操作系统总共申请的内存大小;
// NumGC： 垃圾回收运行的次数。
func printMemStats() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("Alloc = %v TotalAlloc = %v Sys = %v NumGC = %v\n",
		m.Alloc/1024, m.TotalAlloc/1024, m.Sys/1024, m.NumGC)
}
