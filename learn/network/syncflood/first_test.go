package main

import (
	"github.com/kuangcp/logger"
	"testing"
	"time"
)

func TestTcp(t *testing.T) {
	max := 20000

	addr = "192.168.130.209:3306"
	con = 2000
	for i := 0; i < max/2; i++ {
		go createTcp()
		go createTcp()
	}

	for range time.NewTicker(time.Second).C {
		cur := set.Len()
		logger.Info("size: ", cur)
		if cur < max {
			for i := cur; i < max/2; i++ {
				go createTcp()
				go createTcp()
			}
		}
	}
	time.Sleep(time.Second * 120)
}
