package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"github.com/kuangcp/logger"
	"net/http"
)

var (
	port        int
	reloadConf  bool
	debug       bool
	queryServer bool
	queryPort   int
)

func startQueryServer() {
	logger.Info("Start query server on 127.0.0.1:%d", queryPort)

	http.HandleFunc("/", searchPage)
	http.HandleFunc("/list", handleInterceptor(JSONFunc(pageListReqHistory)))
	http.HandleFunc("/curl", buildCurlCommandApi)
	http.HandleFunc("/replay", replayRequest)
	http.HandleFunc("/del", delRequest)
	http.HandleFunc("/flushAll", flushAllData)

	http.ListenAndServe(fmt.Sprintf(":%v", queryPort), nil)
}

func main() {
	flag.IntVar(&port, "p", 1234, "port")
	flag.BoolVar(&reloadConf, "r", false, "auto reload changed config")
	flag.BoolVar(&debug, "d", false, "debug mode")

	flag.BoolVar(&queryServer, "q", false, "query log")
	flag.IntVar(&queryPort, "qp", 1235, "port")
	flag.Parse()

	initConfig()
	InitConnection()

	if queryServer {
		go startQueryServer()
	}

	logger.Info("Start proxy server on 127.0.0.1:%d", port)
	cert, err := genCertificate()
	if err != nil {
		logger.Fatal(err)
	}

	server := &http.Server{
		Addr:      fmt.Sprintf("0.0.0.0:%d", port),
		TLSConfig: &tls.Config{Certificates: []tls.Certificate{cert}},
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			proxyHandler(w, r)
		}),
	}

	logger.Fatal(server.ListenAndServe())

	//err := http.ListenAndServe(fmt.Sprintf(":%d", port),
	//	http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	//		proxyHandler(w, r)
	//	}))
	//if err != nil {
	//	logger.Error(err)
	//}
	//os.Exit(0)
}
