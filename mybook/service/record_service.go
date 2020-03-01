package service

import (
	"github.com/kuangcp/gobase/cuibase"
	"github.com/kuangcp/gobase/mybook/constant"
	"github.com/kuangcp/gobase/mybook/dal"
	"github.com/kuangcp/gobase/mybook/domain"
	"github.com/kuangcp/gobase/mybook/util"
	"github.com/kuangcp/gobase/mybook/vo"
	"github.com/wonderivan/logger"
	"strconv"
	"time"
)

func addRecord(record *domain.Record) {
	db := dal.GetDB()
	// TODO support multiple book
	record.BookId = 1
	db.Create(record)
}

func checkParam(record *domain.Record) (vo.ResultVO, *domain.Category, *domain.Account) {
	category := FindCategoryById(record.CategoryId)
	if category == nil || !category.Leaf {
		return vo.FailedWithMsg("分类id无效"), nil, nil
	}

	account := FindAccountById(record.AccountId)
	if account == nil {
		return vo.FailedWithMsg("账户无效"), category, nil
	}

	if record.Amount <= 0 {
		return vo.FailedWithMsg("金额无效"), category, account
	}
	if !constant.IsValidRecordType(record.Type) {
		return vo.FailedWithMsg("类别无效"), category, account
	}
	return vo.Success(), category, account
}

func CreateRecord(record *domain.Record) vo.ResultVO {
	if nil == record {
		return vo.Failed()
	}
	resultVO, _, _ := checkParam(record)
	if resultVO.IsFailed() {
		return resultVO
	}

	addRecord(record)
	return vo.Success()
}

func createTransRecord(origin *domain.Record, target *domain.Record) vo.ResultVO {
	if nil == origin || nil == target {
		return vo.Failed()
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
		return vo.Failed()
	}
	return vo.Success()
}

func CreateIncomeRecordByParams(params [] string) {
	cuibase.AssertParamCount(5, "参数缺失: -ri AccountId CategoryId Amount Date [Comment]")
	p := params[2:]
	p = append([]string{strconv.Itoa(int(constant.RECORD_INCOME))}, p...)
	record := buildRecordByParams(p)
	resultVO := CreateRecord(record)
	if resultVO.IsFailed() {
		logger.Error(resultVO)
	}
}

func CreateExpenseRecordByParams(params [] string) {
	cuibase.AssertParamCount(5, "参数缺失: -re AccountId CategoryId Amount Date [Comment]")
	p := params[2:]
	p = append([]string{strconv.Itoa(int(constant.RECORD_EXPENSE))}, p...)
	record := buildRecordByParams(p)
	resultVO := CreateRecord(record)
	if resultVO.IsFailed() {
		logger.Error(resultVO)
	}
}

func CreateTransRecordByParams(params [] string) {
	cuibase.AssertParamCount(6, "参数缺失: -rt OutAccountId CategoryId Amount Date InAccountId [Comment]")
	p := params[2:6]
	p = append([]string{strconv.Itoa(int(constant.RECORD_TRANSFER_OUT))}, p...)
	record := buildRecordByParams(p)
	if record == nil {
		return
	}
	accountId, e := strconv.ParseUint(params[6], 10, 64)
	if e != nil {
		logger.Error(e)
		return
	}

	now := time.Now()
	record.TransferId = uint(now.UnixNano())

	target := util.Copy(record, new(domain.Record)).(*domain.Record)
	if target == nil {
		return
	}

	target.AccountId = uint(accountId)
	target.Type = constant.RECORD_TRANSFER_IN

	checkResult, _, _ := checkParam(target)
	if checkResult.IsFailed() {
		logger.Error(checkResult)
		return
	}

	createResult := createTransRecord(record, target)
	if createResult.IsFailed() {
		logger.Error(createResult)
	}
}

func CreateRecordByParams(params [] string) {
	cuibase.AssertParamCount(6, "参数缺失: -r TypeId AccountId CategoryId Amount Date [Comment]")
	record := buildRecordByParams(params[2:])
	resultVO := CreateRecord(record)
	if resultVO.IsFailed() {
		logger.Error(resultVO)
	}
}

// params: TypeId AccountId CategoryId Amount Date [Comment]
func buildRecordByParams(params []string) *domain.Record {
	comment := ""
	if len(params) == 6 {
		comment = params[5]
	}

	recordVO := vo.RecordVO{TypeId: params[0], AccountId: params[1], CategoryId: params[2],
		Amount: params[3], Date: params[4], Comment: comment}
	return BuildRecordByField(recordVO)
}

func BuildRecordByField(recordVO vo.RecordVO) *domain.Record {
	typeId, e := strconv.Atoi(recordVO.TypeId)
	if e != nil || !constant.IsValidRecordType(int8(typeId)) {
		logger.Error(e)
		return nil
	}
	accountId, e := strconv.ParseUint(recordVO.AccountId, 10, 64)
	if e != nil {
		logger.Error(e)
		return nil
	}
	categoryId, e := strconv.ParseUint(recordVO.CategoryId, 10, 64)
	if e != nil {
		logger.Error(e)
		return nil
	}

	amount, e := strconv.Atoi(recordVO.Amount)
	if e != nil {
		logger.Error(e)
		return nil
	}

	recordDate, e := time.Parse("2006-01-02", recordVO.Date)
	if e != nil {
		logger.Error(e)
		return nil
	}

	record := &domain.Record{
		AccountId:  uint(accountId),
		CategoryId: uint(categoryId),
		Type:       int8(typeId),
		Amount:     amount,
		RecordTime: recordDate,
	}
	if recordVO.Comment != "" {
		record.Comment = recordVO.Comment
	}

	return record
}

func CreateMultipleTypeRecord(recordVO vo.RecordVO) *domain.Record {
	record := BuildRecordByField(recordVO)
	if record == nil {
		return nil
	}

	if recordVO.TargetAccountId != "" && constant.IsTransferRecordType(record.Type) {
		record.Type = constant.RECORD_TRANSFER_OUT
		accountId, e := strconv.ParseUint(recordVO.TargetAccountId, 10, 64)
		if e != nil {
			logger.Error(e)
			return nil
		}

		now := time.Now()
		record.TransferId = uint(now.UnixNano())

		target := util.Copy(record, new(domain.Record)).(*domain.Record)
		if target == nil {
			return nil
		}

		target.AccountId = uint(accountId)
		target.Type = constant.RECORD_TRANSFER_IN

		checkResult, _, _ := checkParam(target)
		if checkResult.IsFailed() {
			logger.Error(checkResult)
			return nil
		}

		createResult := createTransRecord(record, target)
		if createResult.IsFailed() {
			logger.Error(createResult)
			return nil
		}
		return record
	} else {
		resultVO := CreateRecord(record)
		if resultVO.IsFailed() {
			logger.Error(resultVO)
			return nil
		}
		return record
	}
}
