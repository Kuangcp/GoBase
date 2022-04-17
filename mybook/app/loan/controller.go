package loan

import (
	"github.com/gin-gonic/gin"
	"github.com/kuangcp/gobase/pkg/cuibase"
	"github.com/kuangcp/gobase/pkg/ghelp"
	"github.com/kuangcp/logger"
	"mybook/app/common/constant"
	"mybook/app/common/dal"
	"mybook/app/common/util"
	"mybook/app/record"
	"time"
)

type (
	CreateLoanParam struct {
		UserId       int    `json:"userId"`
		AccountId    int    `json:"accountId"`
		LoanType     int    `json:"loanType"`
		Amount       string `json:"amount"` // 支持多个金额输入 例如 21,13,6 最终会求和 ParseMultiPrice
		Date         string `json:"date"`
		ExceptedDate string `json:"exceptedDate"`
		Comment      string `json:"comment"`
	}
)

func CreateLoan(c *gin.Context) {
	var paramVO CreateLoanParam
	err := c.ShouldBind(&paramVO)
	if err != nil {
		ghelp.GinFailedWithMsg(c, err.Error())
		return
	}

	logger.Debug("param: ", util.Json(paramVO))
	resultVO, entity := buildEntityFromParam(paramVO)
	if resultVO.IsFailed() {
		ghelp.GinResultVO(c, resultVO)
		return
	}
	logger.Info(entity)

	reP := record.CreateParamVO{
		AccountId:       paramVO.AccountId,
		TargetAccountId: 0,
		Amount:          paramVO.Amount,
		CategoryId:      int(constant.CategoryTransferId),
		TypeId:          constant.RecordTransferOut,
		Date:            []string{paramVO.Date},
		Comment:         paramVO.Comment,
	}

	// 设置 账户关系
	if entity.LoanType == constant.LoanBorrow {
		reP.AccountId = constant.AccountAPId
		reP.TargetAccountId = paramVO.AccountId
	} else if entity.LoanType == constant.LoanLend {
		reP.AccountId = paramVO.AccountId
		reP.TargetAccountId = constant.AccountARId
	}
	logger.Info(reP)

	recordList := record.CreateMultipleTypeRecord(reP)
	if recordList.IsFailed() {
		ghelp.GinResultVO(c, recordList)
		return
	}

	entities := recordList.Data.([]*record.RecordEntity)
	logger.Info(entities)
	entity.TransferId = entities[0].TransferId
	db := dal.GetDB()
	db.Create(entity)

	ghelp.GinSuccess(c)
}

func buildEntityFromParam(paramVO CreateLoanParam) (ghelp.ResultVO, *Entity) {
	if paramVO.Date == "" || paramVO.Amount == "" || paramVO.UserId == 0 || paramVO.AccountId == 0 || paramVO.LoanType == 0 {
		return ghelp.FailedWithMsg("参数校验失败"), nil
	}
	recordDate, e := time.Parse(cuibase.YYYY_MM_DD, paramVO.Date)
	if e != nil {
		logger.Error(e)
		return ghelp.FailedWithMsg("date 参数错误"), nil
	}

	price := util.ParseMultiPrice(paramVO.Amount)
	if price.IsFailed() {
		return price, nil
	}
	entity := Entity{
		AccountId:  uint(paramVO.AccountId),
		UserId:     uint(paramVO.UserId),
		LoanType:   int8(paramVO.LoanType),
		RecordTime: recordDate,
		Comment:    paramVO.Comment,
		Amount:     price.Data.(int),
	}

	if paramVO.ExceptedDate != "" {
		exceptedDate, e := time.Parse(cuibase.YYYY_MM_DD, paramVO.ExceptedDate)
		if e != nil {
			logger.Error(e)
			return ghelp.FailedWithMsg("date 参数错误"), nil
		}
		entity.ExceptTime = exceptedDate
	}

	return ghelp.Success(), &entity
}
