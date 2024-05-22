package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/kuangcp/gobase/pkg/ctool"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"strings"
	"time"

	_ "net/http/pprof"

	"github.com/kuangcp/logger"
)

var (
	port       int
	version    bool
	serverAddr string
	localHost  string
	localAddr  string
	syncDir    string
	checkSec   int
)

var (
	lastFile   = ctool.NewSet[string]()
	sideList   = ctool.NewSet[string]() // 对端列表 格式 host:port
	configName = ".ksync.config.json"
)

type (
	ArgVO struct {
		ServerAddr string `json:"server_addr"`
		LocalHost  string `json:"local_host"`
		LocalPort  int    `json:"local_port"`
		SyncDir    string `json:"sync_dir"`
	}
)

func (a ArgVO) String() string {
	bt, err := json.Marshal(a)
	if err != nil {
		return "Error: " + err.Error()
	}
	return string(bt)
}
func init() {
	flag.IntVar(&port, "p", 8000, "port")
	flag.IntVar(&checkSec, "c", 5, "check duration second")
	flag.BoolVar(&version, "v", false, "version")
	flag.StringVar(&serverAddr, "s", "", "init server host and port. ag: 192.168.0.1:8000")
	flag.StringVar(&localHost, "l", "", "local side host. ag: 192.168.0.2")
	flag.StringVar(&syncDir, "d", "", "sync dir.")
}

func main() {
	flag.Parse()
	if version {
		fmt.Println("1.0.3")
		return
	}

	initFromConfig()
	normalizeParam()
	go pprofDebug()
	logger.Info("start success. server:", serverAddr, "local:", localHost)

	go syncTimerTask()
	go displayRemoteSide()
	webServer()
}

func displayRemoteSide() {
	for range time.NewTicker(time.Minute).C {
		logger.Info(sideList.Join(","))
	}
}
func pprofDebug() {
	debugPort := "33054"
	logger.Info("http://127.0.0.1:" + debugPort + "/debug/pprof/")
	_ = http.ListenAndServe("0.0.0.0:"+debugPort, nil)
}

// 当命令行未指定才从配置文件加载
func initFromConfig() {
	if serverAddr != "" && localHost != "" {
		return
	}

	home, err := ctool.Home()
	if err == nil {
		configName = home + "/" + configName
	}
	file, err := os.ReadFile(configName)
	if err != nil {
		logger.Error(err)
		return
	}

	var arg ArgVO
	err = json.Unmarshal(file, &arg)
	if err != nil {
		logger.Error(err)
		return
	}

	logger.Info("load config", arg)
	if arg.LocalPort != 0 {
		port = arg.LocalPort
	}
	if serverAddr == "" {
		serverAddr = arg.ServerAddr
	}
	if syncDir == "" && arg.SyncDir != "" {
		syncDir = arg.SyncDir
	}
	if localHost == "" {
		localHost = arg.LocalHost
	}
}

func normalizeParam() {
	if syncDir == "" {
		syncDir = "./"
	}
	if "windows" == runtime.GOOS {
		if !strings.HasSuffix(syncDir, "\\") {
			syncDir += "\\"
		}
	} else {
		if !strings.HasSuffix(syncDir, "/") {
			syncDir += "/"
		}
	}
}

func absPath(filename string) string {
	if "windows" == runtime.GOOS {
		return syncDir + "\\" + filename
	} else {
		return syncDir + "/" + filename
	}
}

func readNeedSyncFile() []string {
	var result []string
	dir, err := ioutil.ReadDir(syncDir)
	if err != nil {
		logger.Error(err)
		return result
	}

	firstInit := lastFile.IsEmpty()
	for _, info := range dir {
		if info.IsDir() {
			continue
		}
		fileName := info.Name()
		if strings.HasPrefix(fileName, ".") {
			continue
		}
		if fileName == configName {
			continue
		}
		contains := lastFile.Contains(fileName)
		lastFile.Add(fileName)

		if !contains || firstInit {
			logger.Info("need sync", fileName, info.ModTime())
			result = append(result, fileName)
		}
	}
	return result
}

func isFileExist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

// 客户端模式：注册自身到服务端
func registerOnServer() {
	if serverAddr == "" {
		return
	}
	// 必须要声明自身可被访问的host
	if localHost == "" {
		logger.Warn("local host is missing")
		os.Exit(1)
	}

	req, _ := http.NewRequest("GET", "http://"+serverAddr+"/register", nil)
	req.Header.Set("self", localAddr)
	rsp, err := http.DefaultClient.Do(req)
	if rsp != nil && rsp.Body != nil {
		defer rsp.Body.Close()
	}
	if err != nil {
		logger.Error(err, rsp)
		return
	}
	sideList.Add(serverAddr)
}

// 检查本地可访问性
func checkLocalHostReachable(localAddr string) {
	resp, err := http.DefaultClient.Get("http://" + localAddr + "/ping")
	if err != nil {
		logger.Warn("local host unreachable", localAddr, err)
		os.Exit(1)
	}

	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
		all, err := io.ReadAll(resp.Body)
		if err != nil {
			logger.Warn("local host unreachable", localAddr, err)
			os.Exit(1)
		}
		respStr := string(all)
		if respStr != pong {
			logger.Warn("un except local host", localAddr, respStr)
			os.Exit(1)
		}
	}
}

func syncTimerTask() {
	localAddr = fmt.Sprintf("%v:%v", localHost, port)
	checkLocalHostReachable(localAddr)
	registerOnServer()
	ticker := time.NewTicker(time.Second * time.Duration(checkSec))
	for range ticker.C {
		registerOnServer()
		syncFile()
	}
}

func syncFile() {
	if sideList.Len() == 0 {
		return
	}

	fileList := readNeedSyncFile()
	if len(fileList) == 0 {
		return
	}
	logger.Info("check sync %v", sideList)
	for _, path := range fileList {
		sideList.Loop(func(addr string) {
			postFile(addr, path)
		})
	}
}

func postFile(server string, path string) {
	name := url.QueryEscape(path)
	existURL := "http://" + server + "/exist?name=" + name

	existRsp, err := http.Get(existURL)
	if err != nil || existRsp == nil || existRsp.Body == nil {
		return
	}
	rspStr, err := ioutil.ReadAll(existRsp.Body)
	if err != nil || string(rspStr) == "EXIST" {
		logger.Info("%s exist", name)
		return
	}
	defer existRsp.Body.Close()

	open, err := os.Open(absPath(path))
	if err != nil {
		logger.Error(err)
		return
	}
	defer open.Close()

	postRsp, err := http.Post("http://"+server+"/upload?name="+name, "", open)
	if err != nil || postRsp == nil || postRsp.Body == nil {
		logger.Error(err)
		return
	}
	defer postRsp.Body.Close()
	logger.Info("send to ", server, name)
}
