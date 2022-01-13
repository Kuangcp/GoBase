package common

import (
	"mybook/app/account"
	"mybook/app/book"
	"mybook/app/category"
	"mybook/app/common/constant"
	"mybook/app/common/dal"
	"mybook/app/record"
)

func AutoMigrateAll() {
	db := dal.GetDB()
	db.AutoMigrate(&account.Account{})
	db.AutoMigrate(&category.Category{})
	db.AutoMigrate(&record.RecordEntity{})
	db.AutoMigrate(&book.BookKeeping{})

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
}
