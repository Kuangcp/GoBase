package service

import (
	"fmt"

	"mybook/app/common/constant"
	"mybook/app/common/dal"
	"mybook/app/common/util"
	"mybook/app/domain"
	"github.com/wonderivan/logger"
)

func ListAccounts() []*domain.Account {
	db := dal.GetDB()

	var accounts []*domain.Account
	e := db.Find(&accounts).Error
	util.RecordError(e)

	return accounts
}

func ListAccountMap() map[uint]*domain.Account {
	accounts := ListAccounts()
	result := make(map[uint]*domain.Account)
	for i := range accounts {
		account := accounts[i]
		result[account.ID] = account
	}
	return result
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

func PrintAccount() {
	db := dal.GetDB()
	var lists []domain.Account
	db.Where("1=1").Order("id", false).Find(&lists)
	for i := range lists {
		account := lists[i]
		chFormat := util.BuildCHCharFormat(12, account.Name)
		fmt.Printf("  %d  %s "+chFormat+" %-7d %s \n", account.ID, account.Name, "",
			account.CurrentAmount, constant.GetAccountTypeByIndex(account.TypeId).Name)
	}
}
