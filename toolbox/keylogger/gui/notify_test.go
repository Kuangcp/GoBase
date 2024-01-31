package main

import (
	"fmt"
	"github.com/kuangcp/logger"
	"github.com/shirou/gopsutil/cpu"
	"math"
	"testing"
	"time"
)

func TestNotify(t *testing.T) {
	//notifyAny()

	duration, err := time.ParseDuration("2m")
	if err != nil {
		logger.Error(err)
	}
	logger.Info(duration)
}

func TestCpu(t *testing.T) {
	for i := 0; i < 100; i++ {
		go func() {
			memInfo, _ := cpu.Percent(time.Second, false)
			y := math.Round((100 - memInfo[0]) * height / 100)
			fmt.Println(y)
		}()
	}
	time.Sleep(time.Minute)
}
