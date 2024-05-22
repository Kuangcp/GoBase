package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/kuangcp/logger"
)

const (
	pong = "Pong"
)

func ping(writer http.ResponseWriter, request *http.Request) {
	writer.Write([]byte(pong))
}

func exist(writer http.ResponseWriter, request *http.Request) {
	name := request.URL.Query().Get("name")
	unescape, err := url.QueryUnescape(name)
	if err != nil {
		logger.Error(err)
		return
	}
	exist := isFileExist(absPath(unescape))
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

	open, err := os.Create(absPath(unescape))
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

func register(writer http.ResponseWriter, request *http.Request) {
	client := request.Header.Get("self")
	if sideList.Contains(client) {
		//logger.Warn("already register", client)
		writer.Write([]byte("EXIST"))
		return
	}
	logger.Info("register new", client)
	sideList.Add(client)
	writer.Write([]byte("OK"))
}

func webServer() {
	http.HandleFunc("/ping", ping)
	http.HandleFunc("/exist", exist)
	// 接收文件
	http.HandleFunc("/upload", upload)
	// 注册
	http.HandleFunc("/register", register)
	// 触发检查文件同步
	//http.HandleFunc("/refresh", notifySyncFile)

	logger.Info("listening on :%v", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		log.Fatal("error: ", err)
	}
}
