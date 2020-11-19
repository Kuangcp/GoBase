package main

import (
	"bytes"
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
	configFile    string
	pidFile       string
	retentionTime time.Duration
	checkPeriod   time.Duration
)

type (
	FileItem struct {
		name      string
		timestamp int64
		file      os.FileInfo
	}
)

func invokeWithCondition(flag bool, action func()) {
	if flag {
		action()
		os.Exit(0)
	}
}

func main() {
	invokeWithCondition(help, info.PrintHelp)
	invokeWithCondition(initConfig, InitConfig)
	invokeWithCondition(listTrash, ListTrashFiles)
	invokeWithCondition(exit, ExitCheckFileDaemon)
	invokeWithCondition(log, PrintLogFile)
	invokeWithCondition(illegalQuit, func() {
		ActualDeleteFile(pidFile)
	})

	if suffix != "" {
		DeleteFileBySuffix(strings.Split(suffix, ","))
		return
	}

	if restore != "" {
		RestoreFile(restore)
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

	args := os.Args
	if len(args) == 1 {
		info.PrintHelp()
	} else {
		DeleteFiles(args[1:])
		CheckWithDaemon()
	}
}

func PrintLogFile() {
	fmt.Println(logFile)
}

func RestoreFile(restoreFile string) {
	items := listFileItem(func(val string) bool {
		return strings.Contains(val, restoreFile)
	})
	length := len(items)
	if length == 0 {
		logger.Info("Not match: " + restoreFile)
	} else if length == 1 {
		restoreFileToCurDir(items[0])
	} else {
		for i := range items {
			fmt.Printf("  %s : %s\n", cuibase.Green.Printf("%4s", strconv.Itoa(i)), items[i].name)
		}
		fmt.Printf("Select one: ")
		selectFile := 0
		_, err := fmt.Scanln(&selectFile)
		if err != nil {
			logger.Error(err)
			return
		}
		if selectFile >= length {
			logger.Error("Out of index")
			return
		}
		restoreFileToCurDir(items[selectFile])
	}
}

func restoreFileToCurDir(item FileItem) {
	logger.Warn("restore ", item.file.Name())
	cmd := exec.Command("mv", trashDir+"/"+item.file.Name(), item.name)
	execCmdWithQuite(cmd)
}

func listFileItem(filter func(string) bool) []FileItem {
	var result []FileItem
	dir, err := ioutil.ReadDir(trashDir)
	if err != nil {
		logger.Error(err)
		return result
	}

	for _, fileInfo := range dir {
		name := fileInfo.Name()
		if !filter(name) {
			continue
		}

		index := strings.Index(name, ".T.")
		if index == -1 {
			continue
		}

		filename := name[:index]
		value := name[index+3:]
		parseInt, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			logger.Error(err)
			continue
		}

		result = append(result, FileItem{
			name:      filename,
			timestamp: parseInt,
			file:      fileInfo,
		})
	}
	return result
}

func ListTrashFiles() {
	err := parseTime()
	checkError(err)

	items := listFileItem(func(s string) bool {
		return true
	})
	current := time.Now().UnixNano()
	if len(items) != 0 {
		fmt.Printf("%-23s %-10s %s\n", "DeleteTime", "Remaining", "File")
	}

	for _, item := range items {
		second := strconv.FormatInt((retentionTime.Nanoseconds()-current+item.timestamp)/1000000000, 10)
		duration, err := time.ParseDuration(second + "s")
		if err != nil {
			duration = 0
		}
		fmt.Println(time.Unix(item.timestamp/1000000000, 0).Format("2006-01-02 15:04:05.000"),
			cuibase.Yellow.Printf("%10s", duration.String()), cuibase.Green.Print(item.name))
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
	if err != nil || exists {
		return
	}

	// avoid repeat delete
	var deleteFlag int32 = 0
	logger.Info("Start check trash, period:", checkPeriod, "pid:", os.Getpid())

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
	return
}

func parseTime() error {
	duration, err := time.ParseDuration(retentionStr)
	if err != nil {
		return err
	}

	retentionTime = duration
	checkPeriod, err = time.ParseDuration(checkStr)
	if err != nil {
		return err
	}
	return nil
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
	if files == nil || len(files) == 0 {
		return
	}
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

func checkError(err error) {
	if err != nil {
		logger.Error(err)
		os.Exit(1)
	}
}

func DeleteFileBySuffix(params []string) {
	if len(params) == 0 {
		return
	}

	pwd, err := os.Getwd()
	checkError(err)
	dir, err := ioutil.ReadDir(pwd)
	checkError(err)

	var files []string
	for _, fileInfo := range dir {
		name := fileInfo.Name()
		for _, suffix := range params {
			if strings.HasSuffix(name, "."+suffix) {
				files = append(files, name)
			}
		}
	}

	fmt.Printf("%v\nDelete the above file? (y/n):", files)
	var input string
	_, err = fmt.Scanf("%s", &input)
	checkError(err)
	if input == "y" {
		DeleteFiles(files)
	}
}

func ExitCheckFileDaemon() {
	exists, err := isPathExists(pidFile)
	if !exists {
		logger.Error("no pid file")
		return
	}
	file, err := ioutil.ReadFile(pidFile)
	checkError(err)

	// TODO may kill error pid
	pid := string(file)
	logger.Info("kill ", pid)
	cmd := exec.Command("kill", pid)
	execCmdWithQuite(cmd)
}

func ActualDeleteFile(path string) {
	err := os.Remove(path)
	checkError(err)
}

// 静默执行 不关心返回值
func execCmdWithQuite(cmd *exec.Cmd) {
	var out bytes.Buffer

	cmd.Stdout = &out
	err := cmd.Run()
	checkError(err)
}
