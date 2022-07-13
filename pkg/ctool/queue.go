package ctool

type (
	Queue[T any] struct {
		queue []T
	}
)

// https://blog.wolfogre.com/posts/slice-queue-vs-list-queue/
func NewQueue[T any]() *Queue[T] {
	return &Queue[T]{}
}

func (e *Queue[T]) Push(ele T) {
	e.queue = append(e.queue, ele)
}

func (e *Queue[T]) Peek() (t T) {
	if e.IsEmpty() {
		return
	}
	return e.queue[0]
}

func (e *Queue[T]) Len() int {
	return len(e.queue)
}

func (e *Queue[T]) Pop() (t T) {
	if e.IsEmpty() {
		return
	}
	result := e.queue[0]
	e.queue = e.queue[1:]
	return result
}

func (e *Queue[T]) IsEmpty() bool {
	return e.Len() == 0
}
