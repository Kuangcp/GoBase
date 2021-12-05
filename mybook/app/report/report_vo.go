package report

import (
	"fmt"
)

type (
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
	CategorySumVO struct {
		CategoryId uint
		Sum        float32
		Period     string
	}
)

func (this *CategorySumVO) BuildKey() string {
	return BuildKey(this.CategoryId, this.Period)
}
func BuildKey(categoryId uint, period string) string {
	return  fmt.Sprint(categoryId) + ":" + period
}
