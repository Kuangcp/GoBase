package web

import (
	"github.com/gin-gonic/gin"
	"github.com/kuangcp/gobase/mybook/constant"
	"github.com/kuangcp/gobase/mybook/service"
	"github.com/kuangcp/gobase/mybook/vo"
	"strconv"
)

func HealthCheck(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func ListRecordType(c *gin.Context) {
	_, list := constant.GetRecordTypeMap()
	c.JSON(200, vo.SuccessWith(list))
}

func ListCategory(c *gin.Context) {
	recordType := c.Query("recordType")
	if recordType == "" {
		c.JSON(200, vo.SuccessWith(service.FindAllCategory()))
	}
	i, _ := strconv.Atoi(recordType)
	typeEnum := constant.GetCategoryTypeByIndex(int8(i))
	if typeEnum != nil {
		list := service.FindCategoryByTypeId(typeEnum.Index)
		c.JSON(200, vo.SuccessWith(list))
	}
	c.JSON(200, vo.Failed())
}
