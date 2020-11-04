package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/kuangcp/gobase/pkg/ginhelper"

	"mybook/app/common/util"
	"mybook/app/dto"
	"mybook/app/service"
)

func ListAccount(c *gin.Context) {
	accounts := service.ListAccounts()
	result := util.Copy(accounts, new([]dto.AccountDTO)).(*[]dto.AccountDTO)
	ginhelper.GinSuccessWith(c, result)
}

func CalculateAccountBalance(c *gin.Context) {
	ginhelper.GinSuccessWith(c, service.CalculateAccountBalance())
}
