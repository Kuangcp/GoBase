package service

import (
	"github.com/kuangcp/gobase/mybook/dal"
	"github.com/kuangcp/gobase/mybook/domain"
	"github.com/kuangcp/gobase/mybook/util"
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

func FindAccountById(id uint) *domain.Account {
	db := dal.GetDB()
	var account domain.Account
	db.Where("id = ?", id).First(&account)
	return &account
}
