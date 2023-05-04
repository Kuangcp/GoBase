package core

import (
	"crypto/tls"
	"fmt"
	"github.com/kuangcp/gobase/pkg/ctool"
	"github.com/kuangcp/logger"
	"net/http"
)

var (
	Port         int
	ReloadConf   bool
	Debug        bool
	QueryPort    int
	BuildVersion string
	HttpProxyUri string
	HttpProxy    bool
)

var HelpInfo = ctool.HelpInfo{
	Description:  "Http proxy for reroute and trace",
	BuildVersion: BuildVersion,
	Version:      "1.0.3",
	Flags: []ctool.ParamVO{
		{Short: "-r", BoolVar: &ReloadConf, Comment: "auto reload changed config"},
		{Short: "-d", BoolVar: &Debug, Comment: "debug mode"},
		{Short: "-x", BoolVar: &HttpProxy, Comment: "proxy mode"},
	},
	Options: []ctool.ParamVO{
		{Short: "-qp", IntVar: &QueryPort, Int: 1235, Value: "port", Comment: "web port"},
		{Short: "-p", IntVar: &Port, Int: 1234, Value: "port", Comment: "port"},
		{Short: "-pu", StringVar: &HttpProxyUri, String: "http://localhost:7890", Value: "uri", Comment: "proxy uri"},
	},
}

func StartMainServer() {
	InitConfig()
	InitConnection()
	defer CloseConnection()

	logger.Info("list key: ", RequestList)
	logger.Info("Start proxy server on 127.0.0.1:%d", Port)
	cert, err := GenCertificate()
	if err != nil {
		logger.Fatal(err)
	}

	server := &http.Server{
		Addr:      fmt.Sprintf("0.0.0.0:%d", Port),
		TLSConfig: &tls.Config{Certificates: []tls.Certificate{cert}},
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ProxyHandler(w, r)
		}),
	}

	//go func() {
	logger.Fatal(server.ListenAndServe())
	//}()

	//systray.Run(gui.OnReady, gui.OnExit)
}
