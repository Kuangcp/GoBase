package service

import (
	"fmt"
	"mybook/app/common/constant"
	"mybook/app/common/dal"
	"mybook/app/common/util"
	"mybook/app/domain"
	"mybook/app/dto"
	"mybook/app/param"
	"mybook/app/vo"
	"strconv"
	"time"
	"unsafe"

	"github.com/kuangcp/gobase/pkg/ginhelper"
	"github.com/wonderivan/logger"
)

func addRecord(record *domain.Record) {
	db := dal.GetDB()
	// TODO support multiple book
	record.BookId = 1
	db.Create(record)
}

func checkParam(record *domain.Record) (ginhelper.ResultVO, *domain.Category, *domain.Account) {
	category := FindCategoryById(record.CategoryId)
	if category == nil || !category.Leaf {
		return ginhelper.FailedWithMsg("分类id无效"), nil, nil
	}

	account := FindAccountById(record.AccountId)
	if account == nil {
		return ginhelper.FailedWithMsg("账户无效"), category, nil
	}

	if record.Amount <= 0 {
		return ginhelper.FailedWithMsg("金额无效"), category, account
	}
	if !constant.IsValidRecordType(record.Type) {
		return ginhelper.FailedWithMsg("类别无效"), category, account
	}
	return ginhelper.Success(), category, account
}

func CreateRecord(record *domain.Record) ginhelper.ResultVO {
	if nil == record {
		return ginhelper.Failed()
	}
	resultVO, _, _ := checkParam(record)
	if resultVO.IsFailed() {
		return resultVO
	}

	addRecord(record)
	return ginhelper.Success()
}

func createTransRecord(origin *domain.Record, target *domain.Record) ginhelper.ResultVO {
	if nil == origin || nil == target {
		return ginhelper.Failed()
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
		return ginhelper.Failed()
	}
	return ginhelper.Success()
}

