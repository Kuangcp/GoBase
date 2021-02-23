package report

import (
	"fmt"
	"mybook/app/account"
	"mybook/app/category"
	"mybook/app/common/constant"
	"mybook/app/common/dal"
	"mybook/app/record"
	"sort"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kuangcp/gobase/pkg/ghelp"
)

var commonLabel = LabelVO{Show: false, Position: "insideRight"}

func buildBalanceReportParam(c *gin.Context) ghelp.ResultVO {
	var param RecordQueryParam
	err := c.ShouldBind(&param)
	if err != nil {
		return ghelp.FailedWithMsg("参数解析失败")
	}
	if param.StartDate == "" || param.EndDate == "" {
		return ghelp.FailedWithMsg("起止时间必填")
	}
	param.paramTimeFmt = "2006-01-02"
	startDate, err := time.Parse(param.paramTimeFmt, param.StartDate)
	if err != nil {
		return ghelp.FailedWithMsg(err.Error())
	}
	endDate, err := time.Parse(param.paramTimeFmt, param.EndDate)
	if err != nil {
		return ghelp.FailedWithMsg(err.Error())
	}
	if startDate.After(endDate) {
		return ghelp.FailedWithMsg("开始时间晚于结束时间")
	}

	param.startDate = startDate
	param.endDate = endDate
	if param.ChartType == "" {
		param.ChartType = "line"
	}

	return ghelp.SuccessWith(param)
}

func BalanceReport(c *gin.Context) {
	paramResult := buildBalanceReportParam(c)
	if paramResult.IsFailed() {
		ghelp.GinResultVO(c, paramResult)
		return
	}

	param := paramResult.Data.(RecordQueryParam)
	accountMap := account.ListAccountMap()
	curAmount := 0
	for _, a := range accountMap {
		curAmount += a.InitAmount
	}
	fmt.Println("init:", curAmount)

	records := record.QueryForBalance()
	fmt.Println(len(records))
	if len(records) == 0 {
		ghelp.GinFailed(c)
		return
	}

	// 按天生成数据
	var sameDays []record.RecordDTO
	lastTimeStr := ""
	var lastTime = param.startDate
	var data []float32
	var periodList []string

	for _, dto := range records {
		// 仅计算余额
		if dto.RecordTime.Unix() < param.startDate.Unix() {
			if dto.RecordType == constant.RecordExpense {
				curAmount -= dto.Amount
			} else {
				curAmount += dto.Amount
			}
			continue
		}

		curTime := dto.RecordTime.Format("2006-01-02")
		if len(sameDays) == 0 {
			sameDays = append(sameDays, dto)
			lastTimeStr = curTime
			lastTime = dto.RecordTime
			continue
		}

		if curTime == lastTimeStr {
			sameDays = append(sameDays, dto)
			continue
		}
		lastTimeStr = curTime

		// 以下为遇到不同天，清空缓存并计算

		// 填入完全没有记录的数据
		if lastTime.YearDay() != dto.RecordTime.YearDay() &&
			lastTime.Unix() > param.startDate.Unix() && lastTime.Unix() < param.endDate.Unix() {
			emptyTime := lastTime.AddDate(0, 0, 1)
			for emptyTime.Unix() < dto.RecordTime.Unix() && emptyTime.Unix() < param.endDate.Unix() {
				data = append(data, float32(curAmount)/100)
				fillTime := emptyTime.Format("2006-01-02")
				periodList = append(periodList, fillTime)
				emptyTime = emptyTime.AddDate(0, 0, 1)
			}
		}

		lastTime = dto.RecordTime
		curAmount += sameDayTotal(sameDays)

		if lastTime.Unix() > param.startDate.Unix() &&
			lastTime.Unix() < param.endDate.Unix() {
			data = append(data, float32(curAmount)/100)
			periodList = append(periodList, curTime)
		}
		sameDays = nil
	}

	ghelp.GinSuccessWith(c, LineChartVO{
		Lines: []LineVO{{
			Type:      param.ChartType,
			Name:      "余额",
			Data:      data,
			Stack:     "all",
			AreaStyle: "{normal: {}}",
			Label:     commonLabel,
			Color:     colorSet[12%len(colorSet)],
		}},
		XAxis:   periodList,
		Legends: []string{"余额"}},
	)
}

