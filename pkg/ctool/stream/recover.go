package stream

import (
	"context"
	"log"
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
		log.Println(p)
	}
}

// RecoverCtx is used with defer to do cleanup on panics.
func RecoverCtx(ctx context.Context, cleanups ...func()) {
	for _, cleanup := range cleanups {
		cleanup()
	}

	if p := recover(); p != nil {
		//logx.WithContext(ctx).Errorf("%+v\n%s", p, debug.Stack())
		log.Printf("%+v\n%s", p, debug.Stack())
	}
}
