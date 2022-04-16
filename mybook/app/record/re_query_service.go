package record

import (
	"github.com/kuangcp/gobase/pkg/stopwatch"
	"github.com/kuangcp/logger"
	"mybook/app/account"
	"mybook/app/category"
	"mybook/app/common/constant"
	"mybook/app/common/dal"
	"mybook/app/common/util"
	"sort"
	"strconv"
	"time"
	"unsafe"
)

func QueryForBalance() []RecordDTO {
	db := dal.GetDB()
	var lists []RecordEntity
	query := db.Where("deleted_at IS NULL AND type IN (?,?)", constant.RecordExpense, constant.RecordIncome)
	query.Order("record_time", true).Find(&lists)
	if len(lists) < 1 {
		return nil
	}
	return entityToDTO(lists)
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

	result := entityToDTO(lists)
	return result
}

// 帐目按类型分组 typeId record_type
func queryCategoryRecord(param QueryRecordParam) *MonthCategoryRecordResult {
	var startDate = param.StartDate
	var endDate = param.EndDate
	var typeId = param.TypeId
	db := dal.GetDB()
	var result []*MonthCategoryRecordDTO
	query := db.Table("record").
		Select("record.category_id, category.name, sum(amount) as amount,type").
		Joins("left join category on record.category_id = category.id")
	if len(typeId) != 0 {
		query = query.Where("record.type =?", typeId)
	}

	query = query.Where("record_time between ? and ?", startDate, endDate).
		Where("record.deleted_at is null").
		Group("category_id").
		Scan(&result)
	if len(result) == 0 {
		return &MonthCategoryRecordResult{}
	}

	var totalAmount = 0
	for i := range result {
		recordDTO := result[i]

		recordDTO.RecordTypeName = constant.GetRecordTypeByIndex(recordDTO.Type).GetName()
		totalAmount += recordDTO.Amount
	}

	// 逆序
	sort.Slice(result, func(i, j int) bool {
		return result[i].Amount >= result[j].Amount
	})
	for _, dto := range result {
		dto.AmountPercent = float32(dto.Amount*100) / float32(totalAmount)
		dto.AmountPercent = float32(int(dto.AmountPercent*100)) / 100
	}

	return &MonthCategoryRecordResult{
		List:        result,
		TotalAmount: totalAmount,
	}
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

	accountMap := account.ListAccountMap()
	for _, accountEntity := range accountMap {
		accountEntity.CurrentAmount = accountEntity.InitAmount
	}

	var list []RecordEntity
	watch := stopwatch.NewWithName("calculate balance")
	watch.Start("query list")
	db.Where("1=1").Find(&list)
	for _, record := range list {
		accountEntity := accountMap[record.AccountId]
		if constant.IsExpense(record.Type) {
			accountEntity.CurrentAmount -= record.Amount
		} else {
			accountEntity.CurrentAmount += record.Amount
		}
	}
	watch.Stop()

	for _, accountEntity := range accountMap {
		watch.Start("update " + accountEntity.Name)
		db.Model(&accountEntity).
			Select("current_amount").
			Updates(map[string]interface{}{
				"current_amount": accountEntity.CurrentAmount,
			})
		//logger.Release(account.Name, account.CurrentAmount, affected)
		watch.Stop()
	}
	logger.Info(watch.PrettyPrint())

	return account.ListAccounts()
}

func entityToDTO(lists []RecordEntity) []RecordDTO {
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
