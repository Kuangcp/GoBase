package main

import (
	"github.com/kuangcp/gobase/cuibase"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func readPortByParam() string {
	param := os.Args
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

func main() {
	portStr := readPortByParam()
	internalIp := getInternalIp()

	// 绑定路由到当前目录
	fs := http.FileServer(http.Dir("./"))
	http.Handle("/", http.StripPrefix("/", fs))

	log.Printf("%v Start webserver success on http://127.0.0.1:%v %vhttp://%v:%v %v ",
		cuibase.Green, portStr, cuibase.Yellow, internalIp, portStr, cuibase.End)
	err := http.ListenAndServe(":"+portStr, nil)
	if err != nil {
		log.Fatal("error: ", err)
	}
}
