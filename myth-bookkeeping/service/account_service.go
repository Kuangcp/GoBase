package service

import (
	"log"

	"github.com/kuangcp/gobase/myth-bookkeeping/data_source"
	"github.com/kuangcp/gobase/myth-bookkeeping/domain"
)

func QueryAllAccounts() []domain.Account {
	db := data_source.GetDB()
	defer data_source.Close(db)
	return nil
}

func Insert(account *domain.Account) {
	db := data_source.GetDB()
	defer data_source.Close(db)

	migrate := db.AutoMigrate(&domain.Account{})
	log.Println(migrate)

	create := db.Create(account)
	log.Println(create)
}
