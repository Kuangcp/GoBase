package domain

import (
	"github.com/jinzhu/gorm"
	"time"
)

type (
	Record struct {
		gorm.Model

		AccountId uint
		// 转账记录id 转账需要
		TransferId uint
		Amount     int
		Comment    string

		// 交易的分类id
		CategoryId uint

		//Type 支出 收入 转出 转入
		Type int8

		// 记录发生时刻
		RecordTime time.Time
	}
)

func (Record) TableName() string {
	return "record"
}
