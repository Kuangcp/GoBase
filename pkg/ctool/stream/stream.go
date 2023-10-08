package stream

import (
	"sort"
	"sync"
)

const (
	// defaultWorkers here set 1 make stream keep order
	defaultWorkers = 1
	minWorkers     = 1
)

type (
	RxOptions struct {
		UnlimitedWorkers bool
		Workers          int
	}

	// A Stream is a stream that can be used to do stream processing.
	Stream struct {
		source <-chan any
	}

	GroupItem struct {
		Key any
		Val []any
	}
)

// Buffer buffers the items into a queue with size n.
// It can balance the producer and the consumer if their processing throughput don't match.
func (s Stream) Buffer(n int) Stream {
	if n < 0 {
		n = 0
	}

	source := make(chan any, n)
	go func() {
		for item := range s.source {
			source <- item
		}
		close(source)
	}()

	return Range(source)
}

// Concat returns a Stream that concatenated other streams
func (s Stream) Concat(others ...Stream) Stream {
	source := make(chan any)

	go func() {
		group := NewRoutineGroup()
		group.Run(func() {
			for item := range s.source {
				source <- item
			}
		})

		for _, each := range others {
			each := each
			group.Run(func() {
				for item := range each.source {
					source <- item
				}
			})
		}

		group.Wait()
		close(source)
	}()

	return Range(source)
}

// Fork maybe very slow, must wait receive all to copy stream
func (s Stream) Fork() (Stream, Stream) {
	list := s.ForkN(2)
	return list[0], list[1]
}

func (s Stream) ForkTri() (Stream, Stream, Stream) {
	list := s.ForkN(3)
	return list[0], list[1], list[2]
}

func (s Stream) ForkN(n int) []Stream {
	if n <= 1 {
		return []Stream{s}
	}
	var cs []chan any
	for i := 0; i < n; i++ {
		cs = append(cs, make(chan any))
	}

	go func() {
		// TODO memory leak?
		var cache []any
		for item := range s.source {
			cache = append(cache, item)
		}

		for _, c := range cs {
			c := c
			go func() {
				for _, i := range cache {
					c <- i
				}
				close(c)
			}()
		}
	}()

	var result []Stream
	for _, c := range cs {
		result = append(result, Range(c))
	}
	return result
}

func (s Stream) ForkAn(consumers ...func(stream Stream)) {
	if len(consumers) <= 1 {
		return
	}

	var cs []chan any
	for i := 0; i < len(consumers); i++ {
		cs = append(cs, make(chan any))
	}

	go func() {
		// TODO memory leak?
		for item := range s.source {
			for _, c := range cs {
				//fmt.Println("fork to ", i, item, " size:", len(c))
				c <- item
			}
		}

		for _, c := range cs {
			close(c)
		}
	}()

	var wait sync.WaitGroup
	for i, c := range cs {
		wait.Add(1)

		iff := i
		cff := c
		go func() {
			consumers[iff](Range(cff))
			defer wait.Done()
		}()
	}

	wait.Wait()
}

// Distinct removes the duplicated items base on the given KeyFunc.
func (s Stream) Distinct(fn KeyFunc) Stream {
	source := make(chan any)

	GoSafe(func() {
		defer close(source)

		keys := make(map[any]PlaceholderType)
		for item := range s.source {
			key := fn(item)
			if _, ok := keys[key]; !ok {
				source <- item
				keys[key] = Placeholder
			}
		}
	})

	return Range(source)
}

// Done waits all upstreaming operations to be done.
func (s Stream) Done() {
	drain(s.source)
}

// Filter filters the items by the given FilterFunc.
func (s Stream) Filter(fn FilterFunc, opts ...Option) Stream {
	return s.Walk(func(item any, pipe chan<- any) {
		if fn(item) {
			pipe <- item
		}
	}, opts...)
}

// ForAll handles the streaming elements from the source and no later streams.
func (s Stream) ForAll(fn ForAllFunc) {
	fn(s.source)
	// avoid goroutine leak on fn not consuming all items.
	go drain(s.source)
}

// ForEach seals the Stream with the ForEachFunc on each item, no successive operations.
func (s Stream) ForEach(fn ForEachFunc) {
	for item := range s.source {
		fn(item)
	}
}

// Group groups the elements into different groups based on their keys.
func (s Stream) Group(fn KeyFunc, opts ...Option) Stream {
	source := make(chan any)
	go func() {
		option := buildOptions(opts...)
		groups := make(map[any][]any)

		if option.UnlimitedWorkers || option.Workers > 1 {
			s.groupParallel(option, fn, groups)
		} else {
			for item := range s.source {
				key := fn(item)
				groups[key] = append(groups[key], item)
			}
		}

		for k, group := range groups {
			source <- GroupItem{Key: k, Val: group}
		}
		close(source)
	}()

	return Range(source)
}

