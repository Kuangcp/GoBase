package dto

import "time"

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
)
