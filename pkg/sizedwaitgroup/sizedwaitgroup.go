// Based upon sync.WaitGroup, SizedWaitGroup allows to start multiple
// routines and to wait for their end using the simple API.

// Package sizedwaitgroup SizedWaitGroup adds the feature of limiting the maximum number of
// concurrently started routines. It could for example be used to start
// multiples routines querying a database but without sending too much
// queries in order to not overload the given database.
//
// Rémy Mathieu © 2016
package sizedwaitgroup

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"
)

// SizedWaitGroup has the same role and close to the
// same API as the Golang sync.WaitGroup but adds a limit of
// the amount of goroutines started concurrently.
type SizedWaitGroup struct {
	Size        int
	Name        string
	current     chan struct{}
	wg          sync.WaitGroup
	queue       chan func()
	futureQueue chan *Future
}

// New creates a SizedWaitGroup.
// The limit parameter is the maximum amount of
// goroutines which can be started concurrently.
func New(limit int) (*SizedWaitGroup, error) {
	return NewWithName(limit, "")
}

func NewWithName(limit int, name string) (*SizedWaitGroup, error) {
	if limit <= 0 {
		return nil, fmt.Errorf("limit must great than 0")
	}

	return &SizedWaitGroup{
		Size:        limit,
		Name:        name,
		current:     make(chan struct{}, limit),
		queue:       make(chan func()),
		futureQueue: make(chan *Future),
		wg:          sync.WaitGroup{},
	}, nil
}

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

// Add increments the internal WaitGroup counter.
// It can be blocking if the limit of spawned goroutines
// has been reached. It will stop blocking when Done is
// been called.
//
// See sync.WaitGroup documentation for more information.
func (s *SizedWaitGroup) Add() error {
	return s.AddWithContext(context.Background())
}

// AddWithContext increments the internal WaitGroup counter.
// It can be blocking if the limit of spawned goroutines
// has been reached. It will stop blocking when Done is
// been called, or when the context is canceled. Returns nil on
// success or an error if the context is canceled before the lock
// is acquired.
//
// See sync.WaitGroup documentation for more information.
func (s *SizedWaitGroup) AddWithContext(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case s.current <- struct{}{}:
		break
	}
	s.wg.Add(1)
	return nil
}

// Done decrements the SizedWaitGroup counter.
// See sync.WaitGroup documentation for more information.
func (s *SizedWaitGroup) Done() {
	<-s.current
	s.wg.Done()
}

// Wait blocks until the SizedWaitGroup counter is zero.
// See sync.WaitGroup documentation for more information.
func (s *SizedWaitGroup) Wait() {
	s.wg.Wait()
}

// Run one function, around with Add and Done
func (s *SizedWaitGroup) Run(action func()) {
	err := s.Add()
	if err != nil {
		return
	}
	go func() {
		defer s.Done()
		action()
	}()
}

func (s *SizedWaitGroup) Submit(action func()) {
	s.queue <- action
}

func (s *SizedWaitGroup) SubmitFuture(action func() (interface{}, error),
	success func(data interface{}), failed func(ex error)) *Future {
	return s.SubmitFutureTimeout(time.Duration(0), action, success, failed)
}

func (s *SizedWaitGroup) SubmitFutureTimeout(timeout time.Duration, action func() (interface{}, error),
	success func(data interface{}), failed func(ex error)) *Future {
	future := &Future{
		Action:      action,
		timeout:     timeout,
		SuccessFunc: success,
		FailedFunc:  failed,
		finish:      make(chan struct{}, 1),
	}
	s.futureQueue <- future
	return future
}
