package common

import (
	"mybook/app/account"
	"mybook/app/common/dal"
	"mybook/app/domain"
	"mybook/app/record"
)

func AutoMigrateAll() {
	db := dal.GetDB()
	db.AutoMigrate(&account.Account{})
	db.AutoMigrate(&domain.Category{})
	db.AutoMigrate(&record.RecordEntity{})
	db.AutoMigrate(&domain.BookKeeping{})
}
