package stream

// A TaskRunner is used to control the concurrency of goroutines.
type TaskRunner struct {
	limitChan chan PlaceholderType
}

// NewTaskRunner returns a TaskRunner.
func NewTaskRunner(concurrency int) *TaskRunner {
	return &TaskRunner{
		limitChan: make(chan PlaceholderType, concurrency),
	}
}

// Schedule schedules a task to run under concurrency control.
func (rp *TaskRunner) Schedule(task func()) {
	rp.limitChan <- Placeholder

	go func() {
		defer Recover(func() {
			<-rp.limitChan
		})

		task()
	}()
}
