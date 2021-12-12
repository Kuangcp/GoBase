package record

import (
	"mybook/app/common/util"
)

type (
	recordVO struct {
		ID             uint
		AccountName    string
		CategoryName   string
		RecordType     int8
		RecordTypeName string
		Amount         int
		Comment        string
		RecordTime     string
	}

	recordWeekOrMonthVO struct {
		StartDate string
		EndDate   string
		Amount    int
	}
)

func convertToVO(from RecordDTO) *recordVO {
	var record recordVO
	util.Copy(from, &record)
	record.RecordTime = from.RecordTime.Format("2006-01-02")
	return &record
}

func convertToVOList(from []RecordDTO) []*recordVO {
	var result []*recordVO
	for _, dto := range from {
		temp := convertToVO(dto)
		if temp != nil {
			result = append(result, temp)
		}
	}
	return result
}
