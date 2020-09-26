package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/kuangcp/gobase/mybook/app/vo"
)

type (
	RecordQueryParam struct {
		StartDate string
		EndDate   string
	}
)

func LineMap(c *gin.Context) {
	param := buildParam(c)
	if param == nil {
		vo.GinFailedWithMsg(c, "invalid param")
		return
	}
}

func buildParam(c *gin.Context) *RecordQueryParam {
	startDate := c.Query("startDate")
	endDate := c.Query("endDate")

	if startDate == "" || endDate == "" {
		return nil
	}
	return &RecordQueryParam{StartDate: startDate, EndDate: endDate}
}
