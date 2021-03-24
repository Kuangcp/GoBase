package stopwatch

import (
	"fmt"
	"time"
)

type (
	TaskInfo struct {
		name        string
		elapsedTime time.Duration
	}
	StopWatch struct {
		name      string
		first     bool
		hasStart  bool
		firstTime time.Time
		startTime time.Time
		stopTime  time.Time
		tasks     []TaskInfo
	}
)

func New() *StopWatch {
	return &StopWatch{first: true}
}

func NewWithName(name string) *StopWatch {
	return &StopWatch{name: name, first: true}
}

func (s *StopWatch) PrettyPrint() string {
	if s.first {
		return ""
	}
	taskStr := ""
	for i, task := range s.tasks {
		taskStr += fmt.Sprintf("%3d %13s : %v\n", i+1, task.name, task.elapsedTime.String())
	}
	return fmt.Sprintf("%s : %v   %v %v\n%v", s.name, s.stopTime.Sub(s.firstTime),
		s.firstTime.Format("15:04:05.000"),
		s.stopTime.Format("15:04:05.000"), taskStr)
}

func (s *StopWatch) Start(name string) {
	now := time.Now()
	if s.first {
		s.firstTime = now
		s.first = false
	}
	s.startTime = now
	s.hasStart = true
	s.tasks = append(s.tasks, TaskInfo{name: name})
}

func (s *StopWatch) Stop() {
	if !s.hasStart {
		return
	}
	now := time.Now()
	s.stopTime = now
	s.hasStart = false
	s.tasks[len(s.tasks)-1].elapsedTime = now.Sub(s.startTime)
}
