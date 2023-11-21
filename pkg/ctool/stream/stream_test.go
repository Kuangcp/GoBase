package stream

import (
	"fmt"
	"github.com/kuangcp/gobase/pkg/ctool"
	"go.uber.org/goleak"
	"io"
	"log"
	"math/rand"
	"reflect"
	"runtime"
	"sort"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBuffer(t *testing.T) {
	runCheckedTest(t, func(t *testing.T) {
		const N = 5
		var count int32
		var wait sync.WaitGroup
		wait.Add(1)
		From(func(source chan<- any) {
			ticker := time.NewTicker(10 * time.Millisecond)
			defer ticker.Stop()

			for i := 0; i < 2*N; i++ {
				select {
				case source <- i:
					atomic.AddInt32(&count, 1)
				case <-ticker.C:
					wait.Done()
					return
				}
			}
		}).Buffer(N).ForAll(func(pipe <-chan any) {
			wait.Wait()
			// why N+1, because take one more to wait for sending into the channel
			assert.Equal(t, int32(N+1), atomic.LoadInt32(&count))
		})
	})
}

func TestBufferNegative(t *testing.T) {
	runCheckedTest(t, func(t *testing.T) {
		var result int
		Just(1, 2, 3, 4).Buffer(-1).Reduce(func(pipe <-chan any) (any, error) {
			for item := range pipe {
				result += item.(int)
			}
			return result, nil
		})
		assert.Equal(t, 10, result)
	})
}

func TestCount(t *testing.T) {
	runCheckedTest(t, func(t *testing.T) {
		tests := []struct {
			name     string
			elements []any
		}{
			{
				name: "no elements with nil",
			},
			{
				name:     "no elements",
				elements: []any{},
			},
			{
				name:     "1 element",
				elements: []any{1},
			},
			{
				name:     "multiple elements",
				elements: []any{1, 2, 3},
			},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				val := Just(test.elements...).Count()
				assert.Equal(t, len(test.elements), val)
			})
		}
	})
}

func TestDone(t *testing.T) {
	runCheckedTest(t, func(t *testing.T) {
		var count int32
		Just(1, 2, 3).Walk(func(item any, pipe chan<- any) {
			time.Sleep(time.Millisecond * 100)
			atomic.AddInt32(&count, int32(item.(int)))
		}).Done()
		assert.Equal(t, int32(6), count)
	})
}

func TestJust(t *testing.T) {
	runCheckedTest(t, func(t *testing.T) {
		var result int
		Just(1, 2, 3, 4).Reduce(func(pipe <-chan any) (any, error) {
			for item := range pipe {
				result += item.(int)
			}
			return result, nil
		})
		assert.Equal(t, 10, result)
	})
}
func TestJustT(t *testing.T) {
	var s = []string{"x", "2"}
	Just(s).ForEach(func(item any) {
		fmt.Println(item)
	})
}

func TestDistinct(t *testing.T) {
	runCheckedTest(t, func(t *testing.T) {
		var result int
		Just(4, 1, 3, 2, 3, 4).Distinct(func(item any) any {
			return item
		}).Reduce(func(pipe <-chan any) (any, error) {
			for item := range pipe {
				result += item.(int)
			}
			return result, nil
		})
		assert.Equal(t, 10, result)
	})
}

func TestFilter(t *testing.T) {
	runCheckedTest(t, func(t *testing.T) {
		var result int
		Just(1, 2, 3, 4).
			Filter(func(item any) bool {
				return item.(int)%2 == 0
			}).Reduce(func(pipe <-chan any) (any, error) {
			for item := range pipe {
				result += item.(int)
			}
			return result, nil
		})
		assert.Equal(t, 6, result)
	})
}

func TestFirst(t *testing.T) {
	runCheckedTest(t, func(t *testing.T) {
		assert.Nil(t, Just[int]().First())
		assert.Equal(t, "foo", Just("foo").First())
		assert.Equal(t, "foo", Just("foo", "bar").First())
	})
}

