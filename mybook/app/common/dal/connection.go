package dal

import (
	"log"
	"mybook/app/common/config"

	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
)

var db *gorm.DB

func OpenDB() *gorm.DB {
	return getConnectionWithConfig()
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

func getConnectionWithConfig() *gorm.DB {
	db, err := gorm.Open(config.AppConf.DriverName, config.AppConf.DBFilePath)
	if err != nil {
		log.Fatal(err)
	}
	db.LogMode(!config.AppConf.Release)

	return db
}
