package wxrobot

import (
	"fmt"
	"log"
	"math/rand"
	"testing"
	"time"
)

func TestRateLimiter(t *testing.T) {
	limiter := NewLimiter(time.Millisecond*1777, 35)

	//limiter.acquire()
	//time.Sleep(time.Millisecond * 50)
	//limiter.acquire()
	//time.Sleep(time.Millisecond * 70)

	go producer(limiter)
	go producer(limiter)
	go producer(limiter)
	go producer(limiter)
	go producer(limiter)
	ticker := time.NewTicker(time.Millisecond * 500)
	for range ticker.C {
		fmt.Println()
		//log.Println("[", limiter.slideQueue.Len(), "]")
		log.Println("[", limiter.queueState(), "]")
	}
}

func producer(limiter *PeriodRateLimiter) {
	ticker := time.NewTicker(time.Millisecond * 37)
	for range ticker.C {
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(70)+37))
		acquire := limiter.acquire()
		if acquire {
			fmt.Printf("%v▶%v", "\033[0;32m", "\033[0m")
		} else {
			fmt.Print(".")
		}
		//log.Println("[", limiter.slideQueue.Len(), "]")
	}
}
