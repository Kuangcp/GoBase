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
	HttpProxy    bool
	JsonPath     string
	PacPath      string
)

var HelpInfo = ctool.HelpInfo{
	Description:  "Http proxy for reroute and trace",
	BuildVersion: BuildVersion,
	Version:      "1.0.4",
	Flags: []ctool.ParamVO{
		{Short: "-r", BoolVar: &ReloadConf, Comment: "auto reload changed config"},
		{Short: "-d", BoolVar: &Debug, Comment: "debug mode"},
		{Short: "-x", BoolVar: &HttpProxy, Comment: "only track http proxy, default capture https packet"},
	},
	Options: []ctool.ParamVO{
		{Short: "-w", IntVar: &QueryPort, Int: 1235, Value: "port", Comment: "web port"},
		{Short: "-p", IntVar: &Port, Int: 1234, Value: "port", Comment: "proxy port"},
		{Short: "-j", StringVar: &JsonPath, String: "", Value: "path", Comment: "json config file abs path"},
		{Short: "-a", StringVar: &PacPath, String: "", Value: "path", Comment: "pac file abs path"},
	},
}

// StartMainServer HTTP代理和修改 HTTPS转发
func StartMainServer() {
	logger.Info("list key: ", RequestList)
	logger.Info("Start HTTP proxy server on 127.0.0.1:%d", Port)
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

	logger.Fatal(server.ListenAndServe())
}
