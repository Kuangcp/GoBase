package wxrobot

import (
	"sync"
	"time"
)

type (
	PeriodRateLimiter struct {
		lock        *sync.RWMutex
		maxCount    int
		slideWindow time.Duration
		curCount    int
		slideQueue  *Queue // queue task entry nano seconds
	}
)

// NewMinuteLimiter 限流: maxCount/min
func NewMinuteLimiter(maxCount int) *PeriodRateLimiter {
	return NewLimiter(time.Minute, maxCount)
}

// NewSecondLimiter 限流：maxCount/s
func NewSecondLimiter(maxCount int) *PeriodRateLimiter {
	return NewLimiter(time.Second, maxCount)
}

func NewLimiter(slideWindow time.Duration, maxCount int) *PeriodRateLimiter {
	return &PeriodRateLimiter{
		maxCount:    maxCount,
		slideWindow: slideWindow,
		curCount:    0,
		slideQueue:  NewQueue(),
		lock:        &sync.RWMutex{},
	}
}

// calculateCount 移除 滑动窗口外的元素
//  简单压测可发现 队列重整理耗时很小
func (l *PeriodRateLimiter) calculateCount() {
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

	//start := time.Now().UnixNano()
	l.calculateCount()
	//end := time.Now().UnixNano()
	//fmt.Println("queue waste: ", end-start)

	count := l.curCount
	maxCount := l.maxCount
	if count >= maxCount {
		return false
	}

	l.slideQueue.Push(time.Now().UnixNano())
	return true
}

func (l *PeriodRateLimiter) queueState() int {
	l.lock.RLock()
	defer l.lock.RUnlock()
	return l.curCount
}
