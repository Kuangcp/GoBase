package spool

import (
	"context"
	"time"
)

type (
	// Future 参考 java.util.concurrent.Future
	Runnable func(context.Context)
	Callable func(context.Context) (any, error)

	Future interface {
		IsDone() bool
		IsCancelled() bool
		Cancel(mayInterruptIfRunning bool)
		Get() any
		GetTimeout(duration time.Duration) any
	}
	// Executor 参考 java.util.concurrent.ExecutorService
	Executor interface {
		Submit(labels map[string]string, runnable Runnable) Future

		SubmitR(labels map[string]string, runnable Runnable, result any) Future

		SubmitC(labels map[string]string, runnable Callable) Future
	}
)
