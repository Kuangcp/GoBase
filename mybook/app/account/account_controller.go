package account

import (
	"github.com/gin-gonic/gin"
	"github.com/kuangcp/gobase/pkg/ghelp"
	"mybook/app/common/util"
)

func ListAccount(c *gin.Context) {
	accounts := ListAccounts()
	result := util.Copy(accounts, new([]AccountDTO)).(*[]AccountDTO)
	ghelp.GinSuccessWith(c, result)
}