func sameDayTotal(sameDays []record.RecordDTO) int {
	dailyAmount := 0
	for _, recordDTO := range sameDays {
		if recordDTO.RecordType == constant.RecordExpense {
			dailyAmount -= recordDTO.Amount
		} else {
			dailyAmount += recordDTO.Amount
		}
	}
	return dailyAmount
}

func CategoryPeriodReport(c *gin.Context) {
	paramResult := buildParam(c)
	if paramResult.IsFailed() {
		ghelp.GinResultVO(c, paramResult)
		return
	}

	param := paramResult.Data.(RecordQueryParam)
	commonLabel.Show = param.ShowLabel

	periodList := buildPeriodList(param)
	finalStart := param.StartDate
	finalEnd := param.EndDate
	if param.Period == yearPeriod {
		finalStart += "-01"
		finalEnd += "-01"
	}

	var sumResult []CategorySumVO
	sumResult = buildQueryData(param, finalStart, finalEnd)
	if len(sumResult) == 0 {
		ghelp.GinFailedWithMsg(c, "数据为空")
		return
	}

	var legends []string
	var existCategoryMap = make(map[uint]int)
	var periodNumMap = make(map[string]float32)
	var lines []LineVO

	if param.TypeId == int(constant.RecordOverview) {
		legends = append(legends, constant.ERecordIncome.Name, constant.ERecordExpense.Name, "结余")
		for _, sum := range sumResult {
			periodNumMap[sum.BuildKey()] = sum.Sum
			_, ok := existCategoryMap[sum.CategoryId]
			if !ok {
				existCategoryMap[sum.CategoryId] = 0
			}
		}

		lines = buildLinesForOverview(periodList, periodNumMap, param)
	} else {
		categoryList := category.FindLeafCategoryByTypeId(int8(param.TypeId))
		var categoryNameMap = make(map[uint]string)
		for _, entity := range *categoryList {
			categoryNameMap[entity.ID] = entity.Name
		}

		for _, sum := range sumResult {
			periodNumMap[sum.BuildKey()] = sum.Sum
			_, ok := existCategoryMap[sum.CategoryId]
			if !ok {
				existCategoryMap[sum.CategoryId] = 0
				legends = append(legends, categoryNameMap[sum.CategoryId])
			}
		}

		lines = buildLines(existCategoryMap, periodList, periodNumMap, param, categoryNameMap)
	}

	ghelp.GinSuccessWith(c, LineChartVO{Lines: lines, XAxis: periodList, Legends: legends})
}

func buildLines(existCategoryMap map[uint]int,
	periodList []string,
	periodNumMap map[string]float32,
	param RecordQueryParam,
	categoryNameMap map[uint]string) []LineVO {

	var existCategoryList []uint
	for k := range existCategoryMap {
		existCategoryList = append(existCategoryList, k)
	}
	sort.Slice(existCategoryList, func(i, j int) bool {
		return existCategoryList[i] < existCategoryList[j]
	})

	var lines []LineVO
	for _, categoryId := range existCategoryList {
		var data []float32
		for _, period := range periodList {
			key := BuildKey(categoryId, period)
			_, exist := periodNumMap[key]
			if exist {
				data = append(data, periodNumMap[key])
			} else {
				data = append(data, 0)
			}
		}
		line := LineVO{
			Type:      param.ChartType,
			Name:      categoryNameMap[categoryId],
			Data:      data,
			Stack:     "all",
			AreaStyle: "{normal: {}}",
			Label:     commonLabel,
			Color:     colorSet[int(categoryId)%len(colorSet)],
		}
		lines = append(lines, line)
	}
	return lines
}

