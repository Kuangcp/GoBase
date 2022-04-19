package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/kuangcp/logger"
)

func exist(writer http.ResponseWriter, request *http.Request) {
	name := request.URL.Query().Get("name")
	unescape, err := url.QueryUnescape(name)
	if err != nil {
		logger.Error(err)
		return
	}
	exist := isFileExist(syncDir + "" + unescape)
	if exist {
		writer.Write([]byte("EXIST"))
	} else {
		writer.Write([]byte("NONE"))
	}
}

func upload(writer http.ResponseWriter, request *http.Request) {
	name := request.URL.Query().Get("name")
	unescape, err := url.QueryUnescape(name)
	if err != nil {
		logger.Error(err)
		return
	}

	open, err := os.Create(syncDir + unescape)
	if err != nil {
		logger.Error(err)
		return
	}
	defer open.Close()

	var buf = make([]byte, 4096)
	for {
		read, err := request.Body.Read(buf)
		if read != 0 {
			open.Write(buf[:read])
		}
		if err != nil {
			break
		}
	}

	writer.Write([]byte(""))
}

//func notifySyncFile(writer http.ResponseWriter, request *http.Request) {
//	syncFile()
//	actionList := request.URL.Query().Get("actionList")
//	if strings.Contains(actionList, localAddr) {
//		//logger.Warn("cycle notify", actionList)
//		writer.Write([]byte(""))
//		return
//	}
//
//	sideList.Loop(func(s interface{}) {
//		logger.Info("notify refresh", s)
//		http.Get("http://" + s.(string) + "/refresh?actionList=" + actionList + "," + localAddr)
//	})
//	writer.Write([]byte("OK"))
//}

func register(writer http.ResponseWriter, request *http.Request) {
	client := request.Header.Get("self")
	logger.Info("register new", client)
	if sideList.Contains(client) {
		logger.Warn("already register")
		writer.Write([]byte("EXIST"))
		return
	}
	sideList.Add(client)
	writer.Write([]byte("OK"))
}

func webServer() {

	http.HandleFunc("/exist", exist)
	// 接收文件
	http.HandleFunc("/upload", upload)
	// 注册
	http.HandleFunc("/register", register)
	// 触发检查文件同步
	//http.HandleFunc("/refresh", notifySyncFile)

	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		log.Fatal("error: ", err)
	}
}
