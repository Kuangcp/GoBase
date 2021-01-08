package main

import (
	"fmt"
	"os"
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

func (t *fileItem) seconds() int64 {
	return t.timestamp / 1000_000_000
}

func (t *fileItem) formatTime() string {
	return time.Unix(t.seconds(), 0).Format("2006-01-02 15:04:05.000")
}

func (t *fileItem) formatForList(currentNano int64) string {
	duration := time.Duration(t.timestamp - currentNano + retentionTime.Nanoseconds())
	if duration < 0 {
		duration = 0
	}
	return fmt.Sprintln(t.formatTime(), cuibase.Yellow.Print(fmtDuration(duration)), cuibase.Green.Print(t.name))
}
