package main

import (
	"crypto/tls"
	"fmt"
	"github.com/google/uuid"
	"github.com/kuangcp/gobase/pkg/ctool"
	"github.com/kuangcp/logger"
	"net/http"
	"os"
)

var (
	port         int
	reloadConf   bool
	debug        bool
	queryPort    int
	buildVersion string
)

func init() {
	hostname, err := os.Hostname()
	if err != nil {
		logger.Error("Random hostname. err:", err)
		hostname = uuid.NewString()
	}

	RequestList += hostname
}

var helpInfo = ctool.HelpInfo{
	Description:  "Http proxy for reroute and trace",
	BuildVersion: buildVersion,
	Version:      "1.0.3",
	Flags: []ctool.ParamVO{
		{Short: "-r", BoolVar: &reloadConf, Comment: "auto reload changed config"},
		{Short: "-d", BoolVar: &debug, Comment: "debug mode"},
	},
	Options: []ctool.ParamVO{
		{Short: "-qp", IntVar: &queryPort, Int: 1235, Value: "port", Comment: "web port"},
		{Short: "-p", IntVar: &port, Int: 1234, Value: "port", Comment: "port"},
	},
}

func main() {
	helpInfo.Parse()

	initConfig()
	InitConnection()

	go startQueryServer()

	logger.Info("list key: ", RequestList)
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
