package service

import (
	"github.com/kuangcp/gobase/cuibase"
	"github.com/kuangcp/gobase/myth-bookkeeping/constant"
	"github.com/kuangcp/gobase/myth-bookkeeping/dal"
	"github.com/kuangcp/gobase/myth-bookkeeping/domain"
	"github.com/kuangcp/gobase/myth-bookkeeping/vo"
	"log"
	"strconv"
)

func addRecord(record *domain.Record) {
	db := dal.GetDB()
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
	resultVO, _, _ := checkParam(record)
	if resultVO.IsFailed() {
		return resultVO
	}

	addRecord(record)
	return vo.Success()
}

func CreateRecordByParams(params [] string) {
	cuibase.AssertParamCount(5, "参数缺失: AccountId CategoryId Type Amount Comment")
	accountId, e := strconv.ParseUint(params[2], 10, 64)
	if e != nil {
		log.Println(e)
		return
	}
	categoryId, e := strconv.ParseUint(params[3], 10, 64)
	if e != nil {
		log.Println(e)
		return
	}
	typeId, e := strconv.Atoi(params[4])
	if e != nil {
		log.Println(e)
		return
	}
	amount, e := strconv.Atoi(params[5])
	if e != nil {
		log.Println(e)
		return
	}

	record := &domain.Record{
		AccountId:  uint(accountId),
		CategoryId: uint(categoryId),
		Type:       int8(typeId),
		Amount:     amount,
	}
	if len(params) == 7 {
		record.Comment = params[6]
	}

	resultVO := CreateRecord(record)
	if resultVO.IsFailed() {
		log.Println(resultVO)
	}
}
