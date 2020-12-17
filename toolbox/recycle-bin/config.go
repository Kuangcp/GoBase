package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/kuangcp/gobase/pkg/cuibase"
)

const (
	maxEmptyTrashCheck = 3
)

var (
	mainDir       = "/.config/app-conf/recycle-bin"
	configDir     string
	logDir        string
	trashDir      string
	logFile       string
	configFile    string
	pidFile       string
	retentionTime time.Duration
	checkPeriod   time.Duration
)

type (
	fileItem struct {
		name      string
		timestamp int64
		file      os.FileInfo
	}
)

func (t *fileItem) formatTime() string {
	return time.Unix(t.timestamp/1000000000, 0).Format("2006-01-02 15:04:05.000")
}

func (t *fileItem) formatForList(current int64) string {
	second := strconv.FormatInt((retentionTime.Nanoseconds()-current+t.timestamp)/1000000000, 10)

	duration, err := time.ParseDuration(second + "s")
	if err != nil {
		duration = 0
	}
	if duration.Seconds() < 0 {
		duration = 0
	}
	return fmt.Sprintln(t.formatTime(), cuibase.Yellow.Print(fmtDuration(duration)), cuibase.Green.Print(t.name))
}
