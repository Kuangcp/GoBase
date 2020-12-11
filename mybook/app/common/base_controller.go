package common

import (
	"mybook/app/common/constant"

	"github.com/gin-gonic/gin"
	"github.com/kuangcp/gobase/pkg/ghelp"
)

// 简单查询

func HealthCheck(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func ListRecordType(c *gin.Context) {
	_, list := constant.GetRecordTypeMap()
	ghelp.GinSuccessWith(c, list)
}

func ListCategoryType(c *gin.Context) {
	_, list := constant.GetCategoryTypeMap()
	ghelp.GinSuccessWith(c, list)
}
