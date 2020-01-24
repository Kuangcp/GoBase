package service

import (
	"github.com/kuangcp/gobase/myth-bookkeeping/dal"
	"github.com/kuangcp/gobase/myth-bookkeeping/domain"
	"github.com/kuangcp/gobase/myth-bookkeeping/util"
	"log"
)

func ListAccounts() []domain.Account {
	db := dal.GetDB()

	var accounts []domain.Account
	e := db.Find(&accounts).Error
	util.AssertNoError(e)

	log.Println(len(accounts))
	return accounts
}

func AddAccount(account *domain.Account) {
	db := dal.GetDB()

	create := db.Create(account)
	log.Println(create)
}

func UpdateAccount(account *domain.Account) {
	db := dal.GetDB()
	db.Update(account)
}
