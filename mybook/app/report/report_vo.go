package report

import (
	"fmt"
)

type (
	LineChartVO struct {
		Legends []string `json:"legends"` // 类别
		Lines   []LineVO `json:"lines"`   // 类别对应的数据
		XAxis   []string `json:"xAxis"`
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
	CategorySumVO struct {
		CategoryId uint
		Sum        int
		Period     string
	}
)

func (this *CategorySumVO) BuildKey() string {
	return BuildKey(this.CategoryId, this.Period)
}
func BuildKey(categoryId uint, period string) string {
	return fmt.Sprint(categoryId) + ":" + period
}
