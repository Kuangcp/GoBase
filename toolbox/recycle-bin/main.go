package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"os/signal"
	"sort"
	"strconv"
	"strings"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/kuangcp/gobase/pkg/cuibase"
	"github.com/kuangcp/logger"
	"github.com/manifoldco/promptui"
)

func main() {
	invokeWithCondition(help, info.PrintHelp)
	invokeWithCondition(initConfig, InitConfig)
	invokeWithCondition(listTrash, ListTrashFiles)
	invokeWithCondition(exit, ExitCheckFileDaemon)
	invokeWithCondition(log, func() {
		fmt.Println(logFile)
	})
	invokeWithCondition(showConfig, func() {
		fmt.Println(configFile)
	})

	invokeWithCondition(illegalQuit, func() {
		ActualDeleteFile(pidFile)
	})
	invokeWithStr(suffix, func(s string) {
		DeleteFileBySuffix(strings.Split(s, ","))
	})
	invokeWithStr(restore, RestoreFile)

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

func RestoreFile(restoreFile string) {
	items := listTrashFileItem(func(val string) bool {
		return strings.Contains(val, restoreFile)
	})

	length := len(items)
	if length == 0 {
		logger.Info("Not match: " + restoreFile)
	} else if length == 1 {
		restoreFileToCurDir(items[0])
	} else {
		file, err := SelectFile(items)
		if err != nil {
			logger.Error(err)
			return
		}

		restoreFileToCurDir(*file)
	}
}

func SelectFile(list []fileItem) (*fileItem, error) {
	type option struct {
		Name    string
		Time    string
		Peppers int
	}

	var peppers []option
	for i, item := range list {
		peppers = append(peppers, option{
			Name:    item.name,
			Time:    item.formatTime(),
			Peppers: i,
		})
	}

	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}?",
		Active:   "✔️{{ .Name | green }} {{ .Time }}",
		Inactive: "  {{ .Name | cyan }} {{ .Time }}",
		Selected: "✔️{{ .Name | green | cyan }}",
	}

	searcher := func(input string, index int) bool {
		pepper := peppers[index]
		name := strings.Replace(strings.ToLower(pepper.Name), " ", "", -1)
		input = strings.Replace(strings.ToLower(input), " ", "", -1)

		return strings.Contains(name, input)
	}

	promptSelect := promptui.Select{
		Label:     "Restore which file",
		Items:     peppers,
		Templates: templates,
		Size:      4,
		Searcher:  searcher,
	}

	i, _, err := promptSelect.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return nil, err
	}

	return &list[i], nil
}

func restoreFileToCurDir(item fileItem) {
	logger.Warn("restore ", item.file.Name())
	cmd := exec.Command("mv", trashDir+"/"+item.file.Name(), item.name)
	execCmdWithQuite(cmd)
}

func listTrashFileItem(filter func(string) bool) []fileItem {
	var result []fileItem
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

		result = append(result, fileItem{
			name:      filename,
			timestamp: parseInt,
			file:      fileInfo,
		})
	}
	return result
}

func ListTrashFiles() {
	err := parseTime()
	if err != nil {
		logger.Error(err)
		os.Exit(1)
	}

	items := listTrashFileItem(func(s string) bool {
		return true
	})
	current := time.Now().UnixNano()
	if len(items) != 0 {
		fmt.Printf("%v%-23s %-10s %s%v\n", cuibase.Cyan, "DeleteTime", "Remaining", "File", cuibase.End)
	}

	if listOrder > 0 {
		sort.Slice(items, func(i, j int) bool {
			if listOrder != 1 {
				return items[i].timestamp > items[j].timestamp
			}
			return items[i].timestamp < items[j].timestamp
		})
		for _, item := range items {
			fmt.Print(item.formatForList(current))
		}
	} else {
		for _, item := range items {
			fmt.Print(item.formatForList(current))
		}
	}
}

func fmtDuration(d time.Duration) string {
	d = d.Round(time.Second)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	d -= m * time.Minute
	s := d / time.Second
	return fmt.Sprintf("%03d:%02d:%02d", h, m, s)
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
	logger.Info("Start check daemon. check:", checkPeriod, "retention:", retentionTime, "pid:", os.Getpid())

	go func() {
		// Wait for interrupt signal to gracefully shutdown the server with a timeout of 5 seconds.
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
	doCheckTrashDir()
}

func doCheckTrashDir() {
	emptyCount := 0
	for true {
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

		cleanFile(dir)
	}
}

func cleanFile(dir []os.FileInfo) {
	current := time.Now().UnixNano()
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
		if current-parseInt <= retentionTime.Nanoseconds() {
			continue
		}

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
		if err != nil {
			logger.Error(err)
			os.Exit(1)
		}
		if !exists {
			logger.Error(filepath, "not found")
			return
		}

		finalName, err := buildTrashFileName(filepath)
		if err != nil {
			continue
		}

		logger.Info("Prepare delete:", filepath)
		//logger.Debug(filepath, trashDir+"/"+finalName)
		cmd := exec.Command("mv", filepath, finalName)
		execCmdWithQuite(cmd)
	}
}

func buildTrashFileName(path string) (string, error) {
	var result = path
	if strings.HasSuffix(result, "//") || strings.HasPrefix(result, "//") {
		logger.Error("Invalid path", path)
		return "", fmt.Errorf("")
	}

	if strings.HasSuffix(result, "/") {
		result = result[:len(result)-1]
	}

	if strings.HasPrefix(result, "/") {
		last := strings.LastIndex(result, "/")
		if last != -1 && last < len(result)-1 {
			result = result[last+1:]
		} else {
			logger.Error("Invalid path", path)
			return "", fmt.Errorf("")
		}
	}

	if result == "/" || result == "" {
		logger.Error("Invalid path", path)
		return "", fmt.Errorf("")
	}

	if strings.Contains(result, "/") {
		last := strings.LastIndex(result, "/")
		result = result[last+1:]
	}
	timestamp := strconv.FormatInt(time.Now().UnixNano(), 10)
	return trashDir + "/" + result + ".T." + timestamp, nil
}

func DeleteFileBySuffix(params []string) {
	if len(params) == 0 {
		return
	}

	pwd, err := os.Getwd()
	if err != nil {
		logger.Error(err)
		os.Exit(1)
	}
	dir, err := ioutil.ReadDir(pwd)
	if err != nil {
		logger.Error(err)
		os.Exit(1)
	}

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
	if err != nil {
		logger.Error(err)
		os.Exit(1)
	}
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
	if err != nil {
		logger.Error(err)
		os.Exit(1)
	}

	// TODO may kill error pid
	pid := string(file)
	logger.Info("kill ", pid)
	cmd := exec.Command("kill", pid)
	execCmdWithQuite(cmd)
}

func ActualDeleteFile(path string) {
	err := os.Remove(path)
	if err != nil {
		logger.Error(path, err)
		os.Exit(1)
	}
}

// 静默执行 不关心返回值
func execCmdWithQuite(cmd *exec.Cmd) {
	var out bytes.Buffer

	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		logger.Error(cmd, err)
		os.Exit(1)
	}
}
