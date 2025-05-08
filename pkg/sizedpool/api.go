package sizedpool

import (
	"context"
	"time"
)

const (
	TraceID = "traceId"
)

type (
	Callable struct {
		TraceId string
		// ctx.Value(TraceID) 获取 traceId
		ActionFunc func(ctx context.Context) (interface{}, error)
		// 使用Future方式时，以下两个函数不声明
		SuccessFunc func(ctx context.Context, data interface{})
		FailedFunc  func(ctx context.Context, ex error)
	}

	// QueuePool async submit task. then run with poll queue
	// NewQueuePool
	QueuePool interface {
		SizedWait
		Submit(action func())
	}

	// FuturePool async submit task. then run with poll queue. support future get data and size timeout
	// NewTmpFuturePool & NewFuturePool
	FuturePool interface {
		SizedWait
		ExecTmpFuturePool(timeout time.Duration)
		ExecFuturePool(ctx context.Context)

		SubmitFuture(callable Callable) *FutureTask
		SubmitFutureTimeout(timeout time.Duration, callable Callable) *FutureTask
	}

	// SizedWait only sized wait size
	// NewWithName
	SizedWait interface {
		GetName() string
		GetSize() int

		AddWithContext(ctx context.Context) error
		Add() error
		Done()
		Wait()

		Run(action func())

		Close()
	}
)
