package sizedpool

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"runtime"
	"sync"
	"testing"
	"time"

	_ "net/http/pprof"
)

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
	future, _ := New(PoolOption{Size: 3})
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
			}, func(ctx context.Context, data interface{}) {
				log.Println("success:", data)
			}, func(ctx context.Context, ex error) {
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

// TraceId 传递及 超时控制
func TestFutureGet(t *testing.T) {
	type VO struct {
		id   int
		name string
	}
	future, _ := New(PoolOption{Size: 2})
	go func() {
		var res []*FutureTask
		for i := 0; i < 7; i++ {
			fi := i
			// 限制并发 提交任务
			submitFuture := future.SubmitFutureTimeout(time.Second*6, Callable{
				TraceId: fmt.Sprint(fi),
				ActionFunc: func(ctx context.Context) (interface{}, error) {
					time.Sleep(time.Millisecond * time.Duration(rand.Intn(900)))
					sec := time.Now().Second()
					if ctx != nil {
						tid := ctx.Value(TraceID)
						log.Println("[" + tid.(string) + "] run task")
					}
					if sec%2 == 0 {
						return VO{id: fi, name: fmt.Sprint("test", sec)}, nil
					}
					return nil, errors.New(fmt.Sprint(fi, " random exception"))
				},
			})

			res = append(res, submitFuture)
		}

		// 收集结果
		d := sync.WaitGroup{}
		for _, re := range res {
			f := re
			d.Add(1)
			go func() {
				// 不限时阻塞等待结果
				//data, err := f.GetData()

				// 限时阻塞等待结果，但是到期后任务的协程还在执行
				data, err := f.GetDataTimeout(time.Millisecond * 600)
				if err != nil {
					log.Println("["+f.TraceId+"]", "future get error: ", err)
				} else {
					log.Println("["+f.TraceId+"]", "future get", data)
				}
				defer func() {
					d.Done()
				}()
			}()
		}
		d.Wait()
		log.Println("finish all task")
		os.Exit(1)
	}()

	go future.ExecFuturePool(context.Background())
	http.ListenAndServe(":9090", nil)
}

func TestFutureGetWithCancel(t *testing.T) {
	type VO struct {
		id   int
		name string
	}
	future, _ := NewFuturePool(PoolOption{Size: 6})

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
			}, func(ctx context.Context, data interface{}) {
				log.Println("su call", data)
			}, func(ctx context.Context, ex error) {
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

// 限时完成一批任务
func TestNewTmpWithFuture(t *testing.T) {
	log.Println("start")
	//future, _ := NewTmpWithFuture(30, time.Second*4)
	future, _ := NewTmpFuturePool(PoolOption{Size: 30, Timeout: time.Second * 7})
	//time.Sleep(time.Second * 5)
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
	log.Println("finish")
}

// 尝试找到内存泄漏的原因，但是没有
func TestNewTmpWithFutureLeak(t *testing.T) {
	go func() {
		// 访问 http://ip:8899/debug/pprof/
		http.ListenAndServe("0.0.0.0:8899", nil)
	}()

	time.Sleep(time.Second * 10)
	for i := 0; i < 200; i++ {
		log.Println("start")

		var tasks []*FutureTask
		future, _ := NewTmpFuturePool(PoolOption{Size: 30, Timeout: time.Second * 10})
		//time.Sleep(time.Second * 5)
		for i := 0; i < 6; i++ {
			fi := i
			task := future.SubmitFutureTimeout(time.Second*5, Callable{
				fmt.Sprint(fi),
				func(ctx context.Context) (interface{}, error) {
					value := ctx.Value(TraceID)
					//log.Println(value, "start")
					sl := rand.Intn(200) + 20
					time.Sleep(time.Millisecond * time.Duration(sl))
					log.Println(value, "finish", sl, "s")
					return fi, nil
				}, nil, nil,
			})
			tasks = append(tasks, task)
		}
		future.Wait()
		log.Println("finish all", i)
		for _, t := range tasks {
			data, err := t.GetData()
			if err == nil {
				fmt.Print(data)
			}
		}
		fmt.Println()
		future.Close()
	}

	time.Sleep(time.Second * 30)
	runtime.GC()

	time.Sleep(time.Minute * 10)
}

func TestLongWait(t *testing.T) {
	pool, err := NewFuturePool(PoolOption{Size: 3})
	if err != nil {
		return
	}

	for i := 0; i < 8; i++ {
		future := pool.SubmitFuture(Callable{ActionFunc: func(ctx context.Context) (interface{}, error) {
			rsp, err2 := http.Get("http://localhost:9911/longrt")
			if err2 != nil {
				return nil, err2
			}
			all, err2 := io.ReadAll(rsp.Body)
			return all, err2
		}})
		log.Println("submit")
		go func() {
			dataTimeout, err := future.GetDataTimeout(time.Second * 5)
			if err != nil {
				log.Println(err)
				return
			}
			log.Println(string(dataTimeout.([]byte)))
		}()
	}
	pool.Wait()
}
