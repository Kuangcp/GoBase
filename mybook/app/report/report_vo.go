package report

import "strconv"

type CategorySumVO struct {
	CategoryId uint
	Sum        float32
	Period     string
}

func (this *CategorySumVO) BuildKey() string {
	return BuildKey(this.CategoryId, this.Period)
}
func BuildKey(categoryId uint, period string) string {
	return strconv.FormatUint(uint64(categoryId), 10) + ":" + period
}
