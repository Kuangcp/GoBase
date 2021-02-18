package queue

type (
	Queue struct {
		queue []interface{}
	}
)

// https://blog.wolfogre.com/posts/slice-queue-vs-list-queue/
func New() *Queue {
	return &Queue{}
}

func (e *Queue) Push(ele interface{}) {
	e.queue = append(e.queue, ele)
}

func (e *Queue) Peek() *interface{} {
	if e.IsEmpty() {
		return nil
	}
	return &e.queue[0]
}
func (e *Queue) Len() int {
	return len(e.queue)
}

func (e *Queue) Pop() *interface{} {
	if e.IsEmpty() {
		return nil
	}
	result := e.queue[0]
	e.queue = e.queue[1:]
	return &result
}

func (e *Queue) IsEmpty() bool {
	return e.Len() == 0
}
