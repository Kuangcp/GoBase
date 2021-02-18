package queue

import (
	"fmt"
	"testing"
	"time"
)

func TestPush(t *testing.T) {

	queue := New()
	//queue.Push(12)
	//tmp := queue.Pop()
	//fmt.Println((*tmp).(int))

	push := time.NewTicker(time.Millisecond * 200)
	check := time.NewTicker(time.Millisecond * 500)
	go func() {
		for now := range push.C {
			queue.Push(now.Unix())
		}
	}()

	var window int64 = 5
	for now := range check.C {
		for {
			peek := queue.Peek()
			if peek == nil {
				break
			}
			nowSec := now.Unix()
			peekVal := (*peek).(int64)
			if nowSec-peekVal > window {
				queue.Pop()
				fmt.Println(peekVal, queue.Len())
			} else {
				break
			}
		}
	}
}
