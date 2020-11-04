package common

import (
	"mybook/app/common/constant"
	"mybook/app/service"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kuangcp/gobase/pkg/ginhelper"
)

// 简单查询

func HealthCheck(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func ListRecordType(c *gin.Context) {
	_, list := constant.GetRecordTypeMap()
	ginhelper.GinSuccessWith(c, list)
}

func ListCategoryType(c *gin.Context) {
	_, list := constant.GetCategoryTypeMap()
	ginhelper.GinSuccessWith(c, list)
}

func ListCategory(c *gin.Context) {
	recordType := c.Query("recordType")
	if recordType == "" {
		ginhelper.GinSuccessWith(c, service.ListCategories())
		return
	}

	i, _ := strconv.Atoi(recordType)
	typeEnum := constant.GetCategoryTypeByRecordTypeIndex(int8(i))
	if typeEnum != nil {
		list := service.FindLeafCategoryByTypeId(typeEnum.Index)
		ginhelper.GinSuccessWith(c, list)
	} else {
		ginhelper.GinFailed(c)
	}
}
