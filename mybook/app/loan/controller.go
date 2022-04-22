package loan

import (
	"github.com/gin-gonic/gin"
	"github.com/kuangcp/gobase/pkg/ghelp"
)

type (
	CreateLoanParam struct {
		UserId       int    `json:"userId"`
		AccountId    int    `json:"accountId"`
		LoanType     int    `json:"loanType"`
		Amount       string `json:"amount"` // 支持多个金额输入 例如 21,13,6 最终会求和 ParseMultiPrice
		Date         string `json:"date"`
		ExceptedDate string `json:"exceptedDate"`
		Comment      string `json:"comment"`
	}

	LoanUserVO struct {
		UserId uint
		Name   string
		Amount int
	}
)

func QueryLoan(c *gin.Context) {
	users := queryAllLoanUser()
	ghelp.GinSuccessWith(c, users)
}

func CreateLoan(c *gin.Context) {
	var paramVO CreateLoanParam
	err := c.ShouldBind(&paramVO)
	if err != nil {
		ghelp.GinFailedWithMsg(c, err.Error())
		return
	}

	ghelp.GinResultVO(c, createLoan(paramVO))
}
