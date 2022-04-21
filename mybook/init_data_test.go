package main

import (
	"mybook/app/account"
	"mybook/app/book"
	"mybook/app/common"
	"mybook/app/common/constant"
	"testing"
)

func TestInitDBAndTable(t *testing.T) {
	common.AutoMigrateAll()
}

func TestInitBook(t *testing.T) {
	book.AddBookkeeping(&book.BookKeeping{Name: "主账本1", Comment: ""})
}

func TestInitAccount(t *testing.T) {
	//account.AddAccount(&account.Account{TypeId: constant.AccountCash, Name: "现金", InitAmount: 0})
	//account.AddAccount(&account.Account{TypeId: constant.AccountCredit, Name: "花呗", InitAmount: 0, MaxAmount: 2000, BillDay: 1, RepaymentDay: 10})
	account.AddAccount(&account.Account{TypeId: constant.AccountCredit, Name: "兴业信用卡", InitAmount: 0, MaxAmount: 2000, BillDay: 1, RepaymentDay: 10})
	//account.AddAccount(&account.Account{TypeId: constant.AccountOnline, Name: "支付宝", InitAmount: 0})
	//account.AddAccount(&account.Account{TypeId: constant.AccountOnline, Name: "微信", InitAmount: 0})
	//account.AddAccount(&account.Account{TypeId: constant.AccountDeposit, Name: "储蓄卡", InitAmount: 0})
	//account.AddAccount(&account.Account{TypeId: constant.AccountDeposit, Name: "应收款", InitAmount: 0})
	//account.AddAccount(&account.Account{TypeId: constant.AccountCredit, Name: "应付款", InitAmount: 0})
}

// 初始化分类数据
func TestInitCategory(t *testing.T) {
	common.InitCategory()
}
