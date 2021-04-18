package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strings"

	_ "embed"

	"github.com/kuangcp/gobase/cuibase"
)

//go:embed up.html
var uploadStaticPage string

var (
	help         bool
	port         int
	buildVersion string
)

var info = cuibase.HelpInfo{
	Description:   "Start static file web server on current path",
	Version:       "1.0.6",
	BuildVersion:  buildVersion,
	SingleFlagLen: -2,
	ValueLen:      -6,
	Flags: []cuibase.ParamVO{
		{Short: "-h", BoolVar: &help, Comment: "help"},
	},
	Options: []cuibase.ParamVO{
		{Short: "-p", Value: "port", Comment: "web server port"},
	}}

func init() {
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
	var maxMib int64 = 10
	err := r.ParseMultipartForm(maxMib << 20)
	if err != nil {
		log.Println(err)
	}
	for _, headers := range r.MultipartForm.File {
		for _, header := range headers {
			log.Printf("upload: %s", header.Filename)
			// 将文件拷贝到指定路径下，或者其他文件操作
			dst, err := os.Create(header.Filename)
			if err != nil {
				// ignore
				log.Println(err)
				return
			}

			open, _ := header.Open()
			_, err = io.Copy(dst, open)
			if err != nil {
				// ignore
				log.Println(err)
			}
		}
	}

	http.Redirect(w, r, "/up", http.StatusMovedPermanently)
}

func startWebServer(port int) {
	internalIP := getInternalIP()

	// 绑定路由到当前目录
	fs := http.FileServer(http.Dir("./"))

	http.Handle("/", http.StripPrefix("/", fs))

	http.HandleFunc("/f", uploadHandler)

	http.HandleFunc("/up", func(writer http.ResponseWriter, request *http.Request) {
		_, err := writer.Write([]byte(uploadStaticPage))
		if err != nil {
			log.Println(err)
		}
	})

	http.HandleFunc("/e", func(writer http.ResponseWriter, request *http.Request) {
		body, _ := ioutil.ReadAll(request.Body)
		log.Printf(string(body))
	})

	// log
	innerURL := fmt.Sprintf("http://%v:%v", internalIP, port)
	log.Printf("static file web server has started.\n")
	log.Printf("%vhttp://127.0.0.1:%v%v\n", cuibase.Green, port, cuibase.End)
	log.Printf("%v%v%v\n", cuibase.Green, innerURL, cuibase.End)
	log.Printf("%v/up%v  %v/up\n", cuibase.Purple, cuibase.End, innerURL)
	log.Printf("%v/f%v   curl -X POST -H 'Content-Type: multipart/form-data' %v/f -F file=@index.html\n",
		cuibase.Purple, cuibase.End, innerURL)
	log.Printf("%v/e%v   curl %v/e -d 'echo hi'\n", cuibase.Purple, cuibase.End, innerURL)

	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		log.Fatal("error: ", err)
	}
}

func main() {
	info.Parse()
	if help {
		info.PrintHelp()
		return
	}

	if port > 65535 || port == 0 {
		log.Fatalf("Please input correct port [1, 65535]. now: %v", port)
	}
	if port < 1024 {
		log.Printf("%vWARN: [1-1024] need run by root user.%v", cuibase.Red, cuibase.End)
	}

	startWebServer(port)
}
