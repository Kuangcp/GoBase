package sizedpool

import (
	"context"
	"fmt"
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
	tmpAbort    bool
}

type PoolOption struct {
	limit   int
	name    string
	timeout time.Duration
}

// New creates a SizedWaitGroup.
// The limit parameter is the maximum amount of
// goroutines which can be started concurrently.
func New(option PoolOption) (*SizedWaitGroup, error) {
	if option.limit <= 0 {
		return nil, fmt.Errorf("limit must great than 0")
	}

	return &SizedWaitGroup{
		Size:        option.limit,
		Name:        option.name,
		current:     make(chan struct{}, option.limit),
		queue:       make(chan func()),
		futureQueue: make(chan *Future),
		wg:          sync.WaitGroup{},
	}, nil
}

func NewWithName(limit int, name string) (*SizedWaitGroup, error) {
	return New(PoolOption{limit: limit, name: name})
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
