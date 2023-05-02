package main

import (
	"crypto/tls"
	"fmt"
	"github.com/kuangcp/gobase/pkg/ctool"
	"github.com/kuangcp/gobase/toolbox/dev-proxy/core"
	"github.com/kuangcp/logger"
	"net/http"
)

var helpInfo = ctool.HelpInfo{
	Description:  "Http proxy for reroute and trace",
	BuildVersion: core.BuildVersion,
	Version:      "1.0.3",
	Flags: []ctool.ParamVO{
		{Short: "-r", BoolVar: &core.ReloadConf, Comment: "auto reload changed config"},
		{Short: "-d", BoolVar: &core.Debug, Comment: "debug mode"},
		{Short: "-x", BoolVar: &core.HttpProxy, Comment: "proxy mode"},
	},
	Options: []ctool.ParamVO{
		{Short: "-qp", IntVar: &core.QueryPort, Int: 1235, Value: "port", Comment: "web port"},
		{Short: "-p", IntVar: &core.Port, Int: 1234, Value: "port", Comment: "port"},
		{Short: "-pu", StringVar: &core.HttpProxyUri, String: "http://localhost:7890", Value: "uri", Comment: "proxy uri"},
	},
}

func main() {
	helpInfo.Parse()

	core.InitConfig()
	core.InitConnection()
	defer core.CloseConnection()

	go core.StartQueryServer()

	logger.Info("list key: ", core.RequestList)
	logger.Info("Start proxy server on 127.0.0.1:%d", core.Port)
	cert, err := core.GenCertificate()
	if err != nil {
		logger.Fatal(err)
	}

	server := &http.Server{
		Addr:      fmt.Sprintf("0.0.0.0:%d", core.Port),
		TLSConfig: &tls.Config{Certificates: []tls.Certificate{cert}},
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			core.ProxyHandler(w, r)
		}),
	}

	//go func() {
	logger.Fatal(server.ListenAndServe())
	//}()

	//systray.Run(gui.OnReady, gui.OnExit)
}
