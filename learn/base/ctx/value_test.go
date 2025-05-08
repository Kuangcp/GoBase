package main

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestCopy(t *testing.T) {
	f := func(ctx context.Context, k string) {
		if v := ctx.Value(k); v != nil {
			fmt.Println("found value:", v)
			return
		}
		fmt.Println("key not found:", k)
	}

	k := "language"
	ctx := context.WithValue(context.Background(), k, "Go")

	f(ctx, k)
	f(ctx, "color")
}
func TestCopyMore(t *testing.T) {
	f := func(ctx context.Context, k string) {
		if v := ctx.Value(k); v != nil {
			fmt.Println("found value:", v)
			return
		}
		fmt.Println("key not found:", k)
	}

	k := "language"
	// 多次调用 WithValue，因为context是不可变对象，每次都是复制和创建
	ctx := context.WithValue(context.Background(), k, "Go")
	ctx = context.WithValue(ctx, "two", "2")

	f(ctx, k)
	f(ctx, "color")
	f(ctx, "two")
}
func TestCopyMore2(t *testing.T) {
	// ContextValue is a context key
	type ContextValue map[string]interface{}

	f := func(ctx context.Context, k string) {
		data := ctx.Value("data")
		if data == nil {
			fmt.Println("no context data", k)
			return
		}

		if v := data.(ContextValue)[k]; v != nil {
			fmt.Println("found value:", v)
			return
		}

		fmt.Println("key not found:", k)
	}

	data := ContextValue{
		"1": "one",
		"2": "two",
	}
	// 但是会暴露修改能力，有并发安全问题
	data["x"] = "ssss"
	// 传递
	ctx := context.WithValue(context.Background(), "data", data)

	f(ctx, "data")
	f(ctx, "1")
	f(ctx, "2")
}

// 封装为不可变对象，规避并发问题
type Values struct {
	m map[string]interface{}
}

func (v Values) Get(key string) interface{} {
	return v.m[key]
}

func TestCopyMore3(t *testing.T) {
	f := func(ctx context.Context, k string) {
		data := ctx.Value("data")
		if data == nil {
			fmt.Println("no context data", k)
			return
		}

		if v := data.(Values).m[k]; v != nil {
			fmt.Println("found value:", v)
			return
		}

		fmt.Println("key not found:", k)
	}

	data := Values{map[string]interface{}{
		"1": "one",
		"2": "two",
	}}
	// 传递
	ctx := context.WithValue(context.Background(), "data", data)

	f(ctx, "data")
	f(ctx, "1")
	f(ctx, "2")
}

func TestRoutineCopy(t *testing.T) {
	f := func(ctx context.Context, k string) {
		if v := ctx.Value(k); v != nil {
			fmt.Println("found value:", v)
			return
		}
		fmt.Println("key not found:", k)
	}

	ctx := context.WithValue(context.TODO(), "lang", "Go")

	go f(ctx, "lang")
	f(ctx, "color")
	time.Sleep(1 * time.Second)
}
