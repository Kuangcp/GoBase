package loan

import (
	"mybook/app/common/constant"
	"mybook/app/common/dal"
	"mybook/app/common/util"
	"mybook/app/record"
	"mybook/app/user"
	"sort"
	"time"

	"github.com/kuangcp/gobase/pkg/cuibase"
	"github.com/kuangcp/gobase/pkg/ghelp"
	"github.com/kuangcp/logger"
)

func queryAllLoanUser() []LoanUserVO {
	db := dal.GetDB()
	var entities []*Entity

	db.Table("record_loan").
		Select("user_id,loan_type,sum(amount) as amount").
		Group("user_id,loan_type").
		Find(&entities)

	var users []LoanUserVO
	if len(entities) == 0 {
		return users
	}

	userMap := user.QueryUserMap()
	amountMap := make(map[uint]int)
	for _, entity := range entities {
		userId := entity.UserId
		_, ok := amountMap[userId]
		if !ok {
			amountMap[userId] = 0
		}

		if entity.LoanType == constant.LoanBorrow || entity.LoanType == constant.LoanLendRe {
			amountMap[userId] = amountMap[userId] - entity.Amount
		} else {
			amountMap[userId] = amountMap[userId] + entity.Amount
		}
	}

	for k, v := range amountMap {
		users = append(users, LoanUserVO{
			UserId: k,
			Name:   userMap[k].Name,
			Amount: v,
		})
	}
	sort.Slice(users, func(i, j int) bool {
		return users[i].Amount > users[j].Amount
	})
	return users
}

func createLoan(paramVO CreateLoanParam) ghelp.ResultVO {
	resultVO, entity := buildEntityFromParam(paramVO)
	if resultVO.IsFailed() {
		logger.Error(resultVO)
		return resultVO
	}

	multiRecordParam := record.CreateParamVO{
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
		multiRecordParam.AccountId = constant.AccountAPId
		multiRecordParam.TargetAccountId = paramVO.AccountId
	} else if entity.LoanType == constant.LoanLend {
		multiRecordParam.AccountId = paramVO.AccountId
		multiRecordParam.TargetAccountId = constant.AccountARId
	} else if entity.LoanType == constant.LoanBorrowRe {
		multiRecordParam.AccountId = paramVO.AccountId
		multiRecordParam.TargetAccountId = constant.AccountAPId
	} else if entity.LoanType == constant.LoanLendRe {
		multiRecordParam.AccountId = constant.AccountARId
		multiRecordParam.TargetAccountId = paramVO.AccountId
	}

	logger.Info(entity, multiRecordParam)
	recordList := record.CreateMultipleTypeRecord(multiRecordParam)
	if recordList.IsFailed() {
		logger.Error(recordList)
		return recordList
	}

	entities := recordList.Data.([]*record.RecordEntity)
	entity.TransferId = entities[0].TransferId
	db := dal.GetDB()
	db.Create(entity)
	return ghelp.Success()
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
