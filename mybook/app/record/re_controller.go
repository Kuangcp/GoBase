package record

import (
	"mybook/app/common/util"

	"github.com/gin-gonic/gin"
	"github.com/kuangcp/gobase/pkg/ghelp"
	"github.com/kuangcp/logger"
)

type (
	CreateParamVO struct {
		AccountId       int      `json:"accountId"`
		TargetAccountId int      `json:"targetAccountId"`
		Amount          string   `json:"amount"` // 支持多个金额输入 例如 21,13,6 最终会求和 ParseMultiPrice
		CategoryId      int      `json:"categoryId"`
		TypeId          int8     `json:"typeId"` // TypeId 含义为 categoryTypeId
		Date            []string `json:"date"`
		Comment         string   `json:"comment"`
	}
	QueryRecordParam struct {
		AccountId  string `form:"accountId"`
		CategoryId string `form:"categoryId"`
		TypeId     string `form:"typeId"` // record_type
		StartDate  string `form:"startDate"`
		EndDate    string `form:"endDate"`
	}
)

func CreateRecord(c *gin.Context) {
	var paramVO CreateParamVO
	err := c.ShouldBind(&paramVO)
	if err != nil {
		ghelp.GinFailedWithMsg(c, err.Error())
		return
	}
	logger.Debug("createRecord param: ", util.Json(paramVO))

	result := CreateMultipleTypeRecord(paramVO)
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

func CalculateAccountBalance(c *gin.Context) {
	ghelp.GinSuccessWith(c, calculateAndQueryAccountBalance())
}
