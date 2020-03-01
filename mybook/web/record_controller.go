package web

import (
	"github.com/gin-gonic/gin"
	"github.com/kuangcp/gobase/mybook/service"
	"github.com/kuangcp/gobase/mybook/vo"
	"github.com/wonderivan/logger"
)

// curl -s -F typeId=4 -F accountId=3 -F categoryId=102 -F amount=2000 -F date='2020-02-03' localhost:10006/record | pretty-json
func CreateRecord(c *gin.Context) {
	typeId := c.PostForm("typeId")
	accountId := c.PostForm("accountId")
	categoryId := c.PostForm("categoryId")
	amount := c.PostForm("amount")
	date := c.PostForm("date")
	comment := c.PostForm("comment")

	logger.Debug(typeId, accountId, categoryId, amount, date, comment)

	record := service.BuildRecordByField(typeId, accountId, categoryId, amount, date, comment)
	logger.Debug(record)
	c.JSON(200, vo.SuccessWith(record))
}
