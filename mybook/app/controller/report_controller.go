package controller

import (
	"mybook/app/common/dal"
	"mybook/app/service"
	"mybook/app/vo"
	"sort"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kuangcp/gobase/pkg/ghelp"
)

type (
	RecordQueryParam struct {
		StartDate string `form:"startDate" json:"startDate"`
		EndDate   string `form:"endDate" json:"endDate"`
		TypeId    int    `form:"typeId" json:"typeId"`
		ChartType string `form:"chartType" json:"chartType"`
		ShowLabel bool   `form:"showLabel" json:"showLabel"`
		Period    string `form:"period" json:"period"`

		startDate  time.Time
		endDate    time.Time
		timeFmt    string
		sqlTimeFmt string
	}

	LineChartVO struct {
		Lines   []LineVO `json:"lines"`
		XAxis   []string `json:"xAxis"`
		Legends []string `json:"legends"`
	}
	LineVO struct {
		Type      string    `json:"type"`
		Name      string    `json:"name"`
		Stack     string    `json:"stack"`
		Data      []float32 `json:"data"`
		Color     string    `json:"color"`
		AreaStyle string    `json:"areaStyle"`
		Label     LabelVO   `json:"label"`
	}
	LabelVO struct {
		Show     bool   `json:"show"`
		Position string `json:"position"`
	}
)

var colorSet = [...]string{
	"#c23531",
	"#2f4554",
	"#61a0a8",
	"#d48265",
	"#91c7ae",
	"#749f83",
	"#ca8622",
	"#bda29a",
	"#6e7074",
	"#546570",
	"#c4ccd3",
}

const (
	yearPeriod  = "year"
	monthPeriod = "month"
	dayPeriod   = "day"
)

var commonLabel = LabelVO{Show: false, Position: "insideRight"}

func gerTimeFmt(period string) (string, string) {
	switch period {
	case yearPeriod:
		return "2006", "%Y"
	case monthPeriod:
		return "2006-01", "%Y-%m"
	case dayPeriod:
		return "2006-01-02", "%Y-%m-%d"
	}
	return "2006-01", "%Y-%m"
}

func CategoryMonthMap(c *gin.Context) {
	paramResult := buildParam(c)
	if paramResult.IsFailed() {
		ghelp.GinResultVO(c, paramResult)
		return
	}

	param := paramResult.Data.(RecordQueryParam)
	commonLabel.Show = param.ShowLabel

	categoryList := service.FindLeafCategoryByTypeId(int8(param.TypeId))
	var categoryNameMap = make(map[uint]string)
	for _, category := range *categoryList {
		categoryNameMap[category.ID] = category.Name
	}

	periodList := buildPeriodList(param)

	var sumResult []vo.CategorySumVO
	db := dal.GetDB()
	db.Table("record").
		Select("category_id, sum(amount)/100.0 sum, strftime('"+param.sqlTimeFmt+"',record_time) as period").
		Where(" type = ?", param.TypeId).
		Where("record_time BETWEEN ? AND ?", param.StartDate, param.EndDate).
		Group("category_id, period").Find(&sumResult)

	var legends []string
	var existCategoryMap = make(map[uint]int)
	for _, sum := range sumResult {
		_, ok := existCategoryMap[sum.CategoryId]
		if !ok {
			existCategoryMap[sum.CategoryId] = 0
			legends = append(legends, categoryNameMap[sum.CategoryId])
		}
	}
	var existCategoryList []uint
	for k, _ := range existCategoryMap {
		existCategoryList = append(existCategoryList, k)
	}
	sort.Slice(existCategoryList, func(i, j int) bool {
		return existCategoryList[i] < existCategoryList[j]
	})
	var lines []LineVO
	for _, categoryId := range existCategoryList {
		var data []float32
		for _, period := range periodList {
			find := false
			for _, sum := range sumResult {
				if sum.CategoryId == categoryId {
					if sum.Period == period {
						data = append(data, sum.Sum)
						find = true
					}
				}
			}
			if !find {
				data = append(data, 0)
			}
		}
		lines = append(lines, LineVO{
			Type:      param.ChartType,
			Name:      categoryNameMap[categoryId],
			Data:      data,
			Stack:     "all",
			AreaStyle: "{normal: {}}",
			Label:     commonLabel,
			Color:     colorSet[int(categoryId)%len(colorSet)],
		})
	}

	ghelp.GinSuccessWith(c, LineChartVO{Lines: lines, XAxis: periodList, Legends: legends})
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
	param.timeFmt, param.sqlTimeFmt = gerTimeFmt(param.Period)

	startDate, err := time.Parse(param.timeFmt, param.StartDate)
	if err != nil {
		return ghelp.FailedWithMsg(err.Error())
	}
	endDate, err := time.Parse(param.timeFmt, param.EndDate)
	if err != nil {
		return ghelp.FailedWithMsg(err.Error())
	}
	if startDate.After(endDate) {
		return ghelp.FailedWithMsg("开始时间早于结束时间")
	}
	param.startDate = startDate
	param.endDate = endDate
	return ghelp.SuccessWith(param)
}
