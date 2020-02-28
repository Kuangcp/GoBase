package domain

import (
	"github.com/jinzhu/gorm"
	"time"
)

type (
	TransferRecord struct {
		gorm.Model

		FromAccount int16
		ToAccount   int16
		Amount int
		// 记录发生时刻
		RecordTime time.Time
	}
)
