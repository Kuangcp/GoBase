package config

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/kuangcp/gobase/pkg/cuibase"
	"github.com/kuangcp/logger"
)

type (
	AppConfig struct {
		DBFilePath  string // SQLite 文件绝对路径
		DriverName  string
		Release     bool // 是否 Release 模式
		DebugStatic bool
		Port        int
	}
)

//var DefaultDBPath = "/tmp/bookkeeping.db"
var DefaultDBPath = "./data/main.db"
var DefaultDriver = "sqlite3"
var DefaultPort = 9090
var DefaultUrlPath = "/api"

var AppConf *AppConfig = &AppConfig{
	Release: false,
	Port: DefaultPort,
	DriverName: DefaultDriver,
	DBFilePath: DefaultDBPath,
}

// InitAppConfig 初始化配置文件
func InitAppConfig() *AppConfig {
	logger.SetLogPathTrim("mybook/app/")

	fileLogger := &logger.FileLogger{
		Filename:   "app.log",
		Level:      logger.DebugDesc,
		Colorful:   true,
		Append:     true,
		PermitMask: "0660",
	}
	consoleLogger := &logger.ConsoleLogger{
		Level:    logger.DebugDesc,
		Colorful: true,
	}
	if AppConf.Release {
		fileLogger.Level = logger.InformationalDesc
		consoleLogger.Level = logger.InformationalDesc
		gin.SetMode(gin.ReleaseMode)
	} else {
		fileLogger.Level = logger.DebugDesc
		consoleLogger.Level = logger.DebugDesc
	}

	err := logger.SetLoggerConfig(&logger.LogConfig{
		TimeFormat: cuibase.YYYY_MM_DD_HH_MM_SS_MS,
		Console:    consoleLogger,
		File:       fileLogger,
	})
	if err != nil {
		logger.Fatal(err)
	}

	show, _ := json.Marshal(AppConf)
	logger.Info("Final config: %v", string(show))
	return AppConf
}
