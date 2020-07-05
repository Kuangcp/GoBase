package domain

import "github.com/jinzhu/gorm"

type (
	Account struct {
		gorm.Model

		Name          string
		InitAmount    int
		CurrentAmount int
		TypeId        int8
		// 账本id
		BookId int

		// 信用卡 最大额度
		MaxAmount int

		// 信用卡 账单日
		BillDay int8

		// 信用卡 还款日 负数表示账单日的相对天数
		RepaymentDay int8
	}
)

func (Account) TableName() string {
	return "account"
}
