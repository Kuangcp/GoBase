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

func ListAccount(c *gin.Context) {
	accounts := service.ListAccounts()
	vo.SuccessForWebWith(c, accounts)
}

func ListRecordType(c *gin.Context) {
	_, list := constant.GetRecordTypeMap()
	vo.SuccessForWebWith(c, list)
}

func ListCategory(c *gin.Context) {
	recordType := c.Query("recordType")
	if recordType == "" {
		vo.SuccessForWebWith(c, service.FindAllCategory())
	}
	i, _ := strconv.Atoi(recordType)
	typeEnum := constant.GetCategoryTypeByRecordTypeIndex(int8(i))
	if typeEnum != nil {
		list := service.FindCategoryByTypeId(typeEnum.Index)
		vo.SuccessForWebWith(c, list)
	} else {
		vo.FailedForWeb(c)
	}
}
