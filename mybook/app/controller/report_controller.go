package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/kuangcp/gobase/pkg/ghelp"
	"mybook/app/service"
	"time"
)

type (
	RecordQueryParam struct {
		StartDate string `form:"startDate" json:"startDate"`
		EndDate   string `form:"endDate" json:"endDate"`
		TypeId    int    `form:"typeId" json:"typeId"`
		ChartType string `form:"chartType" json:"chartType"`
		ShowLabel bool   `form:"showLabel" json:"showLabel"`
		startDate time.Time
		endDate   time.Time
	}

	LineChartVO struct {
		Lines   []LineVO `json:"lines"`
		XAxis   []string `json:"xAxis"`
		Legends []string `json:"legends"`
	}
	LineVO struct {
		Type      string  `json:"type"`
		Name      string  `json:"name"`
		Stack     string  `json:"stack"`
		Data      []int   `json:"data"`
		Color     string  `json:"color"`
		AreaStyle string  `json:"areaStyle"`
		Label     LabelVO `json:"label"`
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
var commonLabel = LabelVO{Show: false, Position: "insideRight"}

func CategoryMonthMap(c *gin.Context) {
	paramResult := buildParam(c)
	if paramResult.IsFailed() {
		ghelp.GinResultVO(c, paramResult)
		return
	}

	param := paramResult.Data.(RecordQueryParam)
	list := service.FindLeafCategoryByTypeId(int8(param.TypeId))
	var legends []string
	for _, category := range *list {
		legends = append(legends, category.Name)
	}
	months := buildMonth(param)
	var lines []LineVO
	commonLabel.Show = param.ShowLabel

	//TODO
	for i := 0; i < 10; i++ {
		var perMonth []int
		for j := 0; j < len(months); j++ {
			perMonth = append(perMonth, (i*3+j)%5)
		}
		lines = append(lines, LineVO{
			Type:      param.ChartType,
			Name:      legends[i],
			Data:      perMonth,
			Stack:     "all",
			AreaStyle: "{normal: {}}",
			Label:     commonLabel,
			Color:     colorSet[i%len(colorSet)],
		})
	}

	ghelp.GinSuccessWith(c, LineChartVO{Lines: lines, XAxis: months, Legends: legends})
}

func buildMonth(param RecordQueryParam) []string {
	start := param.startDate
	var result []string
	for !start.After(param.endDate) {
		result = append(result, start.Format("2006-01"))
		start = start.AddDate(0, 1, 0)
	}
	return result
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
	startDate, err := time.Parse("2006-01", param.StartDate)
	if err != nil {
		return ghelp.FailedWithMsg(err.Error())
	}
	endDate, err := time.Parse("2006-01", param.EndDate)
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
