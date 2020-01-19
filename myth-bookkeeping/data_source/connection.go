package data_source

import (
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
)

type (
	ConnectionConfig struct {
		Path       string
		DriverName string
	}
)

var config = &ConnectionConfig{Path: "/tmp/test.data_source", DriverName: "sqlite3"}

func GetDBConfig() *ConnectionConfig {
	if config != nil {
		return config
	}
	config = &ConnectionConfig{Path: "/tmp/bookkeeping.data_source", DriverName: "sqlite3"}
	return config
}

func GetDB() *gorm.DB {
	return getConnectionWithConfig(GetDBConfig())
}

func getConnectionWithConfig(config *ConnectionConfig) *gorm.DB {
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