func (s Stream) groupParallel(option *RxOptions, fn KeyFunc, groups map[any][]any) {
	if option.UnlimitedWorkers {
		lock := sync.Mutex{}
		var wg sync.WaitGroup
		for item := range s.source {
			wg.Add(1)
			val := item
			GoSafe(func() {
				key := fn(val)
				lock.Lock()
				groups[key] = append(groups[key], val)
				lock.Unlock()
				wg.Done()
			})
		}
		wg.Wait()
	} else if option.Workers > 1 {
		pool := make(chan PlaceholderType, option.Workers)
		lock := sync.Mutex{}
		var wg sync.WaitGroup
		for item := range s.source {
			pool <- Placeholder
			wg.Add(1)

			val := item
			GoSafe(func() {
				defer func() {
					wg.Done()
					<-pool
				}()

				key := fn(val)
				lock.Lock()
				groups[key] = append(groups[key], val)
				lock.Unlock()
			})
		}
		wg.Wait()
	}
}

// Head returns the first n elements in p.
func (s Stream) Head(n int64) Stream {
	if n < 1 {
		panic("n must be greater than 0")
	}

	source := make(chan any)

	go func() {
		for item := range s.source {
			n--
			if n >= 0 {
				source <- item
			}
			if n == 0 {
				// let successive method go ASAP even we have more items to skip
				close(source)
				// why we don't just break the loop, and drain to consume all items.
				// because if breaks, this former goroutine will block forever,
				// which will cause goroutine leak.
				drain(s.source)
			}
		}
		// not enough items in s.source, but we need to let successive method to go ASAP.
		if n > 0 {
			close(source)
		}
	}()

	return Range(source)
}

// Map converts each item to another corresponding item, which means it's a 1:1 model.
func (s Stream) Map(fn MapFunc, opts ...Option) Stream {
	return s.Walk(func(item any, pipe chan<- any) {
		pipe <- fn(item)
	}, opts...)
}
func (s Stream) MapStr(opts ...Option) Stream {
	return s.Map(ToString, opts...)
}

// Flat make item flat to items
func (s Stream) Flat(flat func(any) Stream) Stream {
	// create current action channel
	source := make(chan any)
	// put data to channel by async
	go func() {
		for item := range s.source {
			for innerItem := range flat(item).source {
				source <- innerItem
			}
		}
		// close current channel
		close(source)
	}()

	// generate stream for next action
	return Range(source)
}

// Merge merges all the items into a slice and generates a new stream.
func (s Stream) Merge() Stream {
	var items []any
	for item := range s.source {
		items = append(items, item)
	}

	source := make(chan any, 1)
	source <- items
	close(source)

	return Range(source)
}

// Parallel applies the given ParallelFunc to each item concurrently with given number of Workers.
func (s Stream) Parallel(fn ParallelFunc, opts ...Option) {
	s.Walk(func(item any, pipe chan<- any) {
		fn(item)
	}, opts...).Done()
}

// Reduce is a utility method to let the caller deal with the underlying channel.
func (s Stream) Reduce(fn ReduceFunc) (any, error) {
	return fn(s.source)
}

// Reverse reverses the elements in the stream.
func (s Stream) Reverse() Stream {
	var items []any
	for item := range s.source {
		items = append(items, item)
	}
	// reverse, official method
	for i := len(items)/2 - 1; i >= 0; i-- {
		opp := len(items) - 1 - i
		items[i], items[opp] = items[opp], items[i]
	}

	return Just(items...)
}

// Skip returns a Stream that skips size elements.
func (s Stream) Skip(n int64) Stream {
	if n < 0 {
		panic("n must not be negative")
	}
	if n == 0 {
		return s
	}

	source := make(chan any)

	go func() {
		for item := range s.source {
			n--
			if n >= 0 {
				continue
			} else {
				source <- item
			}
		}
		close(source)
	}()

	return Range(source)
}

// Sort sorts the items from the underlying source.
func (s Stream) Sort(less LessFunc) Stream {
	var items []any
	for item := range s.source {
		items = append(items, item)
	}
	sort.Slice(items, func(i, j int) bool {
		return less(items[i], items[j])
	})

	return Just(items...)
}

