package config

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/wonderivan/logger"
)

type (
	AppConfig struct {
		// SQLite 文件绝对路径
		Path       string
		DriverName string
		// 是否 Debug 模式
		Debug bool
	}
)

var DefaultPath = "/tmp/bookkeeping.db"
var DefaultDriver = "sqlite3"
var config *AppConfig

func GetAppConfig() *AppConfig {
	if config != nil {
		return config
	}
	loadConfig()
	configLogger()

	dbFile := viper.GetString("db.file")
	if dbFile == "" {
		dbFile = DefaultPath
	}
	driver := viper.GetString("driver")
	if driver == "" {
		driver = DefaultDriver
	}

	debug := viper.GetBool("debug")
	config = &AppConfig{Path: dbFile, DriverName: driver, Debug: debug}
	return config
}

func loadConfig() {
	viper.SetConfigName("mybook")
	viper.SetConfigType("yaml")
	// 短路式搜索配置文件
	viper.AddConfigPath("./data")
	viper.AddConfigPath("$HOME/.config")
	err := viper.ReadInConfig()
	if err != nil {
		logger.Error("Fatal error config file: %s \n", err)
	}
}

func configLogger() {
	logger.SetLogPathTrim("mybook/")

	debug := viper.GetBool("debug")
	notDev := viper.GetBool("notDev")
	jsonPath := ""
	if debug {
		jsonPath = "./conf/log-dev.json"
	} else {
		jsonPath = "./conf/log.json"
		gin.SetMode(gin.ReleaseMode)
	}

	var e error
	if !notDev {
		e = logger.SetLogger()
	} else {
		e = logger.SetLogger(jsonPath)
	}
	if e != nil {
		logger.Error(e)
	}
}
