package main

import (
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"

	"github.com/kuangcp/gobase/cuibase"
)

var info = cuibase.HelpInfo{
	Description: "Start simple http server on current path",
	Version:     "1.0.1",
	VerbLen:     -5,
	ParamLen:    -5,
	Params: []cuibase.ParamInfo{
		{
			Verb:    "-h",
			Param:   "",
			Comment: "help",
		}, {
			Verb:    "-p",
			Param:   "port",
			Comment: "specific port",
			Handler: RunWithPort,
		},
	}}

func readPortByParam(param []string) string {
	var port = 8099
	if len(param) > 1 {
		// so happy error handle
		temp, err := strconv.Atoi(param[1])
		if err != nil {
			log.Print("Please input correct port [1, 65535].")
			log.Fatal("input:"+param[1]+". ", err)
		}
		port = temp
	}
	portStr := strconv.Itoa(port)
	if port > 65535 || port == 0 {
		log.Fatal("Please input correct port [1, 65535]. input:" + portStr)
	}
	if port < 1024 {
		log.Print("Please run by root ")
	}
	return portStr
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

func RunWithPort(params []string) {
	portStr := readPortByParam(params[1:])
	run(portStr)
}
func RunWithDefaultPort(params []string) {
	run("8889")
}

func run(port string) {
	internalIp := getInternalIp()

	// 绑定路由到当前目录
	fs := http.FileServer(http.Dir("./"))
	http.Handle("/", http.StripPrefix("/", fs))

	log.Printf("%v Start webserver success on http://127.0.0.1:%v %vhttp://%v:%v %v ",
		cuibase.Green, port, cuibase.Yellow, internalIp, port, cuibase.End)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal("error: ", err)
	}
}

func main() {
	cuibase.RunActionFromInfo(info, RunWithDefaultPort)
}
