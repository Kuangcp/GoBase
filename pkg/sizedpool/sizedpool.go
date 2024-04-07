package sizedpool

import (
	"context"
	"errors"
	"log"
	"time"
)

func NewQueuePool(limit int) (QueuePool, error) {
	group, err := New(PoolOption{Size: limit})
	go group.ExecQueuePool()
	return group, err
}

func NewTmpFuturePool(option PoolOption) (FuturePool, error) {
	group, err := New(option)
	if option.Timeout == 0 {
		return nil, errors.New("not init timeout")
	}
	go group.ExecTmpFuturePool(option.Timeout)
	return group, err
}

func NewFuturePool(option PoolOption) (FuturePool, error) {
	group, err := New(option)
	Go(func() {
		group.ExecFuturePool(nil)
	})
	return group, err
}

// ExecQueuePool 执行队列任务
func (s *SizedWaitGroup) ExecQueuePool() {
	for task := range s.queue {
		action := task
		s.Add()
		Go(func() {
			defer s.Done()
			action()
		})
	}
}

// ExecTmpFuturePool 调度执行池中任务
// timeout: 注意：goroutine 无法被中止，所以需要等已运行的协程执行完成后，才能真正退出池 所以超时时间通常大于设置的值
func (s *SizedWaitGroup) ExecTmpFuturePool(timeout time.Duration) {
	timeoutCtx, cancelFunc := context.WithTimeout(context.TODO(), timeout)

	go func(ctx context.Context) {
		defer cancelFunc()
		s.ExecFuturePool(timeoutCtx)
	}(timeoutCtx)

	select {
	case <-timeoutCtx.Done():
		if timeoutCtx.Err().Error() == "context deadline exceeded" {
			log.Println("total timeout")
			s.tmpAbort = true
		}
		return
	}
}

// ExecFuturePool 调度执行池中任务
// ctx: 限制future执行时间，空表示不限制
func (s *SizedWaitGroup) ExecFuturePool(ctx context.Context) {
	for task := range s.futureQueue {
		future := task

		//log.Println("add task", task,s.tmpAbort)
		if s.tmpAbort {
			log.Println(future.TraceId, "WARN: timeout, task reject.")
			continue
		}

		s.Add()

		if ctx != nil {
			ctx = context.WithValue(ctx, TraceID, future.TraceId)
		}
		if future.timeout.Nanoseconds() == 0 {
			go func(_ context.Context) {
				defer s.Done()
				s.execAction(ctx, future)
			}(ctx)
		} else {
			// run action func with timeout
			go func(ctx context.Context) {
				timeout, cancelFunc := context.WithTimeout(context.TODO(), future.timeout)
				go func(_ context.Context) {
					defer cancelFunc()
					defer s.Done()

					s.execAction(ctx, future)
				}(timeout)
				select {
				case <-timeout.Done():
					if timeout.Err().Error() == "context deadline exceeded" {
						log.Println(future.TraceId, "future timeout")
					}
					return
				}
			}(ctx)
		}
	}
}

func (s *SizedWaitGroup) execAction(ctx context.Context, future *FutureTask) {
	data, actionErr := future.ActionFunc(ctx)
	future.SetData(data, actionErr)

	if actionErr != nil {
		if future.FailedFunc != nil {
			future.FailedFunc(actionErr)
		}
	} else {
		if future.SuccessFunc != nil {
			future.SuccessFunc(data)
		}
	}
}

func (s *SizedWaitGroup) Submit(action func()) {
	s.queue <- action
}

func (s *SizedWaitGroup) SubmitFuture(callable Callable) *FutureTask {
	return s.SubmitFutureTimeout(time.Duration(0), callable)
}

func (s *SizedWaitGroup) SubmitFutureTimeout(timeout time.Duration, callable Callable) *FutureTask {
	if s.tmpAbort {
		return nil
	}
	future := &FutureTask{
		timeout: timeout,
		finish:  make(chan struct{}, 1),
	}
	future.Callable = callable
	s.futureQueue <- future
	return future
}
