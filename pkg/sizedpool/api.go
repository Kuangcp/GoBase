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
		ActionFunc  func(ctx context.Context) (interface{}, error)
		SuccessFunc func(data interface{})
		FailedFunc  func(ex error)
	}

	FuturePool interface {
		ExecTmpFuturePool(timeout time.Duration)
		ExecFuturePool(ctx context.Context)

		Wait()

		SubmitFuture(callable Callable) *Future
		SubmitFutureTimeout(timeout time.Duration, callable Callable) *Future
	}
)
