package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/kuangcp/gobase/pkg/cuibase"
	"github.com/kuangcp/logger"
)

var (
	port       int
	version    bool
	serverAddr string
	localAddr  string
	localHost  string
	syncDir    string
	checkSec   int
)

var (
	lastFile = cuibase.NewSet()
	sideList = cuibase.NewSet() // 对端列表 格式 host:port
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
		fmt.Println("1.0.2")
		return
	}

	normalizeParam()

	go func() {
		for range time.NewTicker(time.Second * 5).C {
			registerOnServer()
		}
	}()

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

	firstInit := lastFile.IsEmpty()
	for _, info := range dir {
		if info.IsDir() {
			continue
		}
		contains := lastFile.Contains(info.Name())
		lastFile.Add(info.Name())

		if !contains || firstInit {
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
	sideList.Add(serverAddr)
}

func syncTimerTask() {
	ticker := time.NewTicker(time.Second * time.Duration(checkSec))
	for range ticker.C {
		syncFile()
	}
}

func syncFile() {
	if sideList.Len() == 0 {
		return
	}

	//logger.Info("check sync %v", sideList)
	fileList := readNeedSyncFile()
	for _, path := range fileList {
		sideList.Loop(func(i interface{}) {
			postFile(i.(string), path)
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

	open, err := os.Open(path)
	if err != nil {
		logger.Error(err)
		return
	}
	defer open.Close()

	post, err := http.Post("http://"+server+"/upload?name="+name, "", open)
	if err != nil || post == nil || post.Body == nil {
		return
	}
	defer post.Body.Close()
	logger.Info("send to ", server, name)
}
