package dal

import (
	"log"

	"github.com/jinzhu/gorm"
	"github.com/kuangcp/gobase/myth-bookkeeping/conf"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/viper"
)

var config *conf.ConnectionConfig

func GetDBConfig() *conf.ConnectionConfig {
	conf.LoadConfig()
	if config != nil {
		return config
	}

	path := viper.GetString("path")
	if path == "" {
		path = conf.DefaultPath
	}
	driver := viper.GetString("driver")
	if driver == "" {
		driver = conf.DefaultDriver
	}

	return &conf.ConnectionConfig{Path: path, DriverName: driver}
}

func GetDB() *gorm.DB {
	return getConnectionWithConfig(GetDBConfig())
}

func getConnectionWithConfig(config *conf.ConnectionConfig) *gorm.DB {
	db, err := gorm.Open(config.DriverName, config.Path)
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func Close(db *gorm.DB) {
	e := db.Close()
	if e != nil {
		log.Fatal(e)
	}
}
