package sizedpool

import (
	"context"
	"log"
)

func NewWithQueue(limit int) (*SizedWaitGroup, error) {
	group, err := NewWithName(limit, "")
	go group.QueuePool()
	return group, err
}

func (s *SizedWaitGroup) QueuePool() {
	for task := range s.queue {
		action := task
		s.Add()
		go func() {
			defer s.Done()
			action()
		}()
	}
}

func NewWithFuture(limit int) (*SizedWaitGroup, error) {
	group, err := New(limit)
	go group.FuturePool()
	return group, err
}

func (s *SizedWaitGroup) FuturePool() {
	for task := range s.futureQueue {
		future := task
		s.Add()
		if future.timeout.Nanoseconds() == 0 {
			go func() {
				defer s.Done()
				s.finishAction(future)
			}()
		} else {
			// run action func with timeout
			timeout, cancelFunc := context.WithTimeout(context.TODO(), future.timeout)
			go func() {
				go func(ctx context.Context) {
					defer cancelFunc()
					defer s.Done()

					s.finishAction(future)
				}(timeout)
				select {
				case <-timeout.Done():
					if timeout.Err().Error() == "context deadline exceeded" {
						log.Println("timeout")
					}
					return
				}
			}()
		}
	}
}

func (s *SizedWaitGroup) finishAction(future *Future) {
	data, actionErr := future.Action()
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

// TODO 具有超时取消功能的协程池
