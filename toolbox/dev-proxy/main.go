package main

import (
	"crypto/tls"
	"fmt"
	"github.com/kuangcp/gobase/pkg/ctool"
	"github.com/kuangcp/logger"
	"net/http"
)

var (
	port         int
	reloadConf   bool
	debug        bool
	queryPort    int
	buildVersion string
	proxyUri     string
	proxy        bool
)

var helpInfo = ctool.HelpInfo{
	Description:  "Http proxy for reroute and trace",
	BuildVersion: buildVersion,
	Version:      "1.0.3",
	Flags: []ctool.ParamVO{
		{Short: "-r", BoolVar: &reloadConf, Comment: "auto reload changed config"},
		{Short: "-d", BoolVar: &debug, Comment: "debug mode"},
		{Short: "-x", BoolVar: &proxy, Comment: "proxy mode"},
	},
	Options: []ctool.ParamVO{
		{Short: "-qp", IntVar: &queryPort, Int: 1235, Value: "port", Comment: "web port"},
		{Short: "-p", IntVar: &port, Int: 1234, Value: "port", Comment: "port"},
		{Short: "-pu", StringVar: &proxyUri, String: "http://localhost:7890", Value: "uri", Comment: "proxy uri"},
	},
}

func main() {
	helpInfo.Parse()

	initConfig()
	InitConnection()
	defer CloseConnection()

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
}
