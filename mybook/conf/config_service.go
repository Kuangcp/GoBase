package conf

import (
	"github.com/spf13/viper"
	"github.com/wonderivan/logger"
)

type (
	ConnectionConfig struct {
		Path       string
		DriverName string
	}
)

var DefaultPath = "/tmp/bookkeeping.db"
var DefaultDriver = "sqlite3"
var loaded = false

func configLogger() {
	logger.SetLogPathTrim("mybook/")
}

func LoadConfig() {
	configLogger()

	if !loaded {
		logger.Info("load config file ~/.config/mybook.yml")
		viper.SetConfigName("mybook")
		viper.SetConfigType("yaml")
		viper.AddConfigPath("$HOME/.config")
		err := viper.ReadInConfig()
		if err != nil {
			logger.Error("Fatal error config file: %s \n", err)
		}
		loaded = true
	}
}
