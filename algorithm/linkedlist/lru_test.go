package linkedlist

import (
	"fmt"
	"testing"
	"time"
)

func TestNormal(t *testing.T) {
	pool := NewLRUCache(2)
	pool.Save("x", "1")
	pool.Save("x2", "2")
	pool.Save("x3", "3")

	fmt.Println(pool.Size())
	fmt.Println(pool.Get("x"))
	fmt.Println(pool.Get("x2"))
	fmt.Println(pool.Get("x3"))
}

func TestConcurrency(t *testing.T) {
	pool := NewLRUCache(5)

	for i := 0; i < 5; i++ {
		go func() {
			for i := 0; i < 4; i++ {
				pool.Save(fmt.Sprint(time.Now().UnixNano()), "x")
				fmt.Println(pool.Size())
			}
		}()
	}
	time.Sleep(time.Second * 3)
	fmt.Println(pool.Size())
}
