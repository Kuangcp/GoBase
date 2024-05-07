package core

import (
	"github.com/kuangcp/logger"
	"testing"
	"time"
)

func TestRegCache(t *testing.T) {
	cache := make(map[string]string)
	go func() {
		for t2 := range time.NewTicker(time.Millisecond * 100).C {
			cache[t2.String()] = ""
		}
	}()

	go func() {
		for range time.NewTicker(time.Millisecond * 100).C {
			for k, v := range cache {
				logger.Info(k, v)
			}
		}
	}()

	time.Sleep(time.Second * 10)
}
