package service

import (
	"mybook/app/common/dal"
	"mybook/app/domain"
)

func AutoMigrateAll() {
	db := dal.GetDB()
	db.AutoMigrate(&domain.Account{})
	db.AutoMigrate(&domain.Category{})
	db.AutoMigrate(&domain.Record{})
	db.AutoMigrate(&domain.BookKeeping{})
}
