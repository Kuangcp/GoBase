package stream

import (
	"context"
	"github.com/kuangcp/logger"
	"runtime/debug"
)

// Recover is used with defer to do cleanup on panics.
// Use it like:
//
//	defer Recover(func() {})
func Recover(cleanups ...func()) {
	for _, cleanup := range cleanups {
		cleanup()
	}

	if p := recover(); p != nil {
		//logx.ErrorStack(p)
		logger.Error(p)
	}
}

// RecoverCtx is used with defer to do cleanup on panics.
func RecoverCtx(ctx context.Context, cleanups ...func()) {
	for _, cleanup := range cleanups {
		cleanup()
	}

	if p := recover(); p != nil {
		//logx.WithContext(ctx).Errorf("%+v\n%s", p, debug.Stack())
		logger.Error("%+v\n%s", p, debug.Stack())
	}
}
