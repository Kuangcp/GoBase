package conf

import (
	"fmt"
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

func LoadConfig() {
	logger.SetLogPathTrim("myth-bookkeeping/")
	if !loaded {

		logger.Info("load config file ~/.config/bookkeeping.yml")
		viper.SetConfigName("bookkeeping")
		viper.SetConfigType("yaml")
		viper.AddConfigPath("$HOME/.config")
		err := viper.ReadInConfig()
		if err != nil {
			logger.Error(fmt.Errorf("Fatal error config file: %s \n", err))
		}
		loaded = true
	}
}
