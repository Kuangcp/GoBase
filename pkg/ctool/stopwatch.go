package ctool

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
	// StopWatch inspire by spring stopwatch
	StopWatch struct {
		name      string
		first     bool
		hasStart  bool
		firstTime time.Time
		startTime *time.Time
		stopTime  *time.Time
		tasks     []TaskInfo
	}
)

func NewStopWatch() *StopWatch {
	return &StopWatch{first: true}
}

func NewStopWatchWithName(name string) *StopWatch {
	return &StopWatch{name: name, first: true}
}

func fmtDuration(d time.Duration) string {
	ms := d.Milliseconds()
	d = d.Round(time.Millisecond)
	if ms < 10_000 {
		return fmt.Sprintf("%vms", ms)
	}
	return d.String()
}

func (s StopWatch) PrettyPrint() string {
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
	totalElapsed := s.stopTime.Sub(s.firstTime)
	totalNano := totalElapsed.Nanoseconds()
	for _, task := range s.tasks {

		var rate int64
		if totalNano == 0 {
			rate = 0
		} else {
			rate = task.elapsedTime.Nanoseconds() * 100 / totalNano
		}

		taskStr += fmt.Sprintf("%9v%3v%% %v\n", fmtDuration(task.elapsedTime), rate, task.name)
	}

	return fmt.Sprintf("\nStopWatch %s : %v | %v => %v\n%v", s.name, fmtDuration(totalElapsed),
		s.firstTime.Format("15:04:05.000"),
		s.stopTime.Format("15:04:05.000"), taskStr)
}

func (s *StopWatch) StartAnon() {
	s.Start("")
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
	s.startTime = &now
	s.hasStart = true
	s.tasks = append(s.tasks, TaskInfo{name: name})
}

func (s *StopWatch) Stop() {
	if !s.hasStart {
		return
	}
	now := time.Now()
	s.stopTime = &now
	s.hasStart = false
	s.tasks[len(s.tasks)-1].elapsedTime = now.Sub(*s.startTime)
}

func (s StopWatch) GetTotalDuration() time.Duration {
	if s.stopTime == nil {
		return time.Now().Sub(s.firstTime)
	}
	return s.stopTime.Sub(s.firstTime)
}

func (s StopWatch) GetTaskCount() int {
	return len(s.tasks)
}
