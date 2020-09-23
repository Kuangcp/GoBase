package main

import (
	"flag"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"

	"github.com/kuangcp/gobase/cuibase"
)

var (
	help bool
	port int
)
var info = cuibase.HelpInfo{
	Description:   "Start simple http server on current path",
	Version:       "1.0.3",
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

func getInternalIp() string {
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

func startWebServer(port int) {
	internalIp := getInternalIp()

	// 绑定路由到当前目录
	fs := http.FileServer(http.Dir("./"))
	http.Handle("/", http.StripPrefix("/", fs))

	log.Printf("Start webserver success. %vhttp://127.0.0.1:%v %vhttp://%v:%v %v ",
		cuibase.Green, port, cuibase.Yellow, internalIp, port, cuibase.End)
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
