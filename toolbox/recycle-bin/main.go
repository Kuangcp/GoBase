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
	maxEmptyTrashCheck = 10
)

var (
	mainDir     = "/.config/app-conf/recycle-bin"
	configDir   string
	logDir      string
	trashDir    string
	logFile     string
	pidFile     string
	liveTime    time.Duration
	checkPeriod time.Duration
)

var (
	help        bool
	suffix      string
	check       bool
	daemon      bool
	debug       bool
	exit        bool
	illegalQuit bool
	liveStr     string // time.ParseDuration()
	checkStr    string
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

	logger.SetLogger("{\"Console\": {\"level\": \"DEBG\",\"color\": true},\"File\":{\"filename\": \"" + logFile + "\",\"level\": \"DEBG\",\"color\": true,\"append\": true,\"permit\": \"0660\"}}")

	flag.BoolVar(&help, "h", false, "")
	flag.BoolVar(&help, "H", false, "")
	flag.BoolVar(&debug, "D", false, "")
	flag.BoolVar(&check, "C", false, "")
	flag.BoolVar(&daemon, "d", false, "")
	flag.BoolVar(&exit, "X", false, "")
	flag.BoolVar(&illegalQuit, "q", false, "")

	flag.StringVar(&liveStr, "l", "30s", "")
	flag.StringVar(&checkStr, "c", "5s", "")
	flag.StringVar(&suffix, "s", "", "")

	flag.Usage = info.PrintHelp
	flag.Parse()
}

func main() {
	if help {
		initDir()
		info.PrintHelp()
		return
	}

	if suffix != "" {
		deleteFileBySuffix(strings.Split(suffix, ","))
		return
	}

	if check {
		if daemon {
			checkWithDaemon()
		} else {
			checkTrashDir()
		}
		return
	}
	if exit {
		exitCheckFileDaemon()
		return
	}
	if illegalQuit {
		actualDeleteFile(pidFile)
		return
	}

	args := os.Args
	if len(args) == 1 {
		info.PrintHelp()
	} else {
		deleteFiles(args[1:])
		checkWithDaemon()
	}
}

func checkWithDaemon() {
	params := fmt.Sprintf(" -c %s -l %s", checkStr, liveStr)
	proc, err := startProc([]string{"/usr/bin/bash", "-c", "recycle-bin -C" + params}, logFile)
	if err != nil {
		logger.Error(proc, err)
	}
}

func checkTrashDir() {
	duration, err := time.ParseDuration(liveStr)
	if err != nil {
		logger.Error(err)
		return
	}

	liveTime = duration
	checkPeriod, err = time.ParseDuration(checkStr)
	if err != nil {
		logger.Error(err)
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
			return
		}

		if len(dir) == 0 {
			emptyCount++
		}
		if emptyCount > maxEmptyTrashCheck {
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
			if current-parseInt > liveTime.Nanoseconds() {
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

func deletePidFile(deleteFlag *int32) {
	logger.Warn("Exit")
	curDelete := atomic.AddInt32(deleteFlag, 1)
	if curDelete == 1 {
		actualDeleteFile(pidFile)
	}
}

// deleteFies 移动文件到回收站
func deleteFiles(files []string) {
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

func deleteFileBySuffix(params []string) {
	if len(params) == 0 {
		return
	}
	fmt.Println(params)
}

func exitCheckFileDaemon() {
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

func actualDeleteFile(path string) {
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
