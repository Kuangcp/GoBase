package report

import (
	"fmt"
	"mybook/app/common/util"
	"time"
)

type (
	RecordQueryParam struct {
		StartDate  string `form:"startDate" json:"startDate"`
		EndDate    string `form:"endDate" json:"endDate"`
		TypeId     int    `form:"typeId" json:"typeId"`
		CategoryId int    `form:"categoryId" json:"categoryId"`
		ChartType  string `form:"chartType" json:"chartType"`
		ShowLabel  bool   `form:"showLabel" json:"showLabel"`
		Period     string `form:"period" json:"period"`

		startDate    time.Time
		endDate      time.Time
		paramTimeFmt string                   // 参数开始结束时间的格式化
		periodFunc   func(t time.Time) string // 横坐标数据格式化
		sqlTimeFmt   string                   // SQL查询时间格式化
	}
)

func (param *RecordQueryParam) FillTimeFmt() {
	if param.Period == "" {
		return
	}

	param.periodFunc = func(t time.Time) string {
		return t.Format(param.paramTimeFmt)
	}

	switch param.Period {
	case yearPeriod:
		param.paramTimeFmt = "2006"
		param.sqlTimeFmt = "%Y"
	case monthPeriod:
		param.paramTimeFmt = "2006-01"
		param.sqlTimeFmt = "%Y-%m"
	case weekPeriod:
		param.paramTimeFmt = "2006-01-02"
		param.sqlTimeFmt = "%Y-%W"
		param.periodFunc = func(t time.Time) string {
			year, week := util.WeekOfYearByDate(t)
			weekStr := ""
			if week < 10 {
				weekStr = fmt.Sprintf("0%d", week)
			} else {
				weekStr = fmt.Sprintf("%d", week)
			}
			return fmt.Sprintf("%d-%s", year, weekStr)
		}
	case dayPeriod:
		param.paramTimeFmt = "2006-01-02"
		param.sqlTimeFmt = "%Y-%m-%d"
	default:
		param.paramTimeFmt = "2006-01"
		param.sqlTimeFmt = "%Y-%m"
	}
}
