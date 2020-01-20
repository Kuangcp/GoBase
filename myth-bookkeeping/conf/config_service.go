package conf

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
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
	if !loaded {
		log.Println("load config file")
		viper.SetConfigName("bookkeeping")
		viper.SetConfigType("yaml")
		viper.AddConfigPath("$HOME/.config")
		err := viper.ReadInConfig()
		if err != nil {
			log.Println(fmt.Errorf("Fatal error config file: %s \n", err))
		}
		loaded = true
	}
}
