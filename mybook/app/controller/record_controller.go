package controller

import (
	"mybook/app/common/util"
	"mybook/app/param"
	"mybook/app/service"

	"github.com/gin-gonic/gin"
	"github.com/kuangcp/gobase/pkg/ginhelper"
	"github.com/wonderivan/logger"
)

func CreateRecord(c *gin.Context) {
	recordVO := param.RecordCreateParamVO{}
	err := c.Bind(&recordVO)
	if err != nil {
		ginhelper.GinFailedWithMsg(c, err.Error())
		return
	}

	logger.Debug("createRecord param: ", util.Json(recordVO))

	result := service.CreateMultipleTypeRecord(recordVO)
	if result.IsFailed() {
		ginhelper.GinResultVO(c, result)
		return
	}
	ginhelper.GinResult(c, result.Data)
}

func ListRecord(c *gin.Context) {
	accountId := c.Query("accountId")
	startDate := c.Query("startDate")
	endDate := c.Query("endDate")
	typeId := c.Query("typeId")

	query := param.QueryRecordParam{AccountId: accountId, StartDate: startDate, EndDate: endDate, TypeId: typeId}
	result := service.FindRecord(query)
	ginhelper.GinResult(c, result)
}

func CategoryRecord(c *gin.Context) {
	startDate := c.Query("startDate")
	endDate := c.Query("endDate")
	typeId := c.Query("typeId")

	result := service.CategoryRecord(startDate, endDate, typeId)
	ginhelper.GinResult(c, result)
}

func CategoryDetailRecord(c *gin.Context) {
	result := service.FindRecord(buildCategoryQueryParam(c))
	ginhelper.GinResult(c, result)
}

func WeekCategoryDetailRecord(c *gin.Context) {
	result := service.WeekCategoryRecord(buildCategoryQueryParam(c))
	ginhelper.GinResult(c, result)
}

func MonthCategoryDetailRecord(c *gin.Context) {
	result := service.MonthCategoryRecord(buildCategoryQueryParam(c))
	ginhelper.GinResult(c, result)
}

func buildCategoryQueryParam(c *gin.Context) param.QueryRecordParam {
	startDate := c.Query("startDate")
	endDate := c.Query("endDate")
	categoryId := c.Query("categoryId")
	typeId := c.Query("typeId")

	return param.QueryRecordParam{StartDate: startDate, EndDate: endDate, CategoryId: categoryId, TypeId: typeId}
}
