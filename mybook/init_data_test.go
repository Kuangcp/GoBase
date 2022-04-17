package main

import (
	"mybook/app/account"
	"mybook/app/book"
	"mybook/app/category"
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
	expenseIndex := 10
	var types = []string{"日常餐", "文娱", "日常开支", "交通"}
	for e := range types {
		expenseIndex++
		tmpCategory := &category.Category{Name: types[e], Leaf: false, TypeId: constant.CategoryExpense}
		tmpCategory.ID = uint(expenseIndex)
		category.AddCategory(tmpCategory)
	}

	expenseIndex = 100
	types = []string{"早餐", "午餐", "晚餐", "餐厅", "零食", "日用品", "室外娱乐", "服饰", "云服务", "水果", "买菜",
		"发红包", "房租", "书籍", "话费网费", "火车", "数码", "礼物", "地铁", "酒店", "医疗", "公交", "打车", "知识付费", "坏账",
		"景点门票", "会员", "水电煤", "美容美发", "快递", "投资亏损", "电影", "保险", "打赏", "还贷", "己方借出"}
	for e := range types {
		expenseIndex++
		tmpCategory := &category.Category{Name: types[e], Leaf: true, TypeId: constant.CategoryExpense}
		tmpCategory.ID = uint(expenseIndex)
		category.AddCategory(tmpCategory)
	}
	category.SetParentId("早餐", 11)
	category.SetParentId("午餐", 11)
	category.SetParentId("晚餐", 11)

	incomeIndex := 200
	types = []string{"应收款", "收红包", "投资收益", "返现", "工资", "个缴社保", "司缴社保", "平台", "奖金", "兼职", "生活费",
		"报销流入", "其他收入", "退款", "问题", "借贷", "对方归还"}
	for e := range types {
		incomeIndex++
		tmpCategory := &category.Category{Name: types[e], Leaf: true, TypeId: constant.CategoryIncome}
		tmpCategory.ID = uint(incomeIndex)
		category.AddCategory(tmpCategory)
	}

	//types = []string{"转账", "加仓", "平仓"}
	types = []string{"转账", "投资", "变现"}
	ids := []uint{constant.CategoryTransferId, constant.CategoryTransferOpenId, constant.CategoryTransferCloseId}
	for e := range types {
		tmpCategory := &category.Category{Name: types[e], Leaf: true, TypeId: constant.CategoryTransfer}
		tmpCategory.ID = ids[e]
		category.AddCategory(tmpCategory)
	}
}

func TestSetParentId(t *testing.T) {
	category.SetParentId("早餐", 1)
}
