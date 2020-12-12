package record

import (
	"mybook/app/account"
	"mybook/app/category"
	"mybook/app/common/constant"
	"mybook/app/common/dal"
	"mybook/app/common/util"
	"sort"
	"strconv"
	"time"
	"unsafe"

	"github.com/kuangcp/gobase/pkg/ghelp"

	"github.com/wonderivan/logger"
)

func addRecord(record *RecordEntity) {
	db := dal.GetDB()
	// TODO support multiple book
	record.BookId = 1
	db.Create(record)
}

func checkParam(record *RecordEntity) (ghelp.ResultVO, *category.Category, *account.Account) {
	category := category.FindCategoryById(record.CategoryId)
	if category == nil || !category.Leaf {
		return ghelp.FailedWithMsg("分类id无效"), nil, nil
	}

	accountEntity := account.FindAccountById(record.AccountId)
	if accountEntity == nil {
		return ghelp.FailedWithMsg("账户无效"), category, nil
	}

	if record.Amount <= 0 {
		return ghelp.FailedWithMsg("金额无效"), category, accountEntity
	}
	if !constant.IsValidRecordType(record.Type) {
		return ghelp.FailedWithMsg("类别无效"), category, accountEntity
	}
	return ghelp.Success(), category, accountEntity
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
		record := &RecordEntity{
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

func findRecord(param QueryRecordParam) []RecordDTO {
	db := dal.GetDB()
	var lists []RecordEntity
	accountId, _ := strconv.Atoi(param.AccountId)
	typeId, _ := strconv.Atoi(param.TypeId)
	categoryId, _ := strconv.Atoi(param.CategoryId)

	query := db.Where("deleted_at is null and record_time BETWEEN ? AND ?", param.StartDate, param.EndDate)
	entity := RecordEntity{}
	if typeId != 0 {
		entity.Type = int8(typeId)
	}
	if accountId != 0 {
		entity.AccountId = uint(accountId)
	}
	if categoryId != 0 {
		entity.CategoryId = uint(categoryId)
	}

	if unsafe.Sizeof(entity) != 0 {
		query = query.Where(&entity)
	}
	query.Order("record_time DESC", true).Find(&lists)
	if len(lists) < 1 {
		return nil
	}

	accountMap := account.ListAccountMap()
	categoryMap := category.ListCategoryMap()
	var result []RecordDTO
	for i := range lists {
		record := lists[i]
		ele := RecordDTO{
			ID:             record.ID,
			RecordType:     record.Type,
			AccountName:    accountMap[record.AccountId].Name,
			CategoryName:   categoryMap[record.CategoryId].Name,
			RecordTypeName: constant.GetRecordTypeByIndex(record.Type).GetName(),
			RecordTime:     record.RecordTime,
			Amount:         record.Amount,
			Comment:        record.Comment,
		}

		result = append(result, ele)
	}
	return result
}

// 帐目按类型分组 typeId record_type
func queryCategoryRecord(param QueryRecordParam) *[]MonthCategoryRecordDTO {
	var startDate = param.StartDate
	var endDate = param.EndDate
	var typeId = param.TypeId
	db := dal.GetDB()
	var result []MonthCategoryRecordDTO
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

	var temp []MonthCategoryRecordDTO
	for i := range result {
		recordDTO := &result[i]

		recordDTO.RecordTypeName = constant.GetRecordTypeByIndex(recordDTO.Type).GetName()

		temp = append(temp, *recordDTO)
	}

	// 逆序
	sort.Slice(temp, func(i, j int) bool {
		return temp[i].Amount >= temp[j].Amount
	})

	return &temp
}

func weekCategoryRecord(param QueryRecordParam) *[]recordWeekOrMonthVO {
	records := findRecord(param)
	if records == nil {
		return nil
	}
	endDateObj, err := time.Parse("2006-01-02", param.EndDate)
	if err != nil {
		return nil
	}

	builder := func(recordDTO RecordDTO) *recordWeekOrMonthVO {
		recordTime := recordDTO.RecordTime
		return &recordWeekOrMonthVO{
			StartDate: recordTime.AddDate(0, 0, -int(recordTime.Weekday()-time.Sunday)).
				Format("2006-01-02"),
			EndDate: recordTime.AddDate(0, 0, int(time.Saturday-recordTime.Weekday())).
				Format("2006-01-02"),
			Amount: recordDTO.Amount,
		}
	}
	return buildCommonWeekOrMonthVO(endDateObj, records, util.WeekByDate, builder)
}

func monthCategoryRecord(param QueryRecordParam) *[]recordWeekOrMonthVO {
	records := findRecord(param)
	if records == nil {
		return nil
	}
	endDateObj, err := time.Parse("2006-01-02", param.EndDate)
	if err != nil {
		return nil
	}

	builder := func(recordDTO RecordDTO) *recordWeekOrMonthVO {
		recordTime := recordDTO.RecordTime
		return &recordWeekOrMonthVO{
			StartDate: recordTime.AddDate(0, 0, -recordTime.Day()+1).
				Format("2006-01-02"),
			Amount: recordDTO.Amount,
		}
	}
	return buildCommonWeekOrMonthVO(endDateObj, records, util.MonthByDate, builder)
}

func buildCommonWeekOrMonthVO(endDateObj time.Time,
	records []RecordDTO,
	timeFun func(time.Time) int,
	builder func(RecordDTO) *recordWeekOrMonthVO) *[]recordWeekOrMonthVO {

	var result []recordWeekOrMonthVO
	var temp *recordWeekOrMonthVO = nil
	var lastAdded *recordWeekOrMonthVO = nil

	var lastCachedIndex = timeFun(endDateObj)
	for _, recordDTO := range records {
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

func calculateAndQueryAccountBalance() []*account.Account {
	db := dal.GetDB()
	var list []RecordEntity

	accountMap := account.ListAccountMap()
	db.Where("1=1").Find(&list)
	for _, accountEntity := range accountMap {
		accountEntity.CurrentAmount = accountEntity.InitAmount
	}

	for _, record := range list {
		accountEntity := accountMap[record.AccountId]
		if constant.IsExpense(record.Type) {
			accountEntity.CurrentAmount -= record.Amount
		} else {
			accountEntity.CurrentAmount += record.Amount
		}
	}

	for _, accountEntity := range accountMap {
		db.Model(&accountEntity).
			Select("current_amount").
			Updates(map[string]interface{}{
				"current_amount": accountEntity.CurrentAmount,
			})
		//logger.Debug(account.Name, account.CurrentAmount, affected)
	}

	return account.ListAccounts()
}