func TestForAll(t *testing.T) {
	runCheckedTest(t, func(t *testing.T) {
		var result int
		Just(1, 2, 3, 4).Filter(func(item any) bool {
			return item.(int)%2 == 0
		}).ForAll(func(pipe <-chan any) {
			for item := range pipe {
				result += item.(int)
			}
		})
		assert.Equal(t, 6, result)
	})
}

func TestStream_ForEach(t *testing.T) {
	stream := JustN(10).MapStr()
	stream.ForEach(func(item any) {
		fmt.Println(item)
	})
	stream.Map(func(item any) any {
		return item.(string) + "x"
	}).ForEach(Println)

}

func TestStream_ForEachNone(t *testing.T) {
	JustN(2).Filter(func(item any) bool {
		return item.(int) > 3
	}).ForEach(func(item any) {
		fmt.Println(item)
	})

	JustN(2).Filter(func(item any) bool {
		return item.(int) > 3
	}).Map(func(item any) any {
		fmt.Println(item)
		return 1
	})

	str := JustN(2).Filter(func(item any) bool {
		return item.(int) > 3
	}).MapStr()
	fmt.Println("result:[" + ToJoins(str, "m") + "]")
}

func TestGroup(t *testing.T) {
	runCheckedTest(t, func(t *testing.T) {
		var groups [][]int
		Just(10, 11, 20, 21).Group(func(item any) any {
			v := item.(int)
			return v / 10
		}).ForEach(func(item any) {
			v := item.(GroupItem)
			var group []int
			for _, each := range v.Val {
				group = append(group, each.(int))
			}
			groups = append(groups, group)
		})

		assert.Equal(t, 2, len(groups))
		for _, group := range groups {
			assert.Equal(t, 2, len(group))
			assert.True(t, group[0]/10 == group[1]/10)
		}
	})
}

func TestStream_GroupParallel(t *testing.T) {
	start := time.Now().UnixMicro()
	total := 300
	JustN(total).Group(func(item any) any {
		v := item.(int)
		time.Sleep(time.Microsecond * 1)
		return v / 3
	}, func(opts *RxOptions) {
		opts.Workers = 10
		//opts.UnlimitedWorkers = true
	}).ForEach(func(item any) {
		//v := item.(GroupItem)
		//fmt.Println(v)
	})
	fmt.Println("parallel ----", time.Now().UnixMicro()-start, "us")

	// 如果数据量小或代码执行成本很低，开并发后锁竞争远大于代码执行，反而会导致耗时的增加
	// 按Group的使用场景来说，绝大多数场景不需要开并发,如果有io阻塞类代码则推荐使用

	start = time.Now().UnixMicro()
	JustN(total).Group(func(item any) any {
		v := item.(int)
		time.Sleep(time.Microsecond * 1)
		return v / 3
	}).ForEach(func(item any) {
		//v := item.(GroupItem)
		//fmt.Println(v)
	})
	fmt.Println("serial ------", time.Now().UnixMicro()-start, "us")
}

type User struct {
	id     int
	name   string
	areaId int
}

func TestGroupConstruct(t *testing.T) {
	Just(1, 8, 10, 11, 20, 21).Map(func(item any) any {
		v := item.(int)
		return User{
			id:     v,
			name:   fmt.Sprint(v),
			areaId: v / 3,
		}
	}).Group(func(item any) any {
		v := item.(User)
		return v.areaId
	}).ForEach(func(item any) {
		l := item.(GroupItem)
		for _, i := range l.Val {
			u := i.(User)
			fmt.Println("area:", l.Key, " user:", u.id, u.name)
		}
	})
}

func TestHead(t *testing.T) {
	runCheckedTest(t, func(t *testing.T) {
		var result int
		Just(1, 2, 3, 4).Head(2).Reduce(func(pipe <-chan any) (any, error) {
			for item := range pipe {
				result += item.(int)
			}
			return result, nil
		})
		assert.Equal(t, 3, result)
	})
}

func TestHeadZero(t *testing.T) {
	runCheckedTest(t, func(t *testing.T) {
		assert.Panics(t, func() {
			Just(1, 2, 3, 4).Head(0).Reduce(func(pipe <-chan any) (any, error) {
				return nil, nil
			})
		})
	})
}

