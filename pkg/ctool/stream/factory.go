package stream

import "fmt"

type (

	// FilterFunc defines the method to filter a Stream.
	FilterFunc func(item any) bool
	// ForAllFunc defines the method to handle all elements in a Stream.
	ForAllFunc func(pipe <-chan any)
	// ForEachFunc defines the method to handle each element in a Stream.
	ForEachFunc func(item any)
	// GenerateFunc defines the method to send elements into a Stream.
	GenerateFunc func(source chan<- any)
	// KeyFunc defines the method to generate keys for the elements in a Stream.
	KeyFunc func(item any) any
	// LessFunc defines the method to compare the elements in a Stream.
	LessFunc func(a, b any) bool
	// MapFunc defines the method to map each element to another object in a Stream.
	MapFunc func(item any) any
	// Option defines the method to customize a Stream.
	Option func(opts *RxOptions)
	// ParallelFunc defines the method to handle elements parallel.
	ParallelFunc func(item any)
	// ReduceFunc defines the method to reduce all the elements in a Stream.
	ReduceFunc func(pipe <-chan any) (any, error)
	// WalkFunc defines the method to walk through all the elements in a Stream.
	WalkFunc func(item any, pipe chan<- any)
)

// Concat returns a concatenated Stream.
func Concat(s Stream, others ...Stream) Stream {
	return s.Concat(others...)
}

// From constructs a Stream from the given GenerateFunc.
func From(generate GenerateFunc) Stream {
	source := make(chan any)

	GoSafe(func() {
		defer close(source)
		generate(source)
	})

	return Range(source)
}
func Empty() Stream {
	source := make(chan any)
	close(source)
	return Range(source)
}

// JustN produce 1...n number serial
func JustN(n int) Stream {
	source := make(chan any, n)
	for i := 1; i <= n; i++ {
		source <- i
	}
	close(source)

	return Range(source)
}

// Just converts the given arbitrary items to a Stream.
func Just[T any](items ...T) Stream {
	source := make(chan any, len(items))
	for _, item := range items {
		source <- item
	}
	close(source)

	return Range(source)
}

// Range converts the given channel to a Stream.
func Range(source <-chan any) Stream {
	return Stream{
		source: source,
	}
}

// UnlimitedWorkers lets the caller use as many Workers as the tasks.
func UnlimitedWorkers() Option {
	return func(opts *RxOptions) {
		opts.UnlimitedWorkers = true
	}
}

// WithWorkers lets the caller customize the concurrent Workers.
func WithWorkers(workers int) Option {
	return func(opts *RxOptions) {
		if workers < minWorkers {
			opts.Workers = minWorkers
		} else {
			opts.Workers = workers
		}
	}
}

// buildOptions returns a RxOptions with given customizations.
func buildOptions(opts ...Option) *RxOptions {
	options := newOptions()
	for _, opt := range opts {
		opt(options)
	}

	return options
}

// drain drains the given channel.
func drain(channel <-chan any) {
	for range channel {
	}
}

// newOptions returns a default RxOptions.
func newOptions() *RxOptions {
	return &RxOptions{
		Workers: defaultWorkers,
	}
}

// simplify function

func ToString(item any) any {
	return fmt.Sprint(item)
}

func Self[R any](a any) R {
	return a.(R)
}

func Println(item any) {
	fmt.Println(item)
}
