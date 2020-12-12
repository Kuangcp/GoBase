package report

import "time"

type (
	RecordQueryParam struct {
		StartDate  string `form:"startDate" json:"startDate"`
		EndDate    string `form:"endDate" json:"endDate"`
		TypeId     int    `form:"typeId" json:"typeId"`
		CategoryId int    `form:"categoryId" json:"categoryId"`
		ChartType  string `form:"chartType" json:"chartType"`
		ShowLabel  bool   `form:"showLabel" json:"showLabel"`
		Period     string `form:"period" json:"period"`

		startDate  time.Time
		endDate    time.Time
		timeFmt    string
		sqlTimeFmt string
	}
)
