package sizedpool

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
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

func TestQueue(t *testing.T) {
	run, _ := NewQueuePool(2)
	for i := 0; i < 7; i++ {
		fi := i
		run.Submit(func() {
			time.Sleep(time.Second * 2)
			log.Println("task run", fi)
		})
		log.Println("submit", i)
	}
	log.Println("submit all")
	time.Sleep(time.Second * 1)
	run.Wait()
}

func TestFuture(t *testing.T) {
	future, _ := New(PoolOption{size: 3})
	var res []*FutureTask
	for i := 0; i < 80; i++ {
		submitFuture := future.SubmitFutureTimeout(time.Second*6, Callable{
			fmt.Sprint(i),
			func(ctx context.Context) (interface{}, error) {
				time.Sleep(time.Second * 1)
				sec := time.Now().Second()
				if sec%2 == 0 {
					return sec, nil
				}
				return nil, errors.New("oo")
			}, func(data interface{}) {
				log.Println("success:", data)
			}, func(ex error) {
				log.Println("fail:", ex)
			}})
		res = append(res, submitFuture)
	}
	time.Sleep(time.Second * 2)
	future.Wait()

	for _, re := range res {
		log.Println(re.GetData())
	}
}

func TestFutureGet(t *testing.T) {
	type VO struct {
		id   int
		name string
	}
	future, _ := New(PoolOption{size: 2})
	var res []*FutureTask
	for i := 0; i < 7; i++ {
		fi := i
		submitFuture := future.SubmitFutureTimeout(time.Second*6, Callable{
			fmt.Sprint(fi),
			func(ctx context.Context) (interface{}, error) {
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(900)))
				sec := time.Now().Second()
				if sec%4 == 0 {
					return VO{id: fi, name: fmt.Sprint("test", sec)}, nil
				}
				return nil, errors.New(fmt.Sprint(fi, " exception"))
			}, func(data interface{}) {
				log.Println("su call", data)
			}, func(ex error) {
				log.Println("fa call", ex)
			},
		})

		res = append(res, submitFuture)
	}

	go func() {
		for _, re := range res {
			f := re
			go func() {
				//data, err := f.GetData()
				// 超时不等待，但是任务还在执行
				data, err := f.GetDataTimeout(time.Millisecond * 600)
				log.Println("future get", data, err)
			}()
		}
	}()

	go future.ExecFuturePool(nil)
	http.ListenAndServe(":9090", nil)
}

func TestFutureGetWithCancel(t *testing.T) {
	type VO struct {
		id   int
		name string
	}
	future, _ := NewFuturePool(PoolOption{size: 6})

	var res []*FutureTask
	for i := 0; i < 30; i++ {
		fi := i

		submitFuture := future.SubmitFuture(Callable{
			fmt.Sprint(fi),
			func(ctx context.Context) (interface{}, error) {
				//submitFuture := future.SubmitFutureTimeout(time.Second*2, func() (interface{}, error) {
				x := rand.Intn(900) + 1600
				//fmt.Println(fi, x)
				time.Sleep(time.Millisecond * time.Duration(x))
				sec := time.Now().Second()
				//if sec%4 == 0 {
				return VO{id: fi, name: fmt.Sprint("test", sec)}, nil
				//}
				//return nil, errors.New(fmt.Sprint(fi, " exception"))
			}, func(data interface{}) {
				log.Println("su call", data)
			}, func(ex error) {
				log.Println("fa call", ex)
			},
		})

		res = append(res, submitFuture)
	}

	go func() {
		for _, re := range res {
			f := re
			go func() {
				data, err := f.GetData()
				// 超时未获取到结果就返回，但是任务还在执行
				//data, err := f.GetDataTimeout(time.Millisecond * 2300)
				log.Println("future get", data, err)
			}()
		}
	}()

	time.Sleep(time.Second * 5)
	future.Wait()
}

func TestNewTmpWithFuture(t *testing.T) {
	log.Println("start")
	//future, _ := NewTmpWithFuture(30, time.Second*4)
	future, err := NewTmpFuturePool(PoolOption{size: 30, timeout: time.Second * 7})
	if err != nil {
		log.Println(err)
		return
	}

	for i := 0; i < 3; i++ {
		fi := i
		future.SubmitFutureTimeout(time.Second*5, Callable{
			fmt.Sprint(fi),
			func(ctx context.Context) (interface{}, error) {
				value := ctx.Value(TraceID)
				log.Println(value, "start")
				sl := rand.Intn(4) + 10
				time.Sleep(time.Second * time.Duration(sl))
				log.Println(value, "finish", sl)
				return fi, nil
			}, nil, nil,
		})
	}
	future.Wait()
}
