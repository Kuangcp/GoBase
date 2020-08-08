package util

import (
	"github.com/kuangcp/gobase/mybook/app/dto"
	"github.com/wonderivan/logger"
	"testing"
)

func TestName(t *testing.T) {

	var list = []dto.MonthCategoryRecordDTO{
		{Amount: 10}, {Amount: 2}, {Amount: 15}, {Amount: 6},
	}

	var temp []interface{}
	for _, recordDTO := range list {
		temp = append(temp, recordDTO)
	}

	Sort(SortWrapper{Data: temp,
		CompareLessFunc: func(a interface{}, b interface{}) bool {
			result := a.(dto.MonthCategoryRecordDTO).Amount < b.(dto.MonthCategoryRecordDTO).Amount
			logger.Info(a, b, result)
			return result
		}, Reverse: true})

	logger.Debug("temp", temp)
	logger.Debug("list", list)
}
