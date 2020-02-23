package service

import (
	"encoding/json"
	"github.com/kuangcp/gobase/cuibase"
	"github.com/kuangcp/gobase/myth-bookkeeping/constant"
	"github.com/kuangcp/gobase/myth-bookkeeping/dal"
	"github.com/kuangcp/gobase/myth-bookkeeping/domain"
	"github.com/kuangcp/gobase/myth-bookkeeping/vo"
	"log"
	"strconv"
	"time"
)

func addRecord(record *domain.Record) {
	db := dal.GetDB()
	db.Create(record)
}

func addBatchRecordsWithTransaction(records ...*domain.Record) error {
	db := dal.GetDB()
	tx := db.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	for i := range records {
		if err := tx.Error; err != nil {
			return err
		}

		if err := tx.Create(records[i]).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit().Error
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

	e := addBatchRecordsWithTransaction(origin, target)
	if e != nil {
		log.Println(e)
		return vo.Failed()
	}
	return vo.Success()
}

func CreateExpenseRecordByParams(params [] string) {
	cuibase.AssertParamCount(5, "参数缺失: -re AccountId CategoryId Amount Time [Comment]")
	p := params[2:]
	p = append([]string{strconv.Itoa(int(constant.EXPENSE))}, p...)
	record := buildRecordByParams(p)
	resultVO := CreateRecord(record)
	if resultVO.IsFailed() {
		log.Println(resultVO)
	}
}

func CreateTransRecordByParams(params [] string) {
	cuibase.AssertParamCount(6, "参数缺失: -rt AccountId CategoryId Amount Time ToAccountId [Comment]")
	p := params[2:6]
	p = append([]string{strconv.Itoa(int(constant.TRANSFER_IN))}, p...)
	record := buildRecordByParams(p)
	if record == nil {
		return
	}
	accountId, e := strconv.ParseUint(params[6], 10, 64)
	if e != nil {
		log.Println(e)
		return
	}

	now := time.Now()
	record.TransferId = uint(now.UnixNano())
	aj, _ := json.Marshal(record)
	target := new(domain.Record)
	_ = json.Unmarshal(aj, target)
	target.AccountId = uint(accountId)
	checkResult, _, _ := checkParam(target)
	if checkResult.IsFailed() {
		return
	}

	createResult := createTransRecord(record, target)
	if createResult.IsFailed() {
		log.Println(createResult)
	}
}

func CreateRecordByParams(params [] string) {
	cuibase.AssertParamCount(6, "参数缺失: -r TypeId AccountId CategoryId Amount Time [Comment]")
	record := buildRecordByParams(params[2:])
	resultVO := CreateRecord(record)
	if resultVO.IsFailed() {
		log.Println(resultVO)
	}
}

// params: TypeId AccountId CategoryId Amount Time [Comment]
func buildRecordByParams(params []string) *domain.Record {
	typeId, e := strconv.Atoi(params[0])
	if e != nil || !constant.IsValidRecordType(int8(typeId)) {
		log.Println(e)
		return nil
	}
	accountId, e := strconv.ParseUint(params[1], 10, 64)
	if e != nil {
		log.Println(e)
		return nil
	}
	categoryId, e := strconv.ParseUint(params[2], 10, 64)
	if e != nil {
		log.Println(e)
		return nil
	}

	amount, e := strconv.Atoi(params[3])
	if e != nil {
		log.Println(e)
		return nil
	}

	recordTime, e := time.Parse("2006-01-02", params[4])
	if e != nil {
		log.Println(e)
		return nil
	}

	record := &domain.Record{
		AccountId:  uint(accountId),
		CategoryId: uint(categoryId),
		Type:       int8(typeId),
		Amount:     amount,
		RecordTime: recordTime,
	}
	if len(params) == 6 {
		record.Comment = params[5]
	}

	return record
}