func BuildRecordByField(param param.CreateRecordParam) *domain.Record {
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

func CreateMultipleTypeRecord(param param.CreateRecordParam) (*domain.Record,error) {
	record := BuildRecordByField(param)
	if record == nil {
		return nil, fmt.Errorf("构建失败")
	}

	if param.TargetAccountId != "" && constant.IsTransferRecordType(record.Type) {
		record.Type = constant.RecordTransferOut
		accountId, e := strconv.ParseUint(param.TargetAccountId, 10, 64)
		if e != nil {
			logger.Error(e)
			return nil, e
		}

		now := time.Now()
		record.TransferId = uint(now.UnixNano())

		target := util.Copy(record, new(domain.Record)).(*domain.Record)
		if target == nil {
			return nil, fmt.Errorf("复制失败")
		}

		target.AccountId = uint(accountId)
		target.Type = constant.RecordTransferIn

		checkResult, _, _ := checkParam(target)
		if checkResult.IsFailed() {
			logger.Error(checkResult)
			return nil, fmt.Errorf("记录校验失败")
		}

		createResult := createTransRecord(record, target)
		if createResult.IsFailed() {
			logger.Error(createResult)
			return nil, fmt.Errorf("转账记录校验失败")
		}
		return record,nil
	} else {
		resultVO := CreateRecord(record)
		if resultVO.IsFailed() {
			logger.Error(resultVO)
			return nil, fmt.Errorf("参数校验失败")
		}
		return record, nil
	}
}

func FindRecord(param param.QueryRecordParam) *[]dto.RecordDTO {
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

// 帐目按类型分组 typeId record_type
func CategoryRecord(startDate string, endDate string, typeId string) *[]interface{} {
	db := dal.GetDB()
	var result []dto.MonthCategoryRecordDTO
	query := db.Table("record").
		Select("record.category_id, category.name, sum(amount) as amount,type").
		Joins("left join category on record.category_id = category.id")
	if len(typeId) != 0 {
		query = query.Where("record.type =?", typeId)
	}

	query = query.Where("record_time between ? and ?", startDate, endDate).
		Where("record.deleted_at is null").
		Group("category_id").Scan(&result)
	if len(result) == 0 {
		return nil
	}

	var temp []interface{}
	for i := range result {
		recordDTO := &result[i]

		recordDTO.Date = startDate
		recordDTO.RecordTypeName = constant.GetRecordTypeByIndex(recordDTO.Type).Name

		temp = append(temp, recordDTO)
	}

	util.Sort(util.SortWrapper{
		Data: temp,
		CompareLessFunc: func(a interface{}, b interface{}) bool {
			return a.(*dto.MonthCategoryRecordDTO).Amount < b.(*dto.MonthCategoryRecordDTO).Amount
		},
		Reverse: true,
	})

	return &temp
}

func WeekCategoryRecord(param param.QueryRecordParam) *[]vo.RecordWeekOrMonthVO {
	records := FindRecord(param)
	if records == nil {
		return nil
	}
	endDateObj, err := time.Parse("2006-01-02", param.EndDate)
	if err != nil {
		return nil
	}

	builder := func(recordDTO dto.RecordDTO) *vo.RecordWeekOrMonthVO {
		recordTime := recordDTO.RecordTime
		return &vo.RecordWeekOrMonthVO{
			StartDate: recordTime.AddDate(0, 0, -int(recordTime.Weekday()-time.Sunday)).
				Format("2006-01-02"),
			EndDate: recordTime.AddDate(0, 0, int(time.Saturday-recordTime.Weekday())).
				Format("2006-01-02"),
			Amount: recordDTO.Amount,
		}
	}
	return buildCommonWeekOrMonthVO(endDateObj, records, util.WeekByDate, builder)
}

func MonthCategoryRecord(param param.QueryRecordParam) *[]vo.RecordWeekOrMonthVO {
	records := FindRecord(param)
	if records == nil {
		return nil
	}
	endDateObj, err := time.Parse("2006-01-02", param.EndDate)
	if err != nil {
		return nil
	}

	builder := func(recordDTO dto.RecordDTO) *vo.RecordWeekOrMonthVO {
		recordTime := recordDTO.RecordTime
		return &vo.RecordWeekOrMonthVO{
			StartDate: recordTime.AddDate(0, 0, -recordTime.Day()+1).
				Format("2006-01-02"),
			Amount: recordDTO.Amount,
		}
	}
	return buildCommonWeekOrMonthVO(endDateObj, records, util.MonthByDate, builder)
}

func buildCommonWeekOrMonthVO(endDateObj time.Time,
	records *[]dto.RecordDTO,
	timeFun func(time.Time) int,
	builder func(dto.RecordDTO) *vo.RecordWeekOrMonthVO) *[]vo.RecordWeekOrMonthVO {

	var result []vo.RecordWeekOrMonthVO
	var temp *vo.RecordWeekOrMonthVO = nil
	var lastAdded *vo.RecordWeekOrMonthVO = nil

	var lastCachedIndex = timeFun(endDateObj)
	for _, recordDTO := range *records {
		curIndex := timeFun(recordDTO.RecordTime)
		if curIndex == lastCachedIndex {
			if temp == nil {
				temp = builder(recordDTO)
			} else {
				temp.Amount += recordDTO.Amount
			}
		} else {
			lastCachedIndex = curIndex
			if temp != nil {
				result = append(result, *temp)
				lastAdded = temp
				logger.Info(recordDTO.RecordTime.Format("2006-01-02"), curIndex)
			}
			temp = builder(recordDTO)
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

func CalculateAccountBalance() []*domain.Account {
	db := dal.GetDB()
	var list []domain.Record

	accountMap := ListAccountMap()
	db.Where("1=1").Find(&list)
	for _, account := range accountMap {
		account.CurrentAmount = account.InitAmount
	}

	for _, record := range list {
		account := accountMap[record.AccountId]
		if constant.IsExpense(record.Type) {
			account.CurrentAmount -= record.Amount
		} else {
			account.CurrentAmount += record.Amount
		}
	}

	for _, account := range accountMap {
		affected := db.Model(&account).
			Select("current_amount").
			Updates(map[string]interface{}{
				"current_amount": account.CurrentAmount,
			}).RowsAffected
		logger.Debug(account.Name, account.CurrentAmount, affected)
	}

	return ListAccounts()
}
