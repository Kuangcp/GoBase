package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/kuangcp/gobase/pkg/cuibase"
	"github.com/kuangcp/logger"
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
	pidFile       string
	retentionTime time.Duration
	checkPeriod   time.Duration
)

var (
	help         bool
	suffix       string
	check        bool
	daemon       bool
	debug        bool
	exit         bool
	illegalQuit  bool
	listTrash    bool
	retentionStr string // time.ParseDuration()
	checkStr     string
)

func init() {
	logger.SetLogPathTrim("recycle-bin")

	home, err := cuibase.Home()
	cuibase.CheckIfError(err)

	mainDir = home + mainDir
	configDir = mainDir + "/conf"
	pidFile = configDir + "/pid"

	logDir = mainDir + "/log"
	logFile = logDir + "/main.log"

	trashDir = mainDir + "/trash"

	_ = logger.SetLoggerConfig(&logger.LogConfig{
		Console: &logger.ConsoleLogger{
			Level:    logger.DebugDesc,
			Colorful: true,
		},
		File: &logger.FileLogger{
			Filename:   logFile,
			Level:      logger.DebugDesc,
			Colorful:   true,
			Append:     true,
			PermitMask: "0660",
		},
	})

	flag.BoolVar(&help, "h", false, "")
	flag.BoolVar(&help, "H", false, "")
	flag.BoolVar(&debug, "D", false, "")
	flag.BoolVar(&check, "C", false, "")
	flag.BoolVar(&daemon, "d", false, "")
	flag.BoolVar(&exit, "X", false, "")
	flag.BoolVar(&illegalQuit, "q", false, "")
	flag.BoolVar(&listTrash, "l", false, "")

	flag.StringVar(&retentionStr, "r", "168h", "")
	flag.StringVar(&checkStr, "c", "1h", "")
	flag.StringVar(&suffix, "s", "", "")

	flag.Usage = info.PrintHelp
	flag.Parse()
}

func main() {
	if help {
		InitDir()
		info.PrintHelp()
		return
	}

	if suffix != "" {
		DeleteFileBySuffix(strings.Split(suffix, ","))
		return
	}

	if listTrash {
		ListTrashFiles()
		return
	}

	if check {
		if daemon {
			CheckWithDaemon()
		} else {
			CheckTrashDir()
		}
		return
	}

	if exit {
		ExitCheckFileDaemon()
		return
	}

	if illegalQuit {
		ActualDeleteFile(pidFile)
		return
	}

	args := os.Args
	if len(args) == 1 {
		info.PrintHelp()
	} else {
		DeleteFiles(args[1:])
		CheckWithDaemon()
	}
}

func ListTrashFiles() {
	err := parseTime()
	if err != nil {
		return
	}
	dir, err := ioutil.ReadDir(trashDir)
	if err != nil {
		logger.Error(err)
		return
	}

	current := time.Now().UnixNano()
	if len(dir)!= 0 {
		fmt.Printf("%-23s %-10s %s\n", "DeleteTime", "Remaining", "File")
	}
	for _, fileInfo := range dir {
		name := fileInfo.Name()
		index := strings.Index(name, ".T.")
		if index == -1 {
			continue
		}

		file := name[:index]

		value := name[index+3:]
		//fmt.Println(value)
		parseInt, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			logger.Error(err)
			return
		}
		duration, err := time.ParseDuration(strconv.FormatInt((retentionTime.Nanoseconds()-current+parseInt)/1000000000, 10) + "s")
		if err != nil {
			duration = 0

		}
		fmt.Println(time.Unix(parseInt/1000000000, 0).Format("2006-01-02 15:04:05.000"),
			cuibase.Yellow.Printf("%10s", duration.String()), cuibase.Green.Print(file))
	}
}

func CheckWithDaemon() {
	params := fmt.Sprintf(" -c %s -r %s", checkStr, retentionStr)
	proc, err := startProc([]string{"/usr/bin/bash", "-c", "recycle-bin -C" + params}, logFile)
	if err != nil {
		logger.Error(proc, err)
	}
}

