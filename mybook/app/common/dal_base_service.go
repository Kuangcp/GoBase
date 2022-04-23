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

	initEnum()
	initData()
}

var (
	expenseRootType = make(map[string]uint)
	expenseType     = make(map[string]uint)
	incomeType      = make(map[string]uint)
	transferType    = make(map[string]uint)
)

func initEnum() {
	index := 10
	var types = []string{"食", "日常", "娱", "行", "资金", "住"}
	for _, s := range types {
		index++
		expenseRootType[s] = uint(index)
	}

	index = 100
	types = []string{"早餐", "午餐", "晚餐", "餐厅", "零食", "日用品", "室外娱乐", "服饰", "云服务", "水果", "买菜",
		"发红包", "房租", "书籍", "话费网费", "数码", "礼物", "地铁", "公交", "打车", "火车", "酒店", "医疗", "知识付费", "坏账",
		"景点门票", "网络会员", "线下会员", "水电煤", "美容美发", "快递", "投资亏损", "电影", "保险", "打赏", "政务"}
	for _, s := range types {
		index++
		expenseType[s] = uint(index)
	}

	index = 200
	types = []string{"收红包", "投资收益", "工资", "奖金", "兼职", "其他收入", "坏账"}
	for _, s := range types {
		index++
		incomeType[s] = uint(index)
	}

	transferType["转账"] = constant.CategoryTransferId
	transferType["加仓"] = constant.CategoryTransferOpenId
	transferType["平仓"] = constant.CategoryTransferCloseId
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
		refreshCategoryBind()
		return
	}

	book.AddBookkeeping(&book.BookKeeping{Name: "主账本", Comment: ""})

	InitCategory()
	refreshCategoryBind()
	InitAccount()
}

func InitAccount() {
	crash := account.Account{TypeId: constant.AccountCash, Name: "现金", InitAmount: 0}
	crash.ID = uint(constant.AccountCash)*100 + 1
	account.AddAccount(&crash)

	huabei := account.Account{TypeId: constant.AccountCredit, Name: "花呗", InitAmount: 0, MaxAmount: 2000, BillDay: 1, RepaymentDay: 10}
	huabei.ID = uint(constant.AccountCredit)*100 + 1
	account.AddAccount(&huabei)

	zhifubao := account.Account{TypeId: constant.AccountOnline, Name: "支付宝", InitAmount: 0}
	zhifubao.ID = uint(constant.AccountOnline)*100 + 1
	account.AddAccount(&zhifubao)
	wechat := account.Account{TypeId: constant.AccountOnline, Name: "微信", InitAmount: 0}
	wechat.ID = uint(constant.AccountOnline)*100 + 2
	account.AddAccount(&wechat)

	deposit := account.Account{TypeId: constant.AccountDeposit, Name: "储蓄卡", InitAmount: 0}
	deposit.ID = uint(constant.AccountDeposit)*100 + 1
	account.AddAccount(&deposit)

	ar := account.Account{TypeId: constant.AccountAR, Name: "应收款", InitAmount: 0}
	ar.ID = constant.AccountARId
	account.AddAccount(&ar)
	ap := account.Account{TypeId: constant.AccountAP, Name: "应付款", InitAmount: 0}
	ap.ID = constant.AccountAPId
	account.AddAccount(&ap)
}

func refreshCategoryBind() {
	category.SetParentId("早餐", expenseRootType["食"])
	category.SetParentId("午餐", expenseRootType["食"])
	category.SetParentId("晚餐", expenseRootType["食"])
	category.SetParentId("餐厅", expenseRootType["食"])
	category.SetParentId("买菜", expenseRootType["食"])
	category.SetParentId("水果", expenseRootType["食"])
	category.SetParentId("零食", expenseRootType["食"])

	category.SetParentId("日用品", expenseRootType["日常"])
	category.SetParentId("宠物", expenseRootType["日常"])
	category.SetParentId("美容美发", expenseRootType["日常"])
	category.SetParentId("快递", expenseRootType["日常"])
	category.SetParentId("政务", expenseRootType["日常"])
	category.SetParentId("发红包", expenseRootType["日常"])
	category.SetParentId("水电煤", expenseRootType["日常"])
	category.SetParentId("话费网费", expenseRootType["日常"])

	category.SetParentId("室外娱乐", expenseRootType["娱"])
	category.SetParentId("云服务", expenseRootType["娱"])
	category.SetParentId("书籍", expenseRootType["娱"])
	category.SetParentId("数码", expenseRootType["娱"])
	category.SetParentId("知识付费", expenseRootType["娱"])
	category.SetParentId("考试报名", expenseRootType["娱"])
	category.SetParentId("景点门票", expenseRootType["娱"])
	category.SetParentId("电影", expenseRootType["娱"])
	category.SetParentId("网络会员", expenseRootType["娱"])
	category.SetParentId("线下会员", expenseRootType["娱"])
	category.SetParentId("打赏", expenseRootType["娱"])
	category.SetParentId("服饰", expenseRootType["娱"])

	category.SetParentId("地铁", expenseRootType["行"])
	category.SetParentId("公交", expenseRootType["行"])
	category.SetParentId("打车", expenseRootType["行"])
	category.SetParentId("火车", expenseRootType["行"])
	category.SetParentId("单车", expenseRootType["行"])

	category.SetParentId("投资亏损", expenseRootType["资金"])
	category.SetParentId("保险", expenseRootType["资金"])
	category.SetParentId("坏账流出", expenseRootType["资金"])

	category.SetParentId("房租", expenseRootType["住"])
	category.SetParentId("酒店", expenseRootType["住"])
}

// InitCategory 初始化分类数据
func InitCategory() {
	for k, v := range expenseRootType {
		tmpCategory := &category.Category{Name: k, Leaf: false, TypeId: constant.CategoryExpense}
		tmpCategory.ID = v
		category.AddCategory(tmpCategory)
	}

	for k, v := range expenseType {
		tmpCategory := &category.Category{Name: k, Leaf: true, TypeId: constant.CategoryExpense}
		tmpCategory.ID = v
		category.AddCategory(tmpCategory)
	}

	for k, v := range incomeType {
		tmpCategory := &category.Category{Name: k, Leaf: true, TypeId: constant.CategoryIncome}
		tmpCategory.ID = v
		category.AddCategory(tmpCategory)
	}

	for k, v := range transferType {
		tmpCategory := &category.Category{Name: k, Leaf: true, TypeId: constant.CategoryTransfer}
		tmpCategory.ID = v
		category.AddCategory(tmpCategory)
	}
}
