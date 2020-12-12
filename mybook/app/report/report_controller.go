package report

import (
	"mybook/app/category"
	"mybook/app/common/constant"
	"mybook/app/common/dal"
	"sort"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kuangcp/gobase/pkg/ghelp"
)

var commonLabel = LabelVO{Show: false, Position: "insideRight"}

// go时间格式，sqlite时间格式
func getTimeFmt(period string) (string, string) {
	switch period {
	case yearPeriod:
		return "2006", "%Y"
	case monthPeriod:
		return "2006-01", "%Y-%m"
	//case weekPeriod:
	//	return "", "%Y-%W"
	case dayPeriod:
		return "2006-01-02", "%Y-%m-%d"
	}
	return "2006-01", "%Y-%m"
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
			Select("category_id, sum(amount)/100.0 sum, strftime('"+param.sqlTimeFmt+"',record_time) as period").
			Where(" type = ?", param.TypeId)
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
		result = append(result, start.Format(param.timeFmt))
		switch param.Period {
		case yearPeriod:
			start = start.AddDate(1, 0, 0)
		case monthPeriod:
			start = start.AddDate(0, 1, 0)
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
	if param.StartDate == "" || param.EndDate == "" || param.ChartType == "" || param.TypeId == 0 {
		return ghelp.FailedWithMsg("参数含空值")
	}
	param.timeFmt, param.sqlTimeFmt = getTimeFmt(param.Period)

	startDate, err := time.Parse(param.timeFmt, param.StartDate)
	if err != nil {
		return ghelp.FailedWithMsg(err.Error())
	}
	endDate, err := time.Parse(param.timeFmt, param.EndDate)
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