func CheckTrashDir() {
	err := parseTime()
	if err != nil {
		return
	}

	exists, err := isPathExists(pidFile)
	if exists {
		return
	}
	logger.Info("Start check trash, period:", checkPeriod, "pid:", os.Getpid())
	if err != nil {
		logger.Error(err)
		return
	}

	// avoid repeat delete
	var deleteFlag int32 = 0

	go func() {
		// Wait for interrupt signal to gracefully shutdown the server with
		// a timeout of 5 seconds.
		quit := make(chan os.Signal)
		// kill (no param) default send syscall.SIGTERM
		// kill -2 is syscall.SIGINT
		// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit

		logger.Warn("Killed")
		deletePidFile(&deleteFlag)
		os.Exit(1)
	}()

	// create pid
	err = ioutil.WriteFile(pidFile, []byte(strconv.Itoa(os.Getpid())), 0644)
	if err != nil {
		logger.Error(err)
		return
	}

	defer deletePidFile(&deleteFlag)

	emptyCount := 0
	for true {
		current := time.Now().UnixNano()
		time.Sleep(checkPeriod)
		logger.Debug("Check")
		dir, err := ioutil.ReadDir(trashDir)
		if err != nil {
			logger.Error(err)
			return
		}

		if len(dir) == 0 {
			emptyCount++
		}
		if emptyCount >= maxEmptyTrashCheck {
			return
		}

		for _, fileInfo := range dir {
			name := fileInfo.Name()
			index := strings.Index(name, ".T.")
			if index == -1 {
				continue
			}

			value := name[index+3:]
			//fmt.Println(value)
			parseInt, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				logger.Error(err)
				return
			}

			//logger.Debug(current, parseInt, current-parseInt)
			if current-parseInt > retentionTime.Nanoseconds() {
				logger.Warn("Delete: ", name[:index])
				actualPath := trashDir + "/" + name
				if actualPath == "/" {
					logger.Error("danger error")
					continue
				}
				cmd := exec.Command("rm", "-rf", actualPath)
				execCmdWithQuite(cmd)
			}
		}
	}
}

func parseTime() error {
	duration, err := time.ParseDuration(retentionStr)
	if err != nil {
		logger.Error(err)
		return nil
	}

	retentionTime = duration
	checkPeriod, err = time.ParseDuration(checkStr)
	if err != nil {
		logger.Error(err)
		return nil
	}
	return err
}

func deletePidFile(deleteFlag *int32) {
	logger.Warn("Exit")
	curDelete := atomic.AddInt32(deleteFlag, 1)
	if curDelete == 1 {
		ActualDeleteFile(pidFile)
	}
}

// deleteFies 移动文件到回收站
func DeleteFiles(files []string) {
	for _, filepath := range files {
		exists, err := isPathExists(filepath)
		cuibase.CheckIfError(err)
		if !exists {
			logger.Error(filepath, "not found")
			return
		}

		logger.Info("Prepare delete:", filepath)

		timestamp := strconv.FormatInt(time.Now().UnixNano(), 10)
		//logger.Debug(filepath, trashDir+"/"+filepath)
		cmd := exec.Command("mv", filepath, trashDir+"/"+filepath+".T."+timestamp)
		execCmdWithQuite(cmd)
	}
}

func DeleteFileBySuffix(params []string) {
	if len(params) == 0 {
		return
	}
	fmt.Println(params)
}

func ExitCheckFileDaemon() {
	exists, err := isPathExists(pidFile)
	if !exists {
		logger.Error("no pid file")
		return
	}
	file, err := ioutil.ReadFile(pidFile)
	if err != nil {
		logger.Error(err)
		return
	}

	pid := string(file)
	logger.Info("kill ", pid)
	cmd := exec.Command("kill", pid)
	execCmdWithQuite(cmd)
}

func ActualDeleteFile(path string) {
	err := os.Remove(path)
	if err != nil {
		logger.Error(err)
	}
}

// 静默执行 不关心返回值
func execCmdWithQuite(cmd *exec.Cmd) {
	var out bytes.Buffer

	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		logger.Error(err)
		return
	}
}
