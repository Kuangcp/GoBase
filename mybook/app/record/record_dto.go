package record

import (
	"time"
)

type (
	RecordDTO struct {
		ID             uint
		AccountName    string
		CategoryName   string
		RecordType     int8
		RecordTypeName string
		Amount         int
		Comment        string
		RecordTime     time.Time
	}

	MonthCategoryRecordResult struct {
		List        []MonthCategoryRecordDTO
		TotalAmount int
	}
	MonthCategoryRecordDTO struct {
		CategoryId     uint
		Name           string
		Type           int8
		RecordTypeName string
		Amount         int
		Date           string
		AmountPercent  float32 // 金额占比
	}
)
