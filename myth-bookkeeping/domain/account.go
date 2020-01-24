package domain

import "github.com/jinzhu/gorm"

type (
	Account struct {
		gorm.Model

		Name       string
		InitAmount int32
		TypeId     int8

		// 信用卡 最大额度
		MaxAmount int32

		// 信用卡 账单日
		BillDay int8

		// 信用卡 还款日
		RepaymentDay int8
	}
)

func (Account) TableName() string {
	return "account"
}
