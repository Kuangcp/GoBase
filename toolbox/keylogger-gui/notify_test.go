package main

import (
	"github.com/kuangcp/logger"
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
