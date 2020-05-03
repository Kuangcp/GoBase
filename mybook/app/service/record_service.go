package service

import (
	"github.com/kuangcp/gobase/cuibase"
	"github.com/kuangcp/gobase/mybook/app/constant"
	"github.com/kuangcp/gobase/mybook/app/dal"
	"github.com/kuangcp/gobase/mybook/app/domain"
	"github.com/kuangcp/gobase/mybook/app/util"
	"github.com/kuangcp/gobase/mybook/app/vo"
	"github.com/kuangcp/gobase/mybook/app/web/dto"
	"github.com/wonderivan/logger"
	"strconv"
	"time"
	"unsafe"
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

func CreateIncomeRecordByParams(params []string) {
	cuibase.AssertParamCount(5, "参数缺失: -ri AccountId CategoryId Amount Date [Comment]")
	p := params[2:]
	p = append([]string{strconv.Itoa(int(constant.RECORD_INCOME))}, p...)
	record := buildRecordByParams(p)
	resultVO := CreateRecord(record)
	if resultVO.IsFailed() {
		logger.Error(resultVO)
	}
}

func CreateExpenseRecordByParams(params []string) {
	cuibase.AssertParamCount(5, "参数缺失: -re AccountId CategoryId Amount Date [Comment]")
	p := params[2:]
	p = append([]string{strconv.Itoa(int(constant.RECORD_EXPENSE))}, p...)
	record := buildRecordByParams(p)
	resultVO := CreateRecord(record)
	if resultVO.IsFailed() {
		logger.Error(resultVO)
	}
}

func CreateTransRecordByParams(params []string) {
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

func CreateRecordByParams(params []string) {
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

	recordVO := vo.CreateRecordParam{TypeId: params[0], AccountId: params[1], CategoryId: params[2],
		Amount: params[3], Date: params[4], Comment: comment}
	return BuildRecordByField(recordVO)
}

func BuildRecordByField(param vo.CreateRecordParam) *domain.Record {
	typeId, e := strconv.Atoi(param.TypeId)
	if e != nil || !constant.IsValidRecordType(int8(typeId)) {
		logger.Error(e)
		return nil
	}
	accountId, e := strconv.ParseUint(param.AccountId, 10, 64)
	if e != nil {
		logger.Error(e)
		return nil
	}
	categoryId, e := strconv.ParseUint(param.CategoryId, 10, 64)
	if e != nil {
		logger.Error(e)
		return nil
	}

	amount, e := strconv.Atoi(param.Amount)
	if e != nil {
		logger.Error(e)
		return nil
	}

	recordDate, e := time.Parse("2006-01-02", param.Date)
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
	if param.Comment != "" {
		record.Comment = param.Comment
	}

	return record
}

func CreateMultipleTypeRecord(param vo.CreateRecordParam) *domain.Record {
	record := BuildRecordByField(param)
	if record == nil {
		return nil
	}

	if param.TargetAccountId != "" && constant.IsTransferRecordType(record.Type) {
		record.Type = constant.RECORD_TRANSFER_OUT
		accountId, e := strconv.ParseUint(param.TargetAccountId, 10, 64)
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

func FindRecord(param vo.QueryRecordParam) *[]dto.RecordDTO {
	db := dal.GetDB()
	var lists []domain.Record
	accountId, _ := strconv.Atoi(param.AccountId)
	typeId, _ := strconv.Atoi(param.TypeId)
	categoryId, _ := strconv.Atoi(param.CategoryId)

	query := db.Where("deleted_at is null and record_time between ? and ?", param.StartDate, param.EndDate)
	record := domain.Record{}
	if typeId != 0 {
		record.Type = int8(typeId)
	}
	if accountId != 0 {
		record.AccountId = uint(accountId)
	}
	if categoryId != 0 {
		record.CategoryId = uint(categoryId)
	}

	if unsafe.Sizeof(record) != 0 {
		query = query.Where(&record)
	}
	query.Order("record_time DESC", true).Find(&lists)
	if len(lists) < 1 {
		return nil
	}

	accountMap := ListAccountMap()
	categoryMap := ListCategoryMap()
	var result []dto.RecordDTO
	for i := range lists {
		record := lists[i]
		ele := dto.RecordDTO{
			ID:             record.ID,
			RecordType:     record.Type,
			AccountName:    accountMap[record.AccountId].Name,
			CategoryName:   categoryMap[record.CategoryId].Name,
			RecordTypeName: constant.GetRecordTypeByIndex(record.Type).Name,
			RecordTime:     record.RecordTime,
			Amount:         record.Amount,
			Comment:        record.Comment,
		}

		result = append(result, ele)
	}
	return &result
}

// typeId record_type
func CategoryRecord(startDate string, endDate string, typeId string) *[]dto.MonthCategoryRecordDTO {
	db := dal.GetDB()
	var result []dto.MonthCategoryRecordDTO
	query := db.Table("record").
		Select("record.category_id, category.name, sum(amount) as amount,type").
		Joins("left join category on record.category_id = category.id")
	if len(typeId) != 0 {
		query = query.Where("record.type =?", typeId)
	}

	query = query.Where("record_time between ? and ?", startDate, endDate)
	query.Where("record.deleted_at is null").Group("category_id").Scan(&result)
	logger.Info(query.QueryExpr())
	if len(result) == 0 {
		return nil
	}

	for i := range result {
		recordDTO := &result[i]
		recordDTO.Date = startDate
		recordDTO.RecordTypeName = constant.GetRecordTypeByIndex(recordDTO.Type).Name
	}
	return &result
}

func WeekCategoryRecord(param vo.QueryRecordParam) *[]vo.RecordWeekVO {
	records := FindRecord(param)
	if records == nil {
		return nil
	}

	var result []vo.RecordWeekVO

	var temp *vo.RecordWeekVO = nil
	var lastAdded *vo.RecordWeekVO = nil
	endDateObj, err := time.Parse("2006-01-02", param.EndDate)
	if err != nil {
		return nil
	}

	var lastCachedWeek = util.WeekByDate(endDateObj)
	for _, recordDTO := range *records {
		recordTime := recordDTO.RecordTime
		curWeek := util.WeekByDate(recordTime)
		if curWeek == lastCachedWeek {
			if temp == nil {
				temp = &vo.RecordWeekVO{
					StartDate: recordTime.AddDate(0, 0, -int(recordTime.Weekday()-time.Sunday)).
						Format("2006-01-02"),
					EndDate: recordTime.AddDate(0, 0, int(time.Saturday-recordTime.Weekday())).
						Format("2006-01-02"),
					Amount: recordDTO.Amount,
				}
				result = append(result, *temp)
				logger.Info(recordTime.Format("2006-01-02"), curWeek)
			} else {
				temp.Amount += recordDTO.Amount
			}
		} else {
			lastCachedWeek = curWeek
			if temp != nil {
				result = append(result, *temp)
				lastAdded = temp
			}
			temp = &vo.RecordWeekVO{
				StartDate: recordTime.AddDate(0, 0, -int(recordTime.Weekday()-time.Sunday)).
					Format("2006-01-02"),
				EndDate: recordTime.AddDate(0, 0, int(time.Saturday-recordTime.Weekday())).
					Format("2006-01-02"),
				Amount: recordDTO.Amount,
			}
			logger.Info(recordTime.Format("2006-01-02"), curWeek)
		}
	}
	if temp != nil && lastAdded != temp {
		result = append(result, *temp)
	}
	if len(result) == 0 {
		return nil
	}
	return &result
}

func MonthCategoryRecord(param vo.QueryRecordParam) *[]vo.RecordWeekVO {
	records := FindRecord(param)
	if records == nil {
		return nil
	}

	var result []vo.RecordWeekVO

	var temp *vo.RecordWeekVO = nil
	var lastAdded *vo.RecordWeekVO = nil
	endDateObj, err := time.Parse("2006-01-02", param.EndDate)
	if err != nil {
		return nil
	}

	var lastCachedMonth = util.MonthByDate(endDateObj)
	for _, recordDTO := range *records {
		recordTime := recordDTO.RecordTime
		curMonth := util.MonthByDate(recordTime)
		if curMonth == lastCachedMonth {
			if temp == nil {
				temp = &vo.RecordWeekVO{
					StartDate: recordTime.AddDate(0, 0, -recordTime.Day()+1).
						Format("2006-01-02"),
					Amount: recordDTO.Amount,
				}
				result = append(result, *temp)
				logger.Debug(recordTime.Format("2006-01-02"), curMonth)
			} else {
				temp.Amount += recordDTO.Amount
			}
		} else {
			lastCachedMonth = curMonth
			if temp != nil {
				result = append(result, *temp)
				lastAdded = temp
			}
			temp = &vo.RecordWeekVO{
				StartDate: recordTime.AddDate(0, 0, -recordTime.Day()+1).
					Format("2006-01-02"),
				Amount: recordDTO.Amount,
			}
			logger.Debug(recordTime.Format("2006-01-02"), curMonth)
		}
	}
	if temp != nil && lastAdded != temp {
		result = append(result, *temp)
	}
	if len(result) == 0 {
		return nil
	}
	return &result
}
