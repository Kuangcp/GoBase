package ctool

type (
	Ring[T any] interface {
		PeekN(int) []T
		Peek() *T
		GetN(int) []T
		Get() *T
		Add(...T)
		SetCapacity(int)
	}

	ARing[T any] struct {
		data     []*T
		capacity int
		start    int
		end      int
	}
)

func NewARing[T any](capacity int) *ARing[T] {
	return &ARing[T]{capacity: capacity, data: make([]*T, capacity), start: 0, end: 0}
}

func (A *ARing[T]) PeekN(length int) []*T {
	if A.end == A.start {
		return nil
	}
	if A.start < A.end {
		delta := min(length, A.end-A.start)
		return A.data[A.start : A.start+delta]
	} else {
		return nil
	}
}

func (A *ARing[T]) Peek() *T {
	//TODO implement me
	panic("implement me")
}

func (A *ARing[T]) GetN(i int) []T {
	//TODO implement me
	panic("implement me")
}

func (A *ARing[T]) Get() *T {
	//TODO implement me
	panic("implement me")
}

func (A *ARing[T]) Add(t ...T) {
	//TODO implement me
	panic("implement me")
}

func (A *ARing[T]) SetCapacity(i int) {
	//TODO implement me
	panic("implement me")
}
