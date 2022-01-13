package queue

import (
	"fmt"
	"testing"
	"time"
)

func TestRollingWindow(t *testing.T) {
	queue := New()

	push := time.NewTicker(time.Millisecond * 20)
	check := time.NewTicker(time.Millisecond * 50)
	go func() {
		for now := range push.C {
			queue.Push(now.Unix())
		}
	}()

	var windowSecSize int64 = 5
	for now := range check.C {
		for {
			peek := queue.Peek()
			if peek == nil {
				break
			}
			nowSec := now.Unix()
			peekVal := (*peek).(int64)
			if nowSec-peekVal > windowSecSize {
				queue.Pop()
				fmt.Println(peekVal, queue.Len())
			} else {
				break
			}
		}
	}
}
