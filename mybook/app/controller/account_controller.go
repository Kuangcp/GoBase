package controller

import (
	"github.com/gin-gonic/gin"
	"mybook/app/common/util"
	"mybook/app/dto"
	"mybook/app/service"
	"mybook/app/vo"
)

func ListAccount(c *gin.Context) {
	accounts := service.ListAccounts()
	result := util.Copy(accounts, new([]dto.AccountDTO)).(*[]dto.AccountDTO)
	vo.GinSuccessWith(c, result)
}

func CalculateAccountBalance(c *gin.Context) {
	vo.GinSuccessWith(c, service.CalculateAccountBalance())
}
