package dal

import (
	"log"

	"github.com/jinzhu/gorm"
	"github.com/kuangcp/gobase/mybook/conf"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/viper"
)

var config *conf.ConnectionConfig
var db *gorm.DB

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

func OpenDB() *gorm.DB {
	return getConnectionWithConfig(GetDBConfig())
}

func GetDB() *gorm.DB {
	if db != nil {
		return db
	} else {
		db = OpenDB()
		return db
	}
}

func BatchSaveWithTransaction(records ...interface{}) error {
	db := GetDB()
	tx := db.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	for i := range records {
		if err := tx.Error; err != nil {
			return err
		}

		if err := tx.Create(records[i]).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit().Error
}

func Close(db *gorm.DB) {
	e := db.Close()
	if e != nil {
		log.Fatal(e)
	}
}

func getConnectionWithConfig(config *conf.ConnectionConfig) *gorm.DB {
	db, err := gorm.Open(config.DriverName, config.Path)
	if err != nil {
		log.Fatal(err)
	}

	return db
}
