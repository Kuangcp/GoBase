package conf

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

	configDir := viper.GetString("configDir")
	dbFile := viper.GetString("db.file")
	if dbFile == "" {
		dbFile = DefaultPath
	} else {
		dbFile = configDir + "/" + dbFile
	}
	driver := viper.GetString("driver")
	if driver == "" {
		driver = DefaultDriver
	}

	debug := viper.GetBool("debug")
	config = &AppConfig{Path: dbFile, DriverName: driver, Debug: debug}

	logger.Info("load config file ~/.config/app-conf/mybook/mybook.yml")
	return config
}

func loadConfig() {
	viper.SetConfigName("mybook")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$HOME/.config/app-conf/mybook")
	err := viper.ReadInConfig()
	if err != nil {
		logger.Error("Fatal error config file: %s \n", err)
	}
}

func configLogger() {
	logger.SetLogPathTrim("mybook/")

	debug := viper.GetBool("debug")
	configDir := viper.GetString("configDir")
	jsonPath := ""
	if debug {
		jsonPath = configDir + "/resources/log-dev.json"
	} else {
		jsonPath = configDir + "/resources/log.json"
		gin.SetMode(gin.ReleaseMode)
	}
	e := logger.SetLogger(jsonPath)
	if e != nil {
		logger.Error(e)
	}
}
