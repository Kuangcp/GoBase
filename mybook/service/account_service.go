package service

import (
	"fmt"
	"github.com/kuangcp/gobase/mybook/constant"
	"github.com/kuangcp/gobase/mybook/dal"
	"github.com/kuangcp/gobase/mybook/domain"
	"github.com/kuangcp/gobase/mybook/util"
	"github.com/wonderivan/logger"
)

func ListAccounts() []domain.Account {
	db := dal.GetDB()

	var accounts []domain.Account
	e := db.Find(&accounts).Error
	util.AssertNoError(e)

	return accounts
}

func AddAccount(account *domain.Account) {
	db := dal.GetDB()

	create := db.Create(account)
	logger.Info(create)
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

func PrintAccount(_ []string) {
	db := dal.GetDB()
	var lists []domain.Account
	db.Where("1=1").Order("id", false).Find(&lists)
	for i := range lists {
		account := lists[i]
		chFormat := util.BuildCHCharFormat(12, account.Name)
		fmt.Printf("  %d  %s "+chFormat+" %s\n", account.ID, account.Name, "",
			constant.GetAccountTypeByIndex(account.TypeId).Name)
	}
}
