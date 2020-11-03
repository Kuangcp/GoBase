package common

import (
	"github.com/gin-gonic/gin"
	"mybook/app/common/constant"
	"mybook/app/service"
	"mybook/app/vo"
	"strconv"
)

// 简单查询

func HealthCheck(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func ListRecordType(c *gin.Context) {
	_, list := constant.GetRecordTypeMap()
	vo.GinSuccessWith(c, list)
}

func ListCategoryType(c *gin.Context) {
	_, list := constant.GetCategoryTypeMap()
	vo.GinSuccessWith(c, list)
}

func ListCategory(c *gin.Context) {
	recordType := c.Query("recordType")
	if recordType == "" {
		vo.GinSuccessWith(c, service.ListCategories())
		return
	}

	i, _ := strconv.Atoi(recordType)
	typeEnum := constant.GetCategoryTypeByRecordTypeIndex(int8(i))
	if typeEnum != nil {
		list := service.FindLeafCategoryByTypeId(typeEnum.Index)
		vo.GinSuccessWith(c, list)
	} else {
		vo.GinFailed(c)
	}
}
