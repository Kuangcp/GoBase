package main

import (
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"

	"github.com/kuangcp/gobase/cuibase"
)

func help(params []string) {
	cuibase.PrintTitleDefault("Start simple http server on current path")
	format := cuibase.BuildFormat(-5, -15)
	cuibase.PrintParams(format, []cuibase.ParamInfo{
		{
			Verb:    "-h",
			Param:   "",
			Comment: "help",
		}, {
			Verb:    "-p",
			Param:   "port",
			Comment: "specific port",
		},
	})
}

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
		if ipnet, ok := addr.(*net.IPNet); ok &&
			!ipnet.IP.IsLoopback() &&
			ipnet.IP.To4() != nil &&
			strings.HasPrefix(ipnet.IP.String(), "192") {
			return ipnet.IP.String()
		}
	}
	return ""
}

func runWithPort(params []string) {
	portStr := readPortByParam(params[1:])
	run(portStr)
}
func runWithDefaultPort(params []string) {
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
	cuibase.RunAction(map[string]func(params []string){
		"-h": help,
		"-p": runWithPort,
	}, runWithDefaultPort)
}
