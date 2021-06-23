package stopwatch

import (
	"fmt"
	"time"
	"unicode/utf8"
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
	if s.hasStart {
		s.Stop()
	}
	taskStr := ""
	maxNameLen := 0
	for _, task := range s.tasks {
		curLen := utf8.RuneCountInString(task.name)
		if maxNameLen < curLen {
			maxNameLen = curLen
		}
	}
	for i, task := range s.tasks {
		taskStr += fmt.Sprintf("%3d %-"+fmt.Sprint(maxNameLen)+"s : %v\n",
			i+1, task.name, task.elapsedTime.String())
	}
	return fmt.Sprintf("\n%s : %v   %v %v\n%v", s.name, s.stopTime.Sub(s.firstTime),
		s.firstTime.Format("15:04:05.000"),
		s.stopTime.Format("15:04:05.000"), taskStr)
}

func (s *StopWatch) Start(name string) {
	now := time.Now()
	if s.first {
		s.firstTime = now
		s.first = false
	}
	if s.hasStart {
		s.Stop()
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
