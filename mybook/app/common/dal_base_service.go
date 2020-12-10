package common

import (
	"mybook/app/common/dal"
	"mybook/app/domain"
	"mybook/app/record"
)

func AutoMigrateAll() {
	db := dal.GetDB()
	db.AutoMigrate(&domain.Account{})
	db.AutoMigrate(&domain.Category{})
	db.AutoMigrate(&record.RecordEntity{})
	db.AutoMigrate(&domain.BookKeeping{})
}
