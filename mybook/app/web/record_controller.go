package web

import (
	"github.com/gin-gonic/gin"
	"github.com/kuangcp/gobase/mybook/app/service"
	"github.com/kuangcp/gobase/mybook/app/util"
	"github.com/kuangcp/gobase/mybook/app/vo"
	"github.com/wonderivan/logger"
)

// curl -s -F typeId=4 -F accountId=3 -F categoryId=102 -F amount=2000 -F date='2020-02-03' localhost:10006/record | pretty-json
func CreateRecord(c *gin.Context) {
	// typeId 含义为 categoryTypeId
	typeId := c.PostForm("typeId")
	accountId := c.PostForm("accountId")
	targetAccountId := c.PostForm("targetAccountId")
	categoryId := c.PostForm("categoryId")
	amount := c.PostForm("amount")
	date := c.PostForm("date")
	comment := c.PostForm("comment")

	recordVO := vo.CreateRecordVO{TypeId: typeId, AccountId: accountId, CategoryId: categoryId,
		Amount: amount, Date: date, Comment: comment, TargetAccountId: targetAccountId}

	logger.Debug("createRecord: ", util.Json(recordVO))

	record := service.CreateMultipleTypeRecord(recordVO)
	if record != nil {
		logger.Debug("createRecord result: ", util.Json(record))
		c.JSON(200, vo.SuccessWith(record))
	} else {
		c.JSON(200, vo.Failed())
	}
}

func ListRecord(c *gin.Context) {
	accountId := c.Query("accountId")
	startDate := c.Query("startDate")
	endDate := c.Query("endDate")
	typeId := c.Query("typeId")

	query := vo.QueryRecordVO{AccountId: accountId, StartDate: startDate, EndDate: endDate, TypeId: typeId}
	result := service.FindRecord(query)
	if result != nil {
		c.JSON(200, vo.SuccessWith(result))
	} else {
		c.JSON(200, vo.Failed())
	}
}