// Split splits the elements into chunk with size up to n,
// might be less than n on tailing elements.
func (s Stream) Split(n int) Stream {
	if n < 1 {
		panic("n should be greater than 0")
	}

	source := make(chan any)
	go func() {
		var chunk []any
		for item := range s.source {
			chunk = append(chunk, item)
			if len(chunk) == n {
				source <- chunk
				chunk = nil
			}
		}
		if chunk != nil {
			source <- chunk
		}
		close(source)
	}()

	return Range(source)
}

// Tail returns the last n elements in p.
func (s Stream) Tail(n int64) Stream {
	if n < 1 {
		panic("n should be greater than 0")
	}

	source := make(chan any)

	go func() {
		ring := NewRing(int(n))
		for item := range s.source {
			ring.Add(item)
		}
		for _, item := range ring.Take() {
			source <- item
		}
		close(source)
	}()

	return Range(source)
}

// Walk lets the callers handle each item, the caller may write zero, one or more items base on the given item.
func (s Stream) Walk(fn WalkFunc, opts ...Option) Stream {
	option := buildOptions(opts...)
	if option.UnlimitedWorkers {
		return s.walkUnlimited(fn, option)
	}

	return s.walkLimited(fn, option)
}

func (s Stream) walkLimited(fn WalkFunc, option *RxOptions) Stream {
	pipe := make(chan any, option.Workers)

	go func() {
		var wg sync.WaitGroup
		pool := make(chan PlaceholderType, option.Workers)

		for item := range s.source {
			// important, used in another goroutine
			val := item
			pool <- Placeholder
			wg.Add(1)

			// better to safely run caller defined method
			GoSafe(func() {
				defer func() {
					wg.Done()
					<-pool
				}()

				fn(val, pipe)
			})
		}

		wg.Wait()
		close(pipe)
	}()

	return Range(pipe)
}

func (s Stream) walkUnlimited(fn WalkFunc, option *RxOptions) Stream {
	pipe := make(chan any, option.Workers)

	go func() {
		var wg sync.WaitGroup

		for item := range s.source {
			// important, used in another goroutine
			val := item
			wg.Add(1)
			// better to safely run caller defined method
			GoSafe(func() {
				defer wg.Done()
				fn(val, pipe)
			})
		}

		wg.Wait()
		close(pipe)
	}()

	return Range(pipe)
}

//////////
//  The following functions are all final operations
//////////

// First returns the first item, nil if no items.
func (s Stream) First() any {
	for item := range s.source {
		// make sure the former goroutine not block, and current func returns fast.
		go drain(s.source)
		return item
	}

	return nil
}

// Last returns the last item, or nil if no items.
func (s Stream) Last() (item any) {
	for item = range s.source {
	}
	return
}

// Min returns the minimum item from the underlying source.
func (s Stream) Min(less LessFunc) any {
	var min any
	for item := range s.source {
		if min == nil || less(item, min) {
			min = item
		}
	}

	return min
}

// Max returns the maximum item from the underlying source.
func (s Stream) Max(less LessFunc) any {
	var max any
	for item := range s.source {
		if max == nil || less(max, item) {
			max = item
		}
	}

	return max
}

// Count counts the number of elements in the result.
func (s Stream) Count() (count int) {
	for range s.source {
		count++
	}
	return
}

// AllMatch returns whether all elements of this stream match the provided predicate.
// May not evaluate the predicate on all elements if not necessary for determining the result.
// If the stream is empty then true is returned and the predicate is not evaluated.
func (s Stream) AllMatch(predicate func(item any) bool) bool {
	for item := range s.source {
		if !predicate(item) {
			// make sure the former goroutine not block, and current func returns fast.
			go drain(s.source)
			return false
		}
	}

	return true
}

// AnyMatch returns whether any elements of this stream match the provided predicate.
// May not evaluate the predicate on all elements if not necessary for determining the result.
// If the stream is empty then false is returned and the predicate is not evaluated.
func (s Stream) AnyMatch(predicate func(item any) bool) bool {
	for item := range s.source {
		if predicate(item) {
			// make sure the former goroutine not block, and current func returns fast.
			go drain(s.source)
			return true
		}
	}

	return false
}

// NoneMatch returns whether all elements of this stream don't match the provided predicate.
// May not evaluate the predicate on all elements if not necessary for determining the result.
// If the stream is empty then true is returned and the predicate is not evaluated.
func (s Stream) NoneMatch(predicate func(item any) bool) bool {
	for item := range s.source {
		if predicate(item) {
			// make sure the former goroutine not block, and current func returns fast.
			go drain(s.source)
			return false
		}
	}

	return true
}