func buildLinesForOverview(periodList []string, periodNumMap map[string]float32, param RecordQueryParam) []LineVO {
	var lines []LineVO
	var balanceData []int32

	for _, typeId := range []constant.RecordTypeEnum{constant.ERecordIncome, constant.ERecordExpense} {
		categoryId := uint(typeId.Index)
		var data []float32
		for i, period := range periodList {
			key := BuildKey(categoryId, period)
			_, exist := periodNumMap[key]
			var temp float32 = 0.0
			if exist {
				temp = periodNumMap[key]
			}

			// 计算结余
			data = append(data, temp)
			if len(balanceData) <= i {
				balanceData = append(balanceData, 0)
			}
			if categoryId == uint(constant.RecordExpense) {
				balanceData[i] += -int32(temp * 100)
			} else {
				balanceData[i] += int32(temp * 100)
			}
		}

		line := LineVO{
			Type:      param.ChartType,
			Name:      typeId.Name,
			Data:      data,
			AreaStyle: "{normal: {}}",
			Label:     commonLabel,
			Color:     typeId.Color,
		}
		lines = append(lines, line)
	}

	var finalBalanceData []float32
	for _, datum := range balanceData {
		finalBalanceData = append(finalBalanceData, float32(datum)/100.0)
	}
	line := LineVO{
		Type:      param.ChartType,
		Name:      "结余",
		Data:      finalBalanceData,
		AreaStyle: "{normal: {}}",
		Label:     commonLabel,
		Color:     "#97B552",
	}
	lines = append(lines, line)
	return lines
}

func buildQueryData(param RecordQueryParam, finalStart string, finalEnd string) []CategorySumVO {
	var sumResult []CategorySumVO
	db := dal.GetDB()
	if param.TypeId == int(constant.RecordOverview) {
		db.Table("record").
			Select("type as category_id, sum(amount)/100.0 sum, strftime('"+param.sqlTimeFmt+"',record_time) as period").
			Where(" type in (?,?)", constant.RecordExpense, constant.RecordIncome).
			Where("record_time BETWEEN ? AND ?", finalStart, finalEnd).
			Group("type, period").Find(&sumResult)
	} else {
		where := db.Table("record").
			Select("category_id, sum(amount)/100.0 sum, strftime('" + param.sqlTimeFmt + "',record_time) as period")
		if param.TypeId != 0 {
			where = where.Where(" type = ?", param.TypeId)
		}
		if param.CategoryId != 0 {
			where = where.Where(" category_id = ?", param.CategoryId)
		}
		where.Where("record_time BETWEEN ? AND ?", finalStart, finalEnd).
			Group("category_id, period").Find(&sumResult)
	}
	return sumResult
}

func buildPeriodList(param RecordQueryParam) []string {
	start := param.startDate

	var result []string
	for !start.After(param.endDate) {
		result = append(result, param.periodFunc(start))
		switch param.Period {
		case yearPeriod:
			start = start.AddDate(1, 0, 0)
		case monthPeriod:
			start = start.AddDate(0, 1, 0)
		case weekPeriod:
			start = start.AddDate(0, 0, 7)
		case dayPeriod:
			start = start.AddDate(0, 0, 1)
		default:
			start = start.AddDate(0, 1, 0)
		}
	}
	return result[:len(result)-1]
}

func buildParam(c *gin.Context) ghelp.ResultVO {
	var param RecordQueryParam
	err := c.ShouldBind(&param)
	if err != nil {
		return ghelp.FailedWithMsg("参数解析失败")
	}
	if param.StartDate == "" || param.EndDate == "" || param.ChartType == "" || (param.TypeId == 0 && param.CategoryId == 0) {
		return ghelp.FailedWithMsg("参数含空值")
	}

	param.FillTimeFmt()

	startDate, err := time.Parse(param.paramTimeFmt, param.StartDate)
	if err != nil {
		return ghelp.FailedWithMsg(err.Error())
	}
	endDate, err := time.Parse(param.paramTimeFmt, param.EndDate)
	if err != nil {
		return ghelp.FailedWithMsg(err.Error())
	}
	if startDate.After(endDate) {
		return ghelp.FailedWithMsg("开始时间晚于结束时间")
	}
	param.startDate = startDate
	param.endDate = endDate
	return ghelp.SuccessWith(param)
}
