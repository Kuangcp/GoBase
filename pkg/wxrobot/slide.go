package wxrobot

import (
	"sync"
	"time"
)

const DefaultSlideWindow = time.Minute // 滑动窗口 默认值

type (
	// 大周期限速，例如 50次/min
	PeriodRateLimiter struct {
		lock        *sync.Mutex
		maxCount    int
		curCount    int
		slideQueue  *Queue
		calPeriod   time.Duration
		slideWindow time.Duration
	}
)

func NewLimiter(maxCount int) *PeriodRateLimiter {
	return NewCustomLimiter(DefaultSlideWindow, maxCount)
}

func NewCustomLimiter(slideWindow time.Duration, maxCount int) *PeriodRateLimiter {
	return &PeriodRateLimiter{
		maxCount:    maxCount,
		slideWindow: slideWindow,
		curCount:    0,
		slideQueue:  NewQueue(),
		lock:        &sync.Mutex{},
	}
}

func (l *PeriodRateLimiter) calculateCount() {
	// remove element that outer of time window
	nowNs := time.Now().UnixNano()
	for {
		peek := l.slideQueue.Peek()
		if peek == nil {
			break
		}
		peekVal := (*peek).(int64)
		if nowNs-peekVal < l.slideWindow.Nanoseconds() {
			break
		}
		l.slideQueue.Pop()
	}
	l.curCount = l.slideQueue.Len()
}

func (l *PeriodRateLimiter) acquire() bool {
	l.lock.Lock()
	defer l.lock.Unlock()

	l.calculateCount()
	count := l.curCount
	maxCount := l.maxCount
	acquire := count < maxCount
	if acquire {
		l.slideQueue.Push(time.Now().UnixNano())
	}

	return acquire
}
