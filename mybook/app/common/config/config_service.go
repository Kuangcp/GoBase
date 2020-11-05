package config

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/wonderivan/logger"
)

type (
	AppConfig struct {
		Path       string `json:"path"` // SQLite 文件绝对路径
		DriverName string `json:"driver"`
		Debug      bool   `json:"debug"` // 是否 Debug 模式
		Port       int    `json:"port"`
	}
)

var DefaultDBPath = "/tmp/bookkeeping.db"
var DefaultDriver = "sqlite3"
var DefaultPort = 9090
var DefaultUrlPath = "/api"

var config *AppConfig

// GetAppConfig 加载配置文件
func GetAppConfig() *AppConfig {
	if config != nil {
		return config
	}

	loadConfigFile()
	configLogger()
	config = buildAppConfig()

	show, _ := json.Marshal(config)
	logger.Info("Final config: %v", string(show))
	return config
}

func loadConfigFile() {
	logger.SetLogPathTrim("mybook/app/")
	viper.SetConfigName("mybook")
	viper.SetConfigType("yaml")
	// 短路式搜索配置文件
	viper.AddConfigPath("./data")
	viper.AddConfigPath("$HOME/.config")
	err := viper.ReadInConfig()
	if err != nil {
		logger.Warn("Use default config. %s", err)
	}
}

func configLogger() {
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

func buildAppConfig() *AppConfig {
	dbFile := viper.GetString("db.file")
	if dbFile == "" {
		dbFile = DefaultDBPath
	}
	driver := viper.GetString("driver")
	if driver == "" {
		driver = DefaultDriver
	}
	port := viper.GetInt("port")
	if port == 0 {
		port = DefaultPort
	}
	debug := viper.GetBool("debug")
	return &AppConfig{
		Path:       dbFile,
		DriverName: driver,
		Debug:      debug,
		Port:       port,
	}
}