func TestHeadMore(t *testing.T) {
	runCheckedTest(t, func(t *testing.T) {
		var result int
		Just(1, 2, 3, 4).Head(6).Reduce(func(pipe <-chan any) (any, error) {
			for item := range pipe {
				result += item.(int)
			}
			return result, nil
		})
		assert.Equal(t, 10, result)
	})
}

func TestLast(t *testing.T) {
	runCheckedTest(t, func(t *testing.T) {
		goroutines := runtime.NumGoroutine()
		assert.Nil(t, Just[int]().Last())
		assert.Equal(t, "foo", Just("foo").Last())
		assert.Equal(t, "bar", Just("foo", "bar").Last())
		// let scheduler schedule first
		runtime.Gosched()
		assert.Equal(t, goroutines, runtime.NumGoroutine())
	})
}

func TestMap(t *testing.T) {
	runCheckedTest(t, func(t *testing.T) {
		log.SetOutput(io.Discard)

		tests := []struct {
			mapper MapFunc
			expect int
		}{
			{
				mapper: func(item any) any {
					v := item.(int)
					return v * v
				},
				expect: 30,
			},
			{
				mapper: func(item any) any {
					v := item.(int)
					if v%2 == 0 {
						return 0
					}
					return v * v
				},
				expect: 10,
			},
			{
				mapper: func(item any) any {
					v := item.(int)
					if v%2 == 0 {
						panic(v)
					}
					return v * v
				},
				expect: 10,
			},
		}

		// Map(...) works even WithWorkers(0)
		for i, test := range tests {
			t.Run(ctool.RandomAlpha(5), func(t *testing.T) {
				var result int
				var workers int
				if i%2 == 0 {
					workers = 0
				} else {
					workers = runtime.NumCPU()
				}
				From(func(source chan<- any) {
					for i := 1; i < 5; i++ {
						source <- i
					}
				}).Map(test.mapper, WithWorkers(workers)).Reduce(
					func(pipe <-chan any) (any, error) {
						for item := range pipe {
							result += item.(int)
						}
						return result, nil
					})

				assert.Equal(t, test.expect, result)
			})
		}
	})
}

func TestMerge(t *testing.T) {
	runCheckedTest(t, func(t *testing.T) {
		Just(1, 2, 3, 4).Merge().ForEach(func(item any) {
			assert.ElementsMatch(t, []any{1, 2, 3, 4}, item.([]any))
		})
	})
}

func TestParallelJust(t *testing.T) {
	runCheckedTest(t, func(t *testing.T) {
		var count int32
		Just(1, 2, 3).Parallel(func(item any) {
			time.Sleep(time.Millisecond * 100)
			atomic.AddInt32(&count, int32(item.(int)))
		}, UnlimitedWorkers())
		assert.Equal(t, int32(6), count)
	})
}

func TestReverse(t *testing.T) {
	runCheckedTest(t, func(t *testing.T) {
		Just(1, 2, 3, 4).Reverse().Merge().ForEach(func(item any) {
			assert.ElementsMatch(t, []any{4, 3, 2, 1}, item.([]any))
		})
	})
}

func TestSort(t *testing.T) {
	runCheckedTest(t, func(t *testing.T) {
		var prev int
		Just(5, 3, 7, 1, 9, 6, 4, 8, 2).Sort(func(a, b any) bool {
			return a.(int) < b.(int)
		}).ForEach(func(item any) {
			next := item.(int)
			assert.True(t, prev < next)
			prev = next
		})
	})
}

func TestSortWithOriginMemory(t *testing.T) {
	var s = []int{1, 4, 5, 2, 11, 2, 4, 5}
	Just(s...).Sort(func(a, b any) bool {
		return a.(int) < b.(int)
	})
	fmt.Println(s)
	rs := Just(s...).Sort(func(a, b any) bool {
		return a.(int) < b.(int)
	})
	ss := ToList[int](rs)
	fmt.Println(s)
	fmt.Println(ss)
}

