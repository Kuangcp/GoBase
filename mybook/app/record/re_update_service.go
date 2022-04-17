package record

import (
	"github.com/kuangcp/gobase/pkg/cuibase"
	"mybook/app/account"
	"mybook/app/category"
	"mybook/app/common/constant"
	"mybook/app/common/dal"
	"mybook/app/common/util"
	"time"

	"github.com/kuangcp/gobase/pkg/ghelp"

	"github.com/kuangcp/logger"
)

func addRecord(record *RecordEntity) {
	db := dal.GetDB()
	// TODO support multiple book
	record.BookId = 1
	db.Create(record)
}

func checkRecordValid(record *RecordEntity) (ghelp.ResultVO, *category.Category, *account.Account) {
	categoryEntity := category.FindCategoryById(record.CategoryId)
	if categoryEntity == nil || !categoryEntity.Leaf {
		return ghelp.FailedWithMsg("分类id无效"), nil, nil
	}

	accountEntity := account.FindAccountById(record.AccountId)
	if accountEntity == nil {
		return ghelp.FailedWithMsg("账户无效"), categoryEntity, nil
	}

	if record.Amount <= 0 {
		return ghelp.FailedWithMsg("金额无效"), categoryEntity, accountEntity
	}
	if !constant.IsValidRecordType(record.Type) {
		return ghelp.FailedWithMsg("类别无效"), categoryEntity, accountEntity
	}
	return ghelp.Success(), categoryEntity, accountEntity
}

func DoCreateRecord(record *RecordEntity) ghelp.ResultVO {
	if nil == record {
		return ghelp.Failed()
	}
	resultVO, _, _ := checkRecordValid(record)
	if resultVO.IsFailed() {
		return resultVO
	}

	addRecord(record)
	return ghelp.Success()
}

func createTransRecord(origin *RecordEntity, target *RecordEntity) ghelp.ResultVO {
	if nil == origin || nil == target {
		return ghelp.Failed()
	}

	resultVO, _, _ := checkRecordValid(origin)
	if resultVO.IsFailed() {
		return resultVO
	}
	resultVO, _, _ = checkRecordValid(target)
	if resultVO.IsFailed() {
		return resultVO
	}

	e := dal.BatchSaveWithTransaction(origin, target)
	if e != nil {
		logger.Error(e)
		return ghelp.Failed()
	}
	return ghelp.Success()
}

func buildRecordListFromParam(param CreateParamVO) ghelp.ResultVO {
	if len(param.Date) == 0 {
		return ghelp.FailedWithMsg("日期为空")
	}
	var recordList []*RecordEntity

	// 多日期 同一个金额 和 其他所有帐目细节
	for _, date := range param.Date {
		recordDate, e := time.Parse(cuibase.YYYY_MM_DD, date)
		if e != nil {
			logger.Error(e)
			return ghelp.FailedWithMsg("date 参数错误")
		}

		priceRe := util.ParseMultiPrice(param.Amount)
		if priceRe.IsFailed() {
			return priceRe
		}
		var totalAmount = priceRe.Data.(int)

		record := &RecordEntity{
			AccountId:  uint(param.AccountId),
			CategoryId: uint(param.CategoryId),
			Type:       param.TypeId,
			Amount:     totalAmount,
			RecordTime: recordDate,
		}
		if param.Comment != "" {
			record.Comment = param.Comment
		}
		recordList = append(recordList, record)
	}

	return ghelp.SuccessWith(recordList)
}

func CreateMultipleTypeRecord(param CreateParamVO) ghelp.ResultVO {
	result := buildRecordListFromParam(param)
	if result.IsFailed() {
		return result
	}

	list := result.Data.([]*RecordEntity)
	var successList []*RecordEntity
	var failResults []*RecordEntity

	for _, record := range list {
		// 转账
		if param.TargetAccountId != 0 &&
			constant.IsTransferRecordType(record.Type) {
			record.Type = constant.RecordTransferOut

			now := time.Now()
			record.TransferId = uint(now.UnixNano())

			target := util.Copy(record, new(RecordEntity)).(*RecordEntity)
			if target == nil {
				failResults = append(failResults, record)
				logger.Error("复制失败")
				continue
			}

			target.AccountId = uint(param.TargetAccountId)
			target.Type = constant.RecordTransferIn

			checkResult, _, _ := checkRecordValid(target)
			if checkResult.IsFailed() {
				logger.Error(checkResult)
				failResults = append(failResults, record)
				continue
			}

			createResult := createTransRecord(record, target)
			if createResult.IsFailed() {
				logger.Error(createResult)
				failResults = append(failResults, record)
				continue
			}
			successList = append(successList, record)
		} else {
			// 支出或收入
			resultVO := DoCreateRecord(record)
			if resultVO.IsFailed() {
				logger.Error(resultVO)
				return resultVO
			}
			successList = append(successList, record)
		}
	}
	if len(failResults) != 0 {
		logger.Error(failResults)
	}
	return ghelp.SuccessWith(successList)
}
