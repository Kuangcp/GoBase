package common

import (
	"mybook/app/account"
	"mybook/app/book"
	"mybook/app/category"
	"mybook/app/common/constant"
	"mybook/app/common/dal"
	"mybook/app/loan"
	"mybook/app/record"
	"mybook/app/user"

	"github.com/kuangcp/logger"
)

// AutoMigrateAll 建表
func AutoMigrateAll() {
	db := dal.GetDB()
	db.AutoMigrate(&account.Account{})
	db.AutoMigrate(&category.Category{})
	db.AutoMigrate(&record.RecordEntity{})
	db.AutoMigrate(&book.BookKeeping{})
	db.AutoMigrate(&loan.Entity{})
	db.AutoMigrate(&user.User{})

	initData()
}

func initData() {
	db := dal.GetDB()
	var b book.BookKeeping
	first := db.Limit(1).Find(&b)
	if first.Error != nil && first.Error.Error() != "record not found" {
		logger.Error(first)
		return
	}

	if b.Name != "" {
		logger.Info("already init")
		return
	}

	book.AddBookkeeping(&book.BookKeeping{Name: "主账本", Comment: ""})

	InitCategory()
	InitAccount()
}

func InitAccount() {
	account.AddAccount(&account.Account{TypeId: constant.AccountCash, Name: "现金", InitAmount: 0})
	account.AddAccount(&account.Account{TypeId: constant.AccountCredit, Name: "花呗", InitAmount: 0, MaxAmount: 2000, BillDay: 1, RepaymentDay: 10})
	account.AddAccount(&account.Account{TypeId: constant.AccountOnline, Name: "支付宝", InitAmount: 0})
	account.AddAccount(&account.Account{TypeId: constant.AccountOnline, Name: "微信", InitAmount: 0})
	account.AddAccount(&account.Account{TypeId: constant.AccountDeposit, Name: "储蓄卡", InitAmount: 0})

	ar := account.Account{TypeId: constant.AccountAR, Name: "应收款", InitAmount: 0}
	ar.ID = constant.AccountARId
	account.AddAccount(&ar)
	ap := account.Account{TypeId: constant.AccountAP, Name: "应付款", InitAmount: 0}
	ap.ID = constant.AccountAPId
	account.AddAccount(&ap)
}

// InitCategory 初始化分类数据
func InitCategory() {
	expenseIndex := 10
	var types = []string{"食", "住", "行", "娱", "资金"}
	for e := range types {
		expenseIndex++
		tmpCategory := &category.Category{Name: types[e], Leaf: false, TypeId: constant.CategoryExpense}
		tmpCategory.ID = uint(expenseIndex)
		category.AddCategory(tmpCategory)
	}

	expenseIndex = 100
	types = []string{"早餐", "午餐", "晚餐", "餐厅", "零食", "日用品", "室外娱乐", "服饰", "云服务", "水果", "买菜",
		"发红包", "房租", "书籍", "话费网费", "数码", "礼物", "地铁", "公交", "打车", "火车", "酒店", "医疗", "知识付费", "坏账",
		"景点门票", "网络会员", "线下会员", "水电煤", "美容美发", "快递", "投资亏损", "电影", "保险", "打赏", "还贷"}
	for e := range types {
		expenseIndex++
		tmpCategory := &category.Category{Name: types[e], Leaf: true, TypeId: constant.CategoryExpense}
		tmpCategory.ID = uint(expenseIndex)
		category.AddCategory(tmpCategory)
	}

	category.SetParentId("早餐", 11)
	category.SetParentId("午餐", 11)
	category.SetParentId("晚餐", 11)
	category.SetParentId("餐厅", 11)
	category.SetParentId("买菜", 11)
	category.SetParentId("水果", 11)
	category.SetParentId("零食", 11)

	category.SetParentId("水电煤", 12)
	category.SetParentId("日用品", 12)
	category.SetParentId("房租", 12)
	category.SetParentId("酒店", 12)

	category.SetParentId("地铁", 13)
	category.SetParentId("公交", 13)
	category.SetParentId("打车", 13)
	category.SetParentId("火车", 13)

	category.SetParentId("室外娱乐", 14)
	category.SetParentId("云服务", 14)
	category.SetParentId("书籍", 14)
	category.SetParentId("数码", 14)
	category.SetParentId("知识付费", 14)
	category.SetParentId("景点门票", 14)
	category.SetParentId("电影", 14)
	category.SetParentId("网络会员", 14)
	category.SetParentId("线下会员", 14)
	category.SetParentId("打赏", 14)

	incomeIndex := 200
	types = []string{"应收款", "收红包", "投资收益", "返现", "工资", "奖金", "兼职", "生活费",
		"报销流入", "其他收入", "退款", "问题"}
	for e := range types {
		incomeIndex++
		tmpCategory := &category.Category{Name: types[e], Leaf: true, TypeId: constant.CategoryIncome}
		tmpCategory.ID = uint(incomeIndex)
		category.AddCategory(tmpCategory)
	}

	//
	types = []string{"转账", "加仓", "平仓"}
	ids := []uint{constant.CategoryTransferId, constant.CategoryTransferOpenId, constant.CategoryTransferCloseId}
	for e := range types {
		tmpCategory := &category.Category{Name: types[e], Leaf: true, TypeId: constant.CategoryTransfer}
		tmpCategory.ID = ids[e]
		category.AddCategory(tmpCategory)
	}
}