func TestSplit(t *testing.T) {
	runCheckedTest(t, func(t *testing.T) {
		assert.Panics(t, func() {
			Just(1, 2, 3, 4, 5, 6, 7, 8, 9, 10).Split(0).Done()
		})
		var chunks [][]any
		Just(1, 2, 3, 4, 5, 6, 7, 8, 9, 10).Split(4).ForEach(func(item any) {
			chunk := item.([]any)
			chunks = append(chunks, chunk)
		})
		assert.EqualValues(t, [][]any{
			{1, 2, 3, 4},
			{5, 6, 7, 8},
			{9, 10},
		}, chunks)
	})
}

func TestTail(t *testing.T) {
	runCheckedTest(t, func(t *testing.T) {
		var result int
		Just(1, 2, 3, 4).Tail(2).Reduce(func(pipe <-chan any) (any, error) {
			for item := range pipe {
				result += item.(int)
			}
			return result, nil
		})
		assert.Equal(t, 7, result)
	})
}

func TestTailZero(t *testing.T) {
	runCheckedTest(t, func(t *testing.T) {
		assert.Panics(t, func() {
			Just(1, 2, 3, 4).Tail(0).Reduce(func(pipe <-chan any) (any, error) {
				return nil, nil
			})
		})
	})
}

func TestWalk(t *testing.T) {
	runCheckedTest(t, func(t *testing.T) {
		var result int
		Just(1, 2, 3, 4, 5).Walk(func(item any, pipe chan<- any) {
			if item.(int)%2 != 0 {
				pipe <- item
			}
		}, UnlimitedWorkers()).ForEach(func(item any) {
			result += item.(int)
		})
		assert.Equal(t, 9, result)
	})
}

func TestStream_AnyMach(t *testing.T) {
	runCheckedTest(t, func(t *testing.T) {
		assetEqual(t, false, Just(1, 2, 3).AnyMatch(func(item any) bool {
			return item.(int) == 4
		}))
		assetEqual(t, false, Just(1, 2, 3).AnyMatch(func(item any) bool {
			return item.(int) == 0
		}))
		assetEqual(t, true, Just(1, 2, 3).AnyMatch(func(item any) bool {
			return item.(int) == 2
		}))
		assetEqual(t, true, Just(1, 2, 3).AnyMatch(func(item any) bool {
			return item.(int) == 2
		}))
	})
}

func TestStream_AllMach(t *testing.T) {
	runCheckedTest(t, func(t *testing.T) {
		assetEqual(
			t, true, Just(1, 2, 3).AllMatch(func(item any) bool {
				return true
			}),
		)
		assetEqual(
			t, false, Just(1, 2, 3).AllMatch(func(item any) bool {
				return false
			}),
		)
		assetEqual(
			t, false, Just(1, 2, 3).AllMatch(func(item any) bool {
				return item.(int) == 1
			}),
		)
	})
}

func TestStream_NoneMatch(t *testing.T) {
	runCheckedTest(t, func(t *testing.T) {
		assetEqual(
			t, true, Just(1, 2, 3).NoneMatch(func(item any) bool {
				return false
			}),
		)
		assetEqual(
			t, false, Just(1, 2, 3).NoneMatch(func(item any) bool {
				return true
			}),
		)
		assetEqual(
			t, true, Just(1, 2, 3).NoneMatch(func(item any) bool {
				return item.(int) == 4
			}),
		)
	})
}

func TestStream_Flat(t *testing.T) {
	flat := JustN(5).Map(func(item any) any {
		return ctool.RandomAlpha(item.(int))
	}).Flat(func(a any) Stream {
		return Just(a, a, a)
	})
	result := ToList[string](flat)
	fmt.Println(result)
}

func TestStream_FlatEmpty(t *testing.T) {
	flat := JustN(5).Map(func(item any) any {
		return ctool.RandomAlpha(item.(int))
	}).Flat(func(a any) Stream {
		if len(a.(string)) > 2 {
			return Just(a, a, "#")
		} else {
			return Empty()
		}
	})
	result := ToList[string](flat)
	fmt.Println(result)
}

