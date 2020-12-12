package account

import (
	"fmt"

	"github.com/wonderivan/logger"
	"mybook/app/common/constant"
	"mybook/app/common/dal"
	"mybook/app/common/util"
)

func ListAccounts() []*Account {
	db := dal.GetDB()

	var accounts []*Account
	e := db.Find(&accounts).Error
	util.RecordError(e)

	return accounts
}

func ListAccountMap() map[uint]*Account {
	accounts := ListAccounts()
	result := make(map[uint]*Account)
	for i := range accounts {
		account := accounts[i]
		result[account.ID] = account
	}
	return result
}

func AddAccount(account *Account) {
	db := dal.GetDB()

	create := db.Create(account)
	logger.Info(create)
}

func UpdateAccount(account *Account) {
	db := dal.GetDB()
	db.Update(account)
}

func FindAccountById(id uint) *Account {
	db := dal.GetDB()
	var account Account
	db.Where("id = ?", id).First(&account)
	return &account
}

func PrintAccount() {
	db := dal.GetDB()
	var lists []Account
	db.Where("1=1").Order("id", false).Find(&lists)
	for i := range lists {
		account := lists[i]
		chFormat := util.BuildCHCharFormat(12, account.Name)
		fmt.Printf("  %d  %s "+chFormat+" %-7d %s \n", account.ID, account.Name, "",
			account.CurrentAmount, constant.GetAccountTypeByIndex(account.TypeId).GetName())
	}
}
