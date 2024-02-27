package stream

import "github.com/kuangcp/gobase/pkg/ctool"

type Collector[T, A, R any] interface {
	supplier() func() A      // 提供结果对象 容器或单值
	accumulator() func(A, T) // 期望结果类型的值 填入结果对象
	combiner() func(A, A) A  // 组合两个结果对象
	finisher() func(A) R     // 结果对象A转换为R 或者 完成A的构造
}

type CollectorItem[T, A, R any] struct {
	sup  func() A
	acc  func(A, T)
	comb func(A, A) A
	fin  func(A) R
}

func (c *CollectorItem[T, A, R]) supplier() func() A {
	return c.sup
}

func (c *CollectorItem[T, A, R]) accumulator() func(A, T) {
	return c.acc
}

func (c *CollectorItem[T, A, R]) combiner() func(A, A) A {
	return c.comb
}

func (c *CollectorItem[T, A, R]) finisher() func(A) R {
	return c.fin
}

// Collect java.util.stream.Collector
func Collect[T, A, R any](s Stream, c Collector[T, A, R]) R {
	container := c.supplier()()
	s.ForEach(func(item any) {
		c.accumulator()(container, item.(T))
	})

	return c.finisher()(container)
}

func Set[T comparable]() Collector[T, *ctool.Set[T], *ctool.Set[T]] {
	return &CollectorItem[T, *ctool.Set[T], *ctool.Set[T]]{
		sup: func() *ctool.Set[T] {
			return ctool.NewSet[T]()
		},
		acc: func(a *ctool.Set[T], t T) {
			a.Add(t)
		},
		comb: func(c *ctool.Set[T], c2 *ctool.Set[T]) *ctool.Set[T] {
			c.Adds(c2)
			return c
		},
		fin: func(c *ctool.Set[T]) *ctool.Set[T] {
			return c
		},
	}
}
