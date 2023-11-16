package main

import (
	"github.com/kuangcp/gobase/pkg/ctool"
	"testing"
	"time"
)

func TestWeekDay(t *testing.T) {
	start := time.Now().AddDate(0, 0, -20)
	w := int(start.Weekday())

	now := time.Now()
	for i := 0; i < 7; i++ {
		off := -(w+6)%7 + 7*i
		tmp := start.AddDate(0, 0, off)
		if tmp.After(now) {
			return
		}

		println(tmp.Format(ctool.YYYY_MM_DD))

	}
}