func TestStream_Parallel(t *testing.T) {
	JustN(10).MapStr().Parallel(func(item any) {
		fmt.Println(item)
	}, UnlimitedWorkers())

	s := JustN(10).Map(func(item any) any {
		time.Sleep(time.Second)
		return item
	}, UnlimitedWorkers()).Map(func(item any) any {
		return ctool.RandomAlpha(item.(int))
	}, UnlimitedWorkers())
	fmt.Println(ToJoins(s, ","))
}

func TestConcat(t *testing.T) {
	runCheckedTest(t, func(t *testing.T) {
		a1 := []any{1, 2, 3}
		a2 := []any{4, 5, 6}
		s1 := Just(a1...)
		s2 := Just(a2...)
		stream := Concat(s1, s2)
		var items []any
		for item := range stream.source {
			items = append(items, item)
		}
		sort.Slice(items, func(i, j int) bool {
			return items[i].(int) < items[j].(int)
		})
		ints := make([]any, 0)
		ints = append(ints, a1...)
		ints = append(ints, a2...)
		assetEqual(t, ints, items)
	})
}

func TestStream_Skip(t *testing.T) {
	runCheckedTest(t, func(t *testing.T) {
		assetEqual(t, 3, Just(1, 2, 3, 4).Skip(1).Count())
		assetEqual(t, 1, Just(1, 2, 3, 4).Skip(3).Count())
		assetEqual(t, 4, Just(1, 2, 3, 4).Skip(0).Count())
		equal(t, Just(1, 2, 3, 4).Skip(3), []any{4})
		assert.Panics(t, func() {
			Just(1, 2, 3, 4).Skip(-1)
		})
	})
}

func TestStream_Concat(t *testing.T) {
	runCheckedTest(t, func(t *testing.T) {
		stream := Just(1).Concat(Just(2), Just(3))
		var items []any
		for item := range stream.source {
			items = append(items, item)
		}
		sort.Slice(items, func(i, j int) bool {
			return items[i].(int) < items[j].(int)
		})
		assetEqual(t, []any{1, 2, 3}, items)

		just := Just(1)
		equal(t, just.Concat(just), []any{1})
	})
}

func TestStream_Max(t *testing.T) {
	runCheckedTest(t, func(t *testing.T) {
		tests := []struct {
			name     string
			elements []any
			max      any
		}{
			{
				name: "no elements with nil",
			},
			{
				name:     "no elements",
				elements: []any{},
				max:      nil,
			},
			{
				name:     "1 element",
				elements: []any{1},
				max:      1,
			},
			{
				name:     "multiple elements",
				elements: []any{1, 2, 9, 5, 8},
				max:      9,
			},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				val := Just(test.elements...).Max(func(a, b any) bool {
					return a.(int) < b.(int)
				})
				assetEqual(t, test.max, val)
			})
		}
	})
}

func TestStream_Min(t *testing.T) {
	runCheckedTest(t, func(t *testing.T) {
		tests := []struct {
			name     string
			elements []any
			min      any
		}{
			{
				name: "no elements with nil",
				min:  nil,
			},
			{
				name:     "no elements",
				elements: []any{},
				min:      nil,
			},
			{
				name:     "1 element",
				elements: []any{1},
				min:      1,
			},
			{
				name:     "multiple elements",
				elements: []any{-1, 1, 2, 9, 5, 8},
				min:      -1,
			},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				val := Just(test.elements...).Min(func(a, b any) bool {
					return a.(int) < b.(int)
				})
				assetEqual(t, test.min, val)
			})
		}
	})
}

