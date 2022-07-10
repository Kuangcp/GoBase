// Based upon sync.WaitGroup, SizedWaitGroup allows to start multiple
// routines and to wait for their end using the simple API.

// SizedWaitGroup adds the feature of limiting the maximum number of
// concurrently started routines. It could for example be used to start
// multiples routines querying a database but without sending too much
// queries in order to not overload the given database.
//
// Rémy Mathieu © 2016
package sizedwaitgroup

import (
	"context"
	"fmt"
	"sync"
)

// SizedWaitGroup has the same role and close to the
// same API as the Golang sync.WaitGroup but adds a limit of
// the amount of goroutines started concurrently.
type SizedWaitGroup struct {
	Size    int
	Name    string
	current chan struct{}
	wg      sync.WaitGroup
	queue   chan func()
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
		Size:    limit,
		Name:    name,
		current: make(chan struct{}, limit),
		queue:   make(chan func()),
		wg:      sync.WaitGroup{},
	}, nil
}

func NewWithQueue(limit int) (*SizedWaitGroup, error) {
	group, err := NewWithName(limit, "")
	go func() {
		for task := range group.queue {
			action := task
			group.Add()
			go func() {
				defer group.Done()
				action()
			}()
		}
	}()
	return group, err
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
	go func() {
		s.queue <- action
	}()
}
