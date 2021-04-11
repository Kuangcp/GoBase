package config

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/kuangcp/gobase/pkg/cuibase"
	"github.com/kuangcp/logger"
)

type (
	AppConfig struct {
		DBFilePath  string `json:"path"` // SQLite 文件绝对路径
		DriverName  string `json:"driver"`
		Debug       bool   `json:"debug"` // 是否 Debug 模式
		DebugStatic bool
		Port        int `json:"port"`
	}
)

//var DefaultDBPath = "/tmp/bookkeeping.db"
var DefaultDBPath = "./data/main.db"
var DefaultDriver = "sqlite3"
var DefaultPort = 9090
var DefaultUrlPath = "/api"

var AppConf = &AppConfig{
	DriverName: DefaultDriver,
	DBFilePath: DefaultDBPath,
}

// InitAppConfig 初始化配置文件
func InitAppConfig() *AppConfig {
	configLogger()

	if AppConf.Debug {

	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	show, _ := json.Marshal(AppConf)
	logger.Info("Final config: %v", string(show))
	return AppConf
}

func configLogger() {
	logger.SetLogPathTrim("mybook/app/")

	err := logger.SetLoggerConfig(&logger.LogConfig{
		TimeFormat: cuibase.YYYY_MM_DD_HH_MM_SS_MS,
		Console: &logger.ConsoleLogger{
			Level:    logger.DebugDesc,
			Colorful: true,
		},
		File: &logger.FileLogger{
			Filename:   "app.log",
			Level:      logger.DebugDesc,
			Colorful:   true,
			Append:     true,
			PermitMask: "0660",
		},
	})
	if err != nil {
		logger.Fatal(err)
	}
}
