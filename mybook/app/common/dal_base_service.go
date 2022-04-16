package common

import (
	"mybook/app/account"
	"mybook/app/book"
	"mybook/app/category"
	"mybook/app/common/dal"
	"mybook/app/loan"
	"mybook/app/record"
	"mybook/app/user"
)

// AutoMigrateAll 建表
func AutoMigrateAll() {
	db := dal.GetDB()
	db.AutoMigrate(&account.Account{})
	db.AutoMigrate(&category.Category{})
	db.AutoMigrate(&record.RecordEntity{})
	db.AutoMigrate(&book.BookKeeping{})
	db.AutoMigrate(&loan.Entity{})
	db.AutoMigrate(&user.User{})
}
