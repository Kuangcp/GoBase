package sizedpool

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// SizedWaitGroup has the same role and close to the
// same API as the Golang sync.WaitGroup but adds a size of
// the amount of goroutines started concurrently.
type SizedWaitGroup struct {
	Size        int
	Name        string
	current     chan struct{}
	wg          sync.WaitGroup
	queue       chan func()
	futureQueue chan *FutureTask
	tmpAbort    bool
}

type PoolOption struct {
	Size    int
	Name    string
	Timeout time.Duration
}

// New creates a SizedWaitGroup. The most flexible way
// The size parameter is the maximum amount of
// goroutines which can be started concurrently.
func New(option PoolOption) (*SizedWaitGroup, error) {
	if option.Size <= 0 {
		return nil, fmt.Errorf("size must great than 0")
	}

	return &SizedWaitGroup{
		Size:        option.Size,
		Name:        option.Name,
		current:     make(chan struct{}, option.Size),
		queue:       make(chan func()),
		futureQueue: make(chan *FutureTask),
		wg:          sync.WaitGroup{},
	}, nil
}

func NewWithName(limit int, name string) (SizedWait, error) {
	return New(PoolOption{Size: limit, Name: name})
}

func (s *SizedWaitGroup) Close() {
	close(s.queue)
	close(s.futureQueue)
	close(s.current)
}

func (s *SizedWaitGroup) GetName() string {
	return s.Name
}

func (s *SizedWaitGroup) GetSize() int {
	return s.Size
}

// Add increments the internal WaitGroup counter.
// It can be blocking if the size of spawned goroutines
// has been reached. It will stop blocking when Done has been called.
//
// See sync.WaitGroup documentation for more information.
func (s *SizedWaitGroup) Add() error {
	return s.AddWithContext(context.Background())
}

// AddWithContext increments the internal WaitGroup counter.
// It can be blocking if the size of spawned goroutines
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
	Go(func() {
		defer s.Done()
		action()
	})
}
