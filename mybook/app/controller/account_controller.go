package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/kuangcp/gobase/mybook/app/dto"
	"github.com/kuangcp/gobase/mybook/app/service"
	"github.com/kuangcp/gobase/mybook/app/util"
	"github.com/kuangcp/gobase/mybook/app/vo"
)

func ListAccount(c *gin.Context) {
	accounts := service.ListAccounts()
	result := util.Copy(accounts, new([]dto.AccountDTO)).(*[]dto.AccountDTO)
	vo.SuccessForWebWith(c, result)
}

func AccountBalance(c *gin.Context) {
	vo.FillResult(c, nil)
}
