package service

import (
	"github.com/kuangcp/gobase/pkg/ghelp"
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

	"github.com/wonderivan/logger"
)

func addRecord(record *domain.Record) {
	db := dal.GetDB()
	// TODO support multiple book
	record.BookId = 1
	db.Create(record)
}

func checkParam(record *domain.Record) (ghelp.ResultVO, *domain.Category, *domain.Account) {
	category := FindCategoryById(record.CategoryId)
	if category == nil || !category.Leaf {
		return ghelp.FailedWithMsg("分类id无效"), nil, nil
	}

	account := FindAccountById(record.AccountId)
	if account == nil {
		return ghelp.FailedWithMsg("账户无效"), category, nil
	}

	if record.Amount <= 0 {
		return ghelp.FailedWithMsg("金额无效"), category, account
	}
	if !constant.IsValidRecordType(record.Type) {
		return ghelp.FailedWithMsg("类别无效"), category, account
	}
	return ghelp.Success(), category, account
}

func CreateRecord(record *domain.Record) ghelp.ResultVO {
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

func createTransRecord(origin *domain.Record, target *domain.Record) ghelp.ResultVO {
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

func BuildRecordByField(param param.RecordCreateParamVO) ghelp.ResultVO {
	if len(param.Date) == 0 {
		return ghelp.FailedWithMsg("日期为空")
	}
	var recordList []*domain.Record
	for _, date := range param.Date {
		recordDate, e := time.Parse("2006-01-02", date)
		if e != nil {
			logger.Error(e)
			return ghelp.FailedWithMsg("date 参数错误")
		}
		record := &domain.Record{
			AccountId:  uint(param.AccountId),
			CategoryId: uint(param.CategoryId),
			Type:       param.TypeId,
			Amount:     param.Amount,
			RecordTime: recordDate,
		}
		if param.Comment != "" {
			record.Comment = param.Comment
		}
		recordList = append(recordList, record)
	}

	return ghelp.SuccessWith(recordList)
}

func CreateMultipleTypeRecord(param param.RecordCreateParamVO) ghelp.ResultVO {
	result := BuildRecordByField(param)
	if result.IsFailed() {
		return result
	}

	list := result.Data.([]*domain.Record)
	var successList []*domain.Record
	var failResults []*domain.Record

	for _, record := range list {
		if param.TargetAccountId != 0 && constant.IsTransferRecordType(record.Type) {
			record.Type = constant.RecordTransferOut

			now := time.Now()
			record.TransferId = uint(now.UnixNano())

			target := util.Copy(record, new(domain.Record)).(*domain.Record)
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
			resultVO := CreateRecord(record)
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

func FindRecord(param param.QueryRecordParam) *[]dto.RecordDTO {
	db := dal.GetDB()
	var lists []domain.Record
	accountId, _ := strconv.Atoi(param.AccountId)
	typeId, _ := strconv.Atoi(param.TypeId)
	categoryId, _ := strconv.Atoi(param.CategoryId)

	query := db.Where("deleted_at is null and record_time BETWEEN ? AND ?", param.StartDate, param.EndDate)
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
		db.Model(&account).
			Select("current_amount").
			Updates(map[string]interface{}{
				"current_amount": account.CurrentAmount,
			})
		//logger.Debug(account.Name, account.CurrentAmount, affected)
	}

	return ListAccounts()
}
