package record

import (
	"github.com/gin-gonic/gin"
	"github.com/kuangcp/gobase/pkg/ghelp"
	"github.com/wonderivan/logger"
	"mybook/app/common/util"
)

func CreateRecord(c *gin.Context) {
	var paramVO RecordCreateParamVO
	err := c.ShouldBind(&paramVO)
	if err != nil {
		ghelp.GinFailedWithMsg(c, err.Error())
		return
	}
	logger.Debug("createRecord param: ", util.Json(paramVO))

	result := createMultipleTypeRecord(paramVO)
	if result.IsFailed() {
		ghelp.GinResultVO(c, result)
		return
	}

	ghelp.GinResult(c, result.Data)
}

func ListRecord(c *gin.Context) {
	var param QueryRecordParam
	_ = c.ShouldBind(&param)
	result := findRecord(param)
	ghelp.GinResult(c, convertToVOList(result))
}

func CategoryRecord(c *gin.Context) {
	var query QueryRecordParam
	_ = c.ShouldBind(&query)

	result := queryCategoryRecord(query)
	ghelp.GinResult(c, result)
}

func CategoryDetailRecord(c *gin.Context) {
	var param QueryRecordParam
	_ = c.ShouldBind(&param)
	result := findRecord(param)
	ghelp.GinResult(c, convertToVOList(result))
}

func WeekCategoryDetailRecord(c *gin.Context) {
	var param QueryRecordParam
	_ = c.ShouldBind(&param)
	result := weekCategoryRecord(param)
	ghelp.GinResult(c, result)
}

func MonthCategoryDetailRecord(c *gin.Context) {
	var param QueryRecordParam
	_ = c.ShouldBind(&param)
	result := monthCategoryRecord(param)
	ghelp.GinResult(c, result)
}

func QueryAccountBalance(c *gin.Context) {
	ghelp.GinSuccessWith(c, calculateAndQueryAccountBalance())
}
