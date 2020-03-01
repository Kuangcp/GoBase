package conf

import (
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
	LoadConfig()

	path := viper.GetString("path")
	if path == "" {
		path = DefaultPath
	}
	driver := viper.GetString("driver")
	if driver == "" {
		driver = DefaultDriver
	}

	debug := viper.GetBool("debug")
	config = &AppConfig{Path: path, DriverName: driver, Debug: debug}
	return config
}

func LoadConfig() {
	logger.Info("load config file ~/.config/mybook.yml")
	viper.SetConfigName("mybook")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$HOME/.config")
	err := viper.ReadInConfig()
	if err != nil {
		logger.Error("Fatal error config file: %s \n", err)
	}
	configLogger()
}

func configLogger() {
	logger.SetLogPathTrim("mybook/")

	debug := viper.GetBool("debug")
	jsonPath := ""
	if debug {
		jsonPath = "./resources/log-dev.json"
	} else {
		jsonPath = "./resources/log.json"
	}
	e := logger.SetLogger(jsonPath)
	if e != nil {
		logger.Error(e)
	}
}
