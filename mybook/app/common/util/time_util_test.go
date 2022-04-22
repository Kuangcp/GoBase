package util

import (
	"fmt"
	"testing"
	"time"

	"github.com/kuangcp/logger"
)

func TestWeekOfYearByDate(t *testing.T) {
	date := time.Now().AddDate(0, 0, 1)
	year, week := WeekOfYearByDate(date)
	fmt.Println(year, week)
}

func TestDayFromWeek(t *testing.T) {
	year := StartDayByWeekAndYear("2022", "23")
	logger.Info(year)
}
