package util

import (
	"fmt"
	"testing"
	"time"
)

func TestWeekOfYearByDate(t *testing.T) {
	date := time.Now().AddDate(0, 0, 1)
	year, week := WeekOfYearByDate(date)
	fmt.Println(year, week)
}
