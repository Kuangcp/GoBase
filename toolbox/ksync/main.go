package main

import (
	"flag"
	"fmt"
	"github.com/kuangcp/logger"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"strings"
	"time"
)

var (
	sideList   []string // 对端列表 格式 host:port
	port       int
	version    bool
	serverAddr string
	localAddr  string
	localHost  string
	syncDir    string
	checkSec   int
)

var (
	lastFile = make(map[string]struct{})
)

func init() {
	flag.IntVar(&port, "p", 8000, "port")
	flag.IntVar(&checkSec, "c", 2, "check duration second")
	flag.BoolVar(&version, "v", false, "version")
	flag.StringVar(&serverAddr, "s", "", "init server host&port. ag: 192.168.0.1:8000")
	flag.StringVar(&localHost, "l", "", "local side host. ag: 192.168.0.2")
	flag.StringVar(&syncDir, "d", "./", "sync dir.")
}

func main() {
	flag.Parse()
	if version {
		fmt.Println("1.0.0")
		return
	}

	normalizeParam()
	registerOnServer()
	go syncTimerTask()
	webServer()
}

func normalizeParam() {
	localAddr = localHost + ":" + fmt.Sprint(port)
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

func readNeedSyncFile() []string {
	var result []string
	dir, err := ioutil.ReadDir(syncDir)
	if err != nil {
		fmt.Println(err)
		return result
	}

	init := len(lastFile) == 0
	for _, info := range dir {
		if info.IsDir() {
			continue
		}
		_, ok := lastFile[info.Name()]
		lastFile[info.Name()] = struct{}{}

		if !ok || init {
			logger.Info("need sync", info.Name(), info.ModTime())
			result = append(result, info.Name())
		}
	}
	return result
}

func isFileExist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

// 注册到服务端
func registerOnServer() {
	// 如果需要注册到服务端，必须要声明自己可被访问的host
	if serverAddr != "" && localHost == "" {
		logger.Warn("local host is missing")
		os.Exit(1)
	}

	if serverAddr == "" {
		return
	}

	client := http.Client{}
	req, err := http.NewRequest("GET", "http://"+serverAddr+"/register", nil)
	req.Header.Set("self", fmt.Sprintf("%v:%v", localHost, port))
	if err != nil {
		fmt.Println(err)
		return
	}
	rsp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(rsp)
	sideList = append(sideList, serverAddr)
}

func syncTimerTask() {
	ticker := time.NewTicker(time.Second * time.Duration(checkSec))
	for range ticker.C {
		syncFile()
	}
}

func syncFile() {
	if len(sideList) == 0 {
		return
	}

	//logger.Info("check sync %v", sideList)
	fileList := readNeedSyncFile()
	for _, path := range fileList {
		for _, r := range sideList {
			postFile(r, path)
		}
	}
}

func postFile(server string, path string) {
	name := url.QueryEscape(path)
	existURL := "http://" + server + "/exist?name=" + name

	existRsp, err := http.Get(existURL)
	if existRsp == nil || err != nil {
		return
	}
	rspStr, err := ioutil.ReadAll(existRsp.Body)
	if string(rspStr) == "EXIST" {
		logger.Info("%s exist", name)
		return
	}

	open, err := os.Open(path)
	if err != nil {
		logger.Error(err)
		return
	}

	defer open.Close()

	syncURL := "http://" + server + "/upload?name=" + name
	post, err := http.Post(syncURL, "", open)
	if err != nil {
		return
	}

	defer post.Body.Close()
	logger.Info("send to ", server, name)
}
