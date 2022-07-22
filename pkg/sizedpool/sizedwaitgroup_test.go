package sizedpool

import (
	"context"
	"fmt"
	"strconv"
	"sync/atomic"
	"testing"
	"time"
)

func TestWait(t *testing.T) {
	swg, _ := New(PoolOption{size: 10})
	var c uint32

	for i := 0; i < 10000; i++ {
		swg.Add()
		go func(c *uint32) {
			defer swg.Done()
			atomic.AddUint32(c, 1)
		}(&c)
	}

	swg.Wait()

	if c != 10000 {
		t.Fatalf("%d, not all routines have been executed.", c)
	}
}

func TestThrottling(t *testing.T) {
	var c uint32

	swg, _ := New(PoolOption{size: 4})

	if len(swg.current) != 0 {
		t.Fatalf("the SizedWaitGroup should start with zero.")
	}

	for i := 0; i < 10000; i++ {
		swg.Add()
		go func(c *uint32) {
			defer swg.Done()
			atomic.AddUint32(c, 1)
			if len(swg.current) > 4 {
				t.Fatalf("not the good amount of routines spawned.")
				return
			}
		}(&c)
	}

	swg.Wait()
}

func TestNoThrottling(t *testing.T) {
	var c uint32
	swg, _ := New(PoolOption{size: 0})
	if len(swg.current) != 0 {
		t.Fatalf("the SizedWaitGroup should start with zero.")
	}
	for i := 0; i < 10000; i++ {
		swg.Add()
		go func(c *uint32) {
			defer swg.Done()
			atomic.AddUint32(c, 1)
		}(&c)
	}
	swg.Wait()
	if c != 10000 {
		t.Fatalf("%d, not all routines have been executed.", c)
	}
}

func TestAddWithContext(t *testing.T) {
	ctx, cancelFunc := context.WithCancel(context.TODO())

	swg, _ := New(PoolOption{size: 1})

	if err := swg.AddWithContext(ctx); err != nil {
		t.Fatalf("AddContext returned error: %v", err)
	}

	cancelFunc()
	if err := swg.AddWithContext(ctx); err != context.Canceled {
		t.Fatalf("AddContext returned non-context.Canceled error: %v", err)
	}
}

func TestRun(t *testing.T) {
	var size int64 = 3
	var loop int64 = 12
	start := time.Now().Unix()
	swg, _ := NewWithName(3, "sleep-group")
	for i := 0; i < 12; i++ {
		index := strconv.Itoa(i)
		swg.Run(func() {
			fmt.Println(swg.GetName(), "run", index)
			time.Sleep(time.Second * 1)
		})
	}
	swg.Wait()
	end := time.Now().Unix()
	if end-start < loop/size {
		t.Fatal("Not sleep enough time")
	}
}