func TestStream_Fork(t *testing.T) {
	watch := ctool.NewStopWatch()

	watch.Start("fork")
	a, b, c := JustN(30).Map(func(item any) any {
		log.Println("xxx", item)
		time.Sleep(time.Millisecond * 40)
		return item
	}).ForkTri()
	watch.Stop()

	// fork 异步执行, 但是会阻塞后序Stream操作,如果此处没有sleep 第一个操作会最耗时,因为需要消费完全部的上游Stream才能进行下一步操作
	time.Sleep(time.Millisecond * 400)

	watch.Run("fork3", func() {
		maxI := c.Max(func(a, b any) bool {
			return a.(int) < b.(int)
		})
		fmt.Println("max val is", maxI)
	})

	watch.Run("fork1", func() {
		a.Filter(func(item any) bool {
			time.Sleep(time.Millisecond * 3)
			return item.(int) > 7
		}).ForEach(Println)
	})

	watch.Run("for2", func() {
		b.Filter(func(item any) bool {
			time.Sleep(time.Millisecond * 3)
			return item.(int) < 5
		}).ForEach(Println)
	})

	fmt.Println(watch.PrettyPrint())
}

func TestStream_ForkParallel(t *testing.T) {
	watch := ctool.NewStopWatch()

	watch.Start("build stream")
	stream := JustN(10).Map(func(item any) any {
		time.Sleep(time.Millisecond * 50)
		re := rand.Intn(item.(int))
		log.Println("provider", re)
		return re
	})
	watch.Stop()

	var result struct {
		odd  []int
		lit  []int
		max  int
		join string
	}

	watch.Start("parallel")
	stream.ForkParallel(func(b Stream) {
		x := b.Filter(func(item any) bool {
			time.Sleep(time.Millisecond * 200)
			//fmt.Println("fork2", item)
			return item.(int) < 5
		})
		result.lit = ToList[int](x)
	}, func(stream Stream) {
		result.odd = ToList[int](stream.Filter(func(item any) bool {
			return item.(int)%2 == 0
		}))
	}, func(c Stream) {
		maxI := c.Max(func(a, b any) bool {
			//fmt.Println("max", a, b)
			return a.(int) < b.(int)
		})
		result.max = maxI.(int)
		//fmt.Println("max val is", maxI)
	}, func(stream Stream) {
		joins := ToJoins(stream, ".")
		result.join = joins
	})

	watch.Stop()
	fmt.Println(watch.PrettyPrint())
	fmt.Println("RESULT ", result)
}

func BenchmarkParallelMapReduce(b *testing.B) {
	b.ReportAllocs()

	mapper := func(v any) any {
		return v.(int64) * v.(int64)
	}
	reducer := func(input <-chan any) (any, error) {
		var result int64
		for v := range input {
			result += v.(int64)
		}
		return result, nil
	}
	b.ResetTimer()
	From(func(input chan<- any) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				input <- int64(rand.Int())
			}
		})
	}).Map(mapper).Reduce(reducer)
}

func BenchmarkMapReduce(b *testing.B) {
	b.ReportAllocs()

	mapper := func(v any) any {
		return v.(int64) * v.(int64)
	}
	reducer := func(input <-chan any) (any, error) {
		var result int64
		for v := range input {
			result += v.(int64)
		}
		return result, nil
	}
	b.ResetTimer()
	From(func(input chan<- any) {
		for i := 0; i < b.N; i++ {
			input <- int64(rand.Int())
		}
	}).Map(mapper).Reduce(reducer)
}

func TestBackpressure(t *testing.T) {
	From(func(source chan<- any) {
		for i := 0; i < 100; i++ {
			time.Sleep(time.Millisecond * 500)
			fmt.Println("...")
			source <- i
		}
	}).Buffer(6).ForEach(func(item any) {
		time.Sleep(time.Second * 2)
		log.Println("for ", item)
	})
}

func assetEqual(t *testing.T, except, data any) {
	if !reflect.DeepEqual(except, data) {
		t.Errorf(" %v, want %v", data, except)
	}
}

func equal(t *testing.T, stream Stream, data []any) {
	items := make([]any, 0)
	for item := range stream.source {
		items = append(items, item)
	}
	if !reflect.DeepEqual(items, data) {
		t.Errorf(" %v, want %v", items, data)
	}
}

func runCheckedTest(t *testing.T, fn func(t *testing.T)) {
	defer goleak.VerifyNone(t)
	fn(t)
}
