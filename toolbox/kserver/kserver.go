package main

import (
	"flag"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/kuangcp/gobase/cuibase"
)

var (
	help    bool
	port    int
	success = []byte("ok\n")
)

var info = cuibase.HelpInfo{
	Description:   "Start simple http server on current path",
	Version:       "1.0.4",
	SingleFlagLen: -3,
	ValueLen:      -7,
	Flags: []cuibase.ParamVO{
		{Short: "-h", Comment: "help"},
	},
	Options: []cuibase.ParamVO{
		{Short: "-p", Value: "<port>", Comment: "web server port"},
	}}

func init() {
	flag.BoolVar(&help, "h", false, "")
	flag.IntVar(&port, "p", 8989, "")
}

func getInternalIP() string {
	address, err := net.InterfaceAddrs()
	if err != nil {
		log.Fatal(err.Error())
	}

	for _, addr := range address {
		if ipNet, ok := addr.(*net.IPNet); ok &&
			!ipNet.IP.IsLoopback() &&
			ipNet.IP.To4() != nil &&
			strings.HasPrefix(ipNet.IP.String(), "192") {
			return ipNet.IP.String()
		}
	}
	return ""
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	// 接受文件
	file, header, err := r.FormFile("file")
	if err != nil || header == nil {
		// ignore the error handler
		log.Println(err)
		return
	}

	log.Printf("upload: %s", header.Filename)
	// 将文件拷贝到指定路径下，或者其他文件操作
	dst, err := os.Create(header.Filename)
	if err != nil {
		// ignore
		log.Println(err)
		return
	}

	_, err = io.Copy(dst, file)
	if err != nil {
		// ignore
		log.Println(err)
	}
	w.Write(success)
}

func startWebServer(port int) {
	internalIP := getInternalIP()

	// 绑定路由到当前目录
	fs := http.FileServer(http.Dir("./"))
	http.Handle("/", http.StripPrefix("/", fs))

	http.HandleFunc("/up", uploadHandler)

	http.HandleFunc("/echo", func(resp http.ResponseWriter, req *http.Request) {
		body, _ := ioutil.ReadAll(req.Body)
		log.Printf(string(body))
		resp.Write(success)
	})

	log.Printf("web server started.\n")
	log.Printf("%vhttp://127.0.0.1:%v%v\n", cuibase.Green, port,cuibase.End)
	log.Printf("%vhttp://%v:%v%v\n", cuibase.Green, internalIP, port,cuibase.End)
	log.Printf("%v/up%v   : upload file.    | curl -X POST -H 'Content-Type: multipart/form-data' -F 'file=@index.html' http://127.0.0.1:%v/up\n",
		cuibase.Purple, cuibase.End, port)
	log.Printf("%v/echo%v : echo. | curl -d 'hi' http://127.0.0.1:8989/echo\n",
		cuibase.Purple, cuibase.End)
	err := http.ListenAndServe(":"+strconv.Itoa(port), nil)
	if err != nil {
		log.Fatal("error: ", err)
	}
}

func main() {
	flag.Parse()
	if help {
		info.PrintHelp()
		return
	}

	if port > 65535 || port == 0 {
		log.Fatalf("Please input correct port [1, 65535]. input:%v", port)
	}
	if port < 1024 {
		log.Print("Please run by root.")
	}

	startWebServer(port)
}
