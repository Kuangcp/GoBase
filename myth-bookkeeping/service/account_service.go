package service

import (
	"log"
	"time"

	"github.com/kuangcp/gobase/myth-bookkeeping/dal"
	"github.com/kuangcp/gobase/myth-bookkeeping/domain"
)

func QueryAllAccounts() []domain.Account {
	db := dal.GetDB()
	defer dal.Close(db)

	migrate := db.AutoMigrate(&domain.Account{})
	log.Println(migrate)

	db.Create(&domain.Account{Name: "test", InitAmount: 0, Type: 1, CreatedAt: time.Now().Unix(), UpdatedAt: time.Now().Unix(), DeletedAt: 0})
	var accounts []domain.Account

	var count int16
	e := db.Find(&accounts).Count(&count).Error
	if e != nil {
		log.Println(e)
	}
	log.Println(count, accounts, db.HasTable("account"))
	return accounts
}

func Insert(account *domain.Account) {
	db := dal.GetDB()
	defer dal.Close(db)

	migrate := db.AutoMigrate(&domain.Account{})
	log.Println(migrate)

	create := db.Create(account)
	log.Println(create)
}
