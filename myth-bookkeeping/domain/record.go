package domain

import "github.com/jinzhu/gorm"

type (
	Record struct {
		gorm.Model

		AccountId  int16
		CategoryId int16

		//Type 支出 收入 转出 转入 record_type
		Type int8
		// 转入或转出的对方账户 支出收入时没有值
		TargetAccount int16
	}
)

func (Record) TableName() string {
	return "record"
}