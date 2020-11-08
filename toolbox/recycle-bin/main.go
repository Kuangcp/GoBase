package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
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
	daemon   bool
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
		}, {
			Short:   "-d",
			Value:   "",
			Comment: "Start check by daemon",
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
	flag.BoolVar(&daemon, "d", false, "")
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

	args := os.Args
	if len(args) == 1 {
		info.PrintHelp()
	} else {
		deleteFiles(args[1:])
		checkWithDaemon()
	}
}

func checkWithDaemon() {
	params := fmt.Sprintf(" -C %s -l %s", checkStr, liveStr)
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

func NewSysProcAttr() *syscall.SysProcAttr {
	return &syscall.SysProcAttr{
		Setsid: true,
	}
}

func startProc(args []string, logFile string) (*exec.Cmd, error) {
	cmd := &exec.Cmd{
		Path:        args[0],
		Args:        args,
		SysProcAttr: NewSysProcAttr(),
	}

	if logFile != "" {
		stdout, err := os.OpenFile(logFile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
		if err != nil {
			log.Println(os.Getpid(), ": 打开日志文件错误:", err)
			return nil, err
		}
		cmd.Stderr = stdout
		cmd.Stdout = stdout
	}

	err := cmd.Start()
	if err != nil {
		return nil, err
	}

	return cmd, nil
}

func checkTrashDaemon() {
	exists, err := isPathExists(pidFile)
	if exists {
		return
	}
	logger.Debug("Start check trash")
	if err != nil {
		logger.Error(err)
		return
	}

	// avoid repeat delete
	var deleteFlag int32 = 0

	logger.Debug("process", os.Getpid())
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

	// create pid
	err = ioutil.WriteFile(pidFile, []byte(strconv.Itoa(os.Getpid())), 0644)
	if err != nil {
		logger.Error(err)
		return
	}

	// delete pid
	defer func() {
		curDelete := atomic.AddInt32(&deleteFlag, 1)
		if curDelete == 1 {
			actualDeleteFile(pidFile)
		}
	}()

	for true {
		current := time.Now().UnixNano()
		time.Sleep(checkPeriod)
		logger.Info("check ...", checkPeriod, os.Getpid())
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
