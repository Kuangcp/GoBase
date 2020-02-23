package domain

import (
	"github.com/jinzhu/gorm"
	"time"
)

type (
	Record struct {
		gorm.Model

		AccountId uint

		// 转账记录时间戳 联系转入和转出
		TransferId uint

		// 分为单位
		Amount     int

		// 备注
		Comment    string

		// 交易的分类id
		CategoryId uint

		//Type record_type 支出 收入 转出 转入
		Type int8

		// 记录发生时刻
		RecordTime time.Time
	}
)

func (Record) TableName() string {
	return "record"
}
