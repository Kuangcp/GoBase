package main

import (
	"fmt"
	"os"
	"time"

	"github.com/kuangcp/gobase/pkg/cuibase"
)

const (
	maxEmptyTrashCheck = 3
	appName            = "recycle-bin"
)

var (
	mainDir       = "/.config/app-conf/" + appName
	configDir     string
	logDir        string
	trashDir      string
	logFile       string
	configFile    string
	pidFile       string
	retentionTime time.Duration
	checkPeriod   time.Duration
	sysDir        = [...]string{"/", "/home", "/home/", "/ect", "/etc/", "/boot", "/boot/", "/sys", "/sys/",
		"/opt", "/opt/", "/bin", "/bin/"}
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
	return time.Unix(t.seconds(), 0).Format(cuibase.YYYY_MM_DD_HH_MM_SS_MS)
}

func (t *fileItem) formatForList(index int, currentNano int64) string {
	duration := time.Duration(t.timestamp - currentNano + retentionTime.Nanoseconds())
	if duration < 0 {
		duration = 0
	}
	return fmt.Sprintf("%-3v %v %v %v\n", index+1, t.formatTime(),
		cuibase.Yellow.Print(fmtDuration(duration)), cuibase.Green.Print(t.name))
}
