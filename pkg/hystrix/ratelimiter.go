package hystrix

import "time"

type RateLimiter interface {
	SetRate(rate int)
	GetRate() int

	Acquire() int64
	AcquireN(n int) int64

	TryAcquire() bool
	TryAcquireN(n int) bool
}

func CreateLeakyLimiter(rate int) RateLimiter {
	l := &LeakyBucket{}
	l.SetRate(rate)
	l.buffer = make(chan int, rate)
	go l.producer()
	return l
}

// LeakyBucket 漏桶实现
type LeakyBucket struct {
	rate   int
	buffer chan int
}

func (l *LeakyBucket) GetRate() int {
	return l.rate
}

func (l *LeakyBucket) producer() {
	for {
		rate := l.GetRate()
		time.Sleep(time.Microsecond * time.Duration(1000_000/rate))
		l.buffer <- 0
	}
}
func (l *LeakyBucket) TryAcquire() bool {
	return len(l.buffer) > 0
}

func (l *LeakyBucket) TryAcquireN(n int) bool {
	if n < 0 {
		return false
	}
	return len(l.buffer) > n
}

func (l *LeakyBucket) SetRate(rate int) {
	if rate < 0 {
		rate = 1
	}
	l.rate = rate
}

func (l *LeakyBucket) Acquire() int64 {
	start := time.Now().UnixMicro()
	<-l.buffer
	return time.Now().UnixMicro() - start
}

func (l LeakyBucket) AcquireN(n int) int64 {
	start := time.Now().UnixMicro()
	for i := 0; i < n; i++ {
		<-l.buffer
	}
	return time.Now().UnixMicro() - start
}
