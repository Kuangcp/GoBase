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
	"github.com/wonderivan/logger"
	"github.com/zh-five/xdaemon"
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
	help     bool
	suffix   string
	check    bool
	debug    bool
	exit     bool
	liveStr  string // time.ParseDuration()
	checkStr string
)
var info = cuibase.HelpInfo{
	Description:   "Recycle bin",
	Version:       "1.0.0",
	SingleFlagLen: -3,
	ValueLen:      -10,
	Flags: []cuibase.ParamVO{
		{
			Short:   "-h",
			Value:   "",
			Comment: "Help info and init",
		}, {
			Short:   "-D",
			Value:   "",
			Comment: "Debug mode",
		}, {
			Short:   "-C",
			Value:   "",
			Comment: "Start check",
		},
	},
	Options: []cuibase.ParamVO{
		{
			Short:   "",
			Value:   "file",
			Comment: "Delete file",
		}, {
			Short:   "-s",
			Value:   "suffix",
			Comment: "Delete *.suffix",
		}, {
			Short:   "-l",
			Value:   "duration",
			Comment: "File live time, default 1m. (unit: ms/s/m/h) ",
		}, {
			Short:   "-c",
			Value:   "duration",
			Comment: "Check trash period, default 1m. (unit: ms/s/m/h)",
		},
	}}

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

	flag.BoolVar(&help, "h", false, "")
	flag.BoolVar(&help, "H", false, "")
	flag.BoolVar(&debug, "D", false, "")
	flag.BoolVar(&check, "C", false, "")
	flag.BoolVar(&exit, "X", false, "")
	flag.StringVar(&liveStr, "l", "1m", "")
	flag.StringVar(&checkStr, "c", "1m", "")
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
		checkTrashDir()
		return
	}
	if exit {
		exitCheckFileDaemon()
		return
	}

	args := os.Args
	if len(args) == 1 {
		info.PrintHelp()
	} else {
		deleteFiles(args[1:])
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

	checkTrashDaemon()
}

func initDir() {
	err := os.MkdirAll(trashDir, 0755)
	cuibase.CheckIfError(err)
	err = os.MkdirAll(configDir, 0755)
	cuibase.CheckIfError(err)
	err = os.MkdirAll(logDir, 0755)
	cuibase.CheckIfError(err)
}

func isPathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
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

		logger.Info(filepath, "move to", trashDir)

		timestamp := strconv.FormatInt(time.Now().UnixNano(), 10)
		//logger.Debug(filepath, trashDir+"/"+filepath)
		cmd := exec.Command("mv", filepath, trashDir+"/"+filepath+".T."+timestamp)
		execCmdWithQuite(cmd)
	}
	checkTrashDaemon()
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

func checkTrashDaemon() {
	exists, err := isPathExists(pidFile)
	if exists {
		return
	}
	if err != nil {
		logger.Error(err)
		return
	}

	//启动一个子进程后主进程退出，之后的代码只有子进程会执行
	if !debug {
		_, err = xdaemon.Background(logFile, true)
		if err != nil {
			logger.Error(err)
			return
		}
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

		logger.Warn("killed")
		curDelete := atomic.AddInt32(&deleteFlag, 1)
		if curDelete == 1 {
			actualDeleteFile(pidFile)
		}
		os.Exit(1)
	}()

	err = ioutil.WriteFile(pidFile, []byte(strconv.Itoa(os.Getpid())), 0644)
	if err != nil {
		logger.Error(err)
		return
	}

	defer func() {
		curDelete := atomic.AddInt32(&deleteFlag, 1)
		if curDelete == 1 {
			actualDeleteFile(pidFile)
		}
	}()

	for true {
		current := time.Now().UnixNano()
		time.Sleep(checkPeriod)
		logger.Info("check ...")
		dir, err := ioutil.ReadDir(trashDir)
		if err != nil {
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
				logger.Info("delete: ", name[:index])
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
