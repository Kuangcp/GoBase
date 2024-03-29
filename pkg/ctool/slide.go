package ctool

import (
	"sync"
	"time"
)

type (
	PeriodRateLimiter struct {
		lock        *sync.RWMutex
		maxCount    int
		slideWindow time.Duration
		curCount    int           // already run task in current windows
		slideQueue  *Queue[int64] // queue task entry nano seconds
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
		slideQueue:  NewQueue[int64](),
		lock:        &sync.RWMutex{},
	}
}

// calculateCount 移除 滑动窗口外的元素
//  简单压测可发现 队列重整理耗时很小
func (l *PeriodRateLimiter) CalculateCount() {
	nowNs := time.Now().UnixNano()
	for {
		peek := l.slideQueue.Peek()
		if peek == 0 {
			break
		}
		peekVal := peek
		if nowNs-peekVal < l.slideWindow.Nanoseconds() {
			break
		}
		l.slideQueue.Pop()
	}
	l.curCount = l.slideQueue.Len()
}

func (l *PeriodRateLimiter) Acquire() bool {
	l.lock.Lock()
	defer l.lock.Unlock()

	//start := time.Now().UnixNano()
	l.CalculateCount()
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

func (l *PeriodRateLimiter) QueueState() int {
	l.lock.RLock()
	defer l.lock.RUnlock()
	return l.curCount
}
