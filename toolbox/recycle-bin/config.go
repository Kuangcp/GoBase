package main

import (
	"fmt"
	"os"
	"strings"
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
	sysMap        = make(map[string]int8)
	sysDir        = [...]string{"/bin/", "/boot/", "/data/", "/dev/", "/etc/", "/home/", "/lib/", "/lib64/",
		"/lost+found/", "/mnt/", "/opt/", "/proc/", "/root/", "/run/", "/sbin/", "/srv/", "/sys/", "/tmp/",
		"/usr/", "/var/"}
)

type (
	fileItem struct {
		name      string
		timestamp int64
		file      os.FileInfo
	}
)

func init() {
	for _, s := range sysDir {
		sysMap[s] = 0
	}
}

// 高危动作目录
func isDangerDir(dir string) bool {
	count := strings.Count(dir, "/")
	if count == 1 {
		return true
	}

	_, ok := sysMap[dir]
	return ok
}

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
