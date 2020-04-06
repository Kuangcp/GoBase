package main

import (
	"github.com/kuangcp/gobase/mybook/app/constant"
	"github.com/kuangcp/gobase/mybook/app/domain"
	"github.com/kuangcp/gobase/mybook/app/service"
	"testing"
)

func TestCreateDBAndTable(t *testing.T) {
	service.AutoMigrateAll()
}

func TestInitBook(t *testing.T) {
	service.AddBookkeeping(&domain.BookKeeping{Name: "主账本1", Comment: ""})
}

func TestInitAccount(t *testing.T) {
	service.AddAccount(&domain.Account{TypeId: constant.ACCOUNT_CASH, Name: "现金", InitAmount: 0})
	service.AddAccount(&domain.Account{TypeId: constant.ACCOUNT_CREDIT, Name: "花呗", InitAmount: 0, MaxAmount: 2000, BillDay: 1, RepaymentDay: 10})
	service.AddAccount(&domain.Account{TypeId: constant.ACCOUNT_ONLINE, Name: "支付宝", InitAmount: 0})
	service.AddAccount(&domain.Account{TypeId: constant.ACCOUNT_ONLINE, Name: "微信", InitAmount: 0})
	service.AddAccount(&domain.Account{TypeId: constant.ACCOUNT_DEPOSIT, Name: "储蓄卡", InitAmount: 0})
}

func TestInitCategory(t *testing.T) {
	expenseIndex := 100
	var types = []string{"日常餐", "文娱", "日常开支", "交通"}
	for e := range types {
		expenseIndex++
		category := &domain.Category{Name: types[e], Leaf: false, TypeId: constant.CATEGORY_EXPENSE}
		category.ID = uint(expenseIndex)
		service.AddCategory(category)
	}

	types = []string{"早餐", "午餐", "晚餐", "餐厅", "零食", "日用品", "室外娱乐", "服饰", "云服务", "水果", "买菜",
		"发红包", "房租", "书籍", "话费网费", "火车", "数码", "礼物", "地铁", "酒店", "医疗", "公交", "打车", "知识付费", "坏账",
		"景点门票", "会员", "水电煤", "美容美发", "快递", "投资亏损", "电影", "保险", "打赏", "还贷", "己方借出"}
	for e := range types {
		expenseIndex++
		category := &domain.Category{Name: types[e], Leaf: true, TypeId: constant.CATEGORY_EXPENSE}
		category.ID = uint(expenseIndex)
		service.AddCategory(category)
	}

	incomeIndex := 200
	types = []string{"应收款", "收红包", "投资收益", "返现", "工资", "个缴社保", "司缴社保", "平台", "奖金", "兼职", "生活费",
		"报销流入", "其他收入", "退款", "问题", "借贷", "对方归还"}
	for e := range types {
		incomeIndex++
		category := &domain.Category{Name: types[e], Leaf: true, TypeId: constant.CATEGORY_INCOME}
		category.ID = uint(incomeIndex)
		service.AddCategory(category)
	}

	transferIndex := 300
	types = []string{"转账", "加仓", "平仓"}
	for e := range types {
		transferIndex++
		category := &domain.Category{Name: types[e], Leaf: true, TypeId: constant.CATEGORY_TRANSFER}
		category.ID = uint(transferIndex)
		service.AddCategory(category)
	}
}

func TestSetParentId(t *testing.T) {
	service.SetParentId("早餐", 1)
}
