package record

import (
	"mybook/app/account"
	"mybook/app/category"
	"mybook/app/common/constant"
	"mybook/app/common/dal"
	"mybook/app/common/util"
	"strconv"
	"strings"
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

func checkParam(record *RecordEntity) (ghelp.ResultVO, *category.Category, *account.Account) {
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
	resultVO, _, _ := checkParam(record)
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

	resultVO, _, _ := checkParam(origin)
	if resultVO.IsFailed() {
		return resultVO
	}
	resultVO, _, _ = checkParam(target)
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

func BuildRecordByField(param RecordCreateParamVO) ghelp.ResultVO {
	if len(param.Date) == 0 {
		return ghelp.FailedWithMsg("日期为空")
	}
	var recordList []*RecordEntity
	for _, date := range param.Date {
		recordDate, e := time.Parse("2006-01-02", date)
		if e != nil {
			logger.Error(e)
			return ghelp.FailedWithMsg("date 参数错误")
		}
		param.Amount = strings.Replace(param.Amount, "，", ",", -1)
		amountList := strings.Split(param.Amount, ",")
		var totalAmount = 0
		for _, one := range amountList {
			parseResult := parseAmount(one)
			if parseResult.IsFailed() {
				return parseResult
			}
			totalAmount += parseResult.Data.(int)
		}

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

func parseAmount(amount string) ghelp.ResultVO {
	floatAmount, e := strconv.ParseFloat(amount, 64)
	if e != nil {
		return ghelp.FailedWithMsg("amount 参数错误")
	}
	if floatAmount <= 0 {
		return ghelp.FailedWithMsg("amount 无效")
	}
	values := strings.Split(amount, ".")
	if len(values) > 1 {
		p := values[0]
		v := values[1]
		vLen := len(v)
		if vLen > 2 {
			return ghelp.FailedWithMsg("amount 仅保留两位小数")
		} else if vLen == 2 {
			pInt, e := strconv.Atoi(p)
			if e != nil {
				return ghelp.FailedWithMsg("amount 参数错误")
			}
			vInt, e := strconv.Atoi(v)
			if e != nil {
				return ghelp.FailedWithMsg("amount 参数错误")
			}
			return ghelp.SuccessWith(pInt*100 + vInt)
		} else if vLen == 1 {
			pInt, e := strconv.Atoi(p)
			if e != nil {
				return ghelp.FailedWithMsg("amount 参数错误")
			}
			vInt, e := strconv.Atoi(v)
			if e != nil {
				return ghelp.FailedWithMsg("amount 参数错误")
			}
			return ghelp.SuccessWith(pInt*100 + vInt*10)
		} else {
			return ghelp.FailedWithMsg("amount 参数错误")
		}
	} else {
		return ghelp.SuccessWith(int(floatAmount) * 100)
	}
}

func createMultipleTypeRecord(param RecordCreateParamVO) ghelp.ResultVO {
	result := BuildRecordByField(param)
	if result.IsFailed() {
		return result
	}

	list := result.Data.([]*RecordEntity)
	var successList []*RecordEntity
	var failResults []*RecordEntity

	for _, record := range list {
		if param.TargetAccountId != 0 && constant.IsTransferRecordType(record.Type) {
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

			checkResult, _, _ := checkParam(target)
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
