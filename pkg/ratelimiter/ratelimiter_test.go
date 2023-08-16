package ratelimiter

import (
	"log"
	"testing"
	"time"
)

func TestLeaky(t *testing.T) {
	limiter := CreateLeakyLimiter(2)
	ls := []int{1, 2, 7, 2, 4, 1, 1, 9, 1}
	for _, i := range ls {
		log.Println("start", i)
		limiter.AcquireN(i)
		log.Println("end", i)
	}
}

func TestConcur(t *testing.T) {
	limiter := CreateLeakyLimiter(7)
	ls := []int{1, 2, 7, 2, 4, 1, 1, 9, 1}
	for _, i := range ls {

		fi := i

		for k := 0; k < fi; k++ {
			go func() {
				log.Println("start", fi)
				limiter.Acquire()
				log.Println("end", fi)
			}()
		}
	}

	time.Sleep(time.Hour)
}

func TestTimeout(t *testing.T) {
	limiter := CreateLeakyLimiter(7)
	ls := []int{4, 6, 7, 2, 8, 5, 1, 9, 8}
	for _, i := range ls {

		fi := i
		go func() {
			for k := 0; k < fi; k++ {
				//log.Println("start", fi)
				wait := limiter.TryAcquireWait(time.Second * 1)
				if wait {
					limiter.Acquire()
					log.Println("end", fi)
				} else {
					log.Println("timeout", fi)
				}
			}
		}()
	}

	time.Sleep(time.Hour)
}
