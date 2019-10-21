package main

import (
	"fmt"
	"net"
	"log"
	"net/http"
	"os"
	"strings"
	"strconv"
)

var green = "\033[0;32m"
var yellow = "\033[0;33m"
var end = "\033[0m"

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

func getInternalIp()(string){
	addrs, err := net.InterfaceAddrs()
    if err != nil {
        log.Fatal(err.Error())
    }

	var internalIp=""
    for _, addr := range addrs {
        if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
            if ipnet.IP.To4() != nil {
				ip := ipnet.IP.String()
				if strings.HasPrefix(ip, "192"){
					internalIp = ip
				}
            }
        }
	}
	return internalIp
}

func main() {
	portStr := readPortByParam()
	internalIp := getInternalIp()

	http.HandleFunc("/index", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hi Mythos !")
	})

	// 绑定路由到当前目录
	fs := http.FileServer(http.Dir("./"))
	http.Handle("/", http.StripPrefix("/", fs))

	log.Printf("%v Start webserver success on http://127.0.0.1:%v %vhttp://%v:%v %v ",
                 green, portStr, yellow, internalIp, portStr, end)
	err := http.ListenAndServe(":"+portStr, nil)
	if err != nil {
		log.Fatal("error: ", err)
	}
}
