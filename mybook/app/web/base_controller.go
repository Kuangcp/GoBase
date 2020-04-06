package web

import (
	"github.com/gin-gonic/gin"
	"github.com/kuangcp/gobase/mybook/app/constant"
	"github.com/kuangcp/gobase/mybook/app/service"
	"github.com/kuangcp/gobase/mybook/app/util"
	"github.com/kuangcp/gobase/mybook/app/vo"
	"github.com/kuangcp/gobase/mybook/app/web/dto"
	"strconv"
)

// 简单查询

func HealthCheck(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func ListAccount(c *gin.Context) {
	accounts := service.ListAccounts()
	result := util.Copy(accounts, new([] dto.AccountDTO)).(*[]dto.AccountDTO)
	vo.SuccessForWebWith(c, result)
}

func ListRecordType(c *gin.Context) {
	_, list := constant.GetRecordTypeMap()
	vo.SuccessForWebWith(c, list)
}

func ListCategoryType(c *gin.Context) {
	_, list := constant.GetCategoryTypeMap()
	vo.SuccessForWebWith(c, list)
}

func ListCategory(c *gin.Context) {
	recordType := c.Query("recordType")
	if recordType == "" {
		vo.SuccessForWebWith(c, service.ListCategories())
		return
	}

	i, _ := strconv.Atoi(recordType)
	typeEnum := constant.GetCategoryTypeByRecordTypeIndex(int8(i))
	if typeEnum != nil {
		list := service.FindLeafCategoryByTypeId(typeEnum.Index)
		vo.SuccessForWebWith(c, list)
	} else {
		vo.FailedForWeb(c)
	}
}
