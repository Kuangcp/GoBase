package sizedpool

import (
	"errors"
	"time"
)

type Future struct {
	Callable
	data       any // 由于 method 无法使用泛型，需要使用在结构体上会导致兼容和使用的复杂性，只能不使用泛型
	ex         error
	finish     chan struct{}
	finishFlag bool
	timeout    time.Duration
}

func NewFuture() *Future {
	return &Future{}
}

func (f *Future) SetData(data interface{}, ex error) {
	f.data = data
	f.ex = ex
	f.finish <- struct{}{}
}

func (f *Future) GetData() (interface{}, error) {
	if f.finishFlag {
		return f.data, f.ex
	}

	select {
	case <-f.finish:
		f.finishFlag = true
		return f.data, f.ex
	}
}
func (f *Future) GetDataTimeout(timeout time.Duration) (interface{}, error) {
	if f.finishFlag {
		return f.data, f.ex
	}

	select {
	case <-f.finish:
		return f.data, f.ex
	case <-time.After(timeout):
		return nil, errors.New("timeout")
	}
}
