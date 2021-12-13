package account

import (
	"mybook/app/common/util"

	"github.com/gin-gonic/gin"
	"github.com/kuangcp/gobase/pkg/ghelp"
)

func ListAccount(c *gin.Context) {
	accounts := ListAccounts()
	result := util.Copy(accounts, new([]AccountDTO)).(*[]AccountDTO)
	ghelp.GinSuccessWith(c, result)
}

func CreateNewAccount(c *gin.Context) {
	//var param createAccountParam
	//err := c.ShouldBind(&param)
	//if err != nil {
	//	ghelp.GinFailedWithMsg(c, "参数解析失败")
	//	return
	//}
	//
	//var account Account
	//util.Copy(param, &account)
	//logger.Info(param, " -> ", account)
	//AddAccount(&account)
	ghelp.GinSuccess(c)
}
