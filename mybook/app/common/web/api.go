package web

import (
	"mybook/app/account"
	"mybook/app/category"
	"mybook/app/common"
	"mybook/app/common/config"
	"mybook/app/loan"
	"mybook/app/record"
	"mybook/app/report"
	"mybook/app/user"

	"github.com/gin-gonic/gin"
)

// 服务端API注册
func registerServerApi(router *gin.Engine) {
	// 分类
	router.GET(buildApi("/category/listCategoryType"), common.ListCategoryType)
	router.GET(buildApi("/category/listCategory"), category.ListCategory)
	router.GET(buildApi("/category/listCategoryTree"), category.ListCategoryTree)

	// 账户
	router.GET(buildApi("/account/listAccount"), account.ListAccount)
	router.POST(buildApi("/account/createAccount"), account.CreateNewAccount)

	// 用户
	router.GET(buildApi("/user/listUser"), user.ListUser)
	router.GET(buildApi("/user/addUser"), user.AddUser)

	// 账单
	router.GET(buildApi("/record/calBalance"), record.CalculateAccountBalance)
	router.POST(buildApi("/record/createRecord"), record.CreateRecord)
	router.POST(buildApi("/loan/create"), loan.CreateLoan)

	router.GET(buildApi("/loan/query"), loan.QueryLoan)

	router.GET(buildApi("/record/listRecord"), record.ListRecord)
	router.GET(buildApi("/record/category"), record.CategoryRecord)

	router.GET(buildApi("/record/categoryDetail"), record.CategoryDetailRecord)
	router.GET(buildApi("/record/categoryWeekDetail"), record.WeekCategoryDetailRecord)
	router.GET(buildApi("/record/categoryMonthDetail"), record.MonthCategoryDetailRecord)

	router.GET(buildApi("/report/categoryPeriodReport"), report.CategoryPeriodReport) // 各分类周期报表
	router.GET(buildApi("/report/balanceReport"), report.BalanceReport)               // 余额报表
}

func buildApi(path string) string {
	return config.DefaultUrlPath + path
}
