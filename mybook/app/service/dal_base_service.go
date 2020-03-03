package service

import (
	"github.com/kuangcp/gobase/mybook/app/dal"
	"github.com/kuangcp/gobase/mybook/app/domain"
)

func AutoMigrateAll() {
	db := dal.GetDB()
	db.AutoMigrate(&domain.Account{})
	db.AutoMigrate(&domain.Category{})
	db.AutoMigrate(&domain.Record{})
	db.AutoMigrate(&domain.BookKeeping{})
}
