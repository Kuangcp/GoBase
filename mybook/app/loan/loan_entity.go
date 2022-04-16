package loan

import (
	"github.com/jinzhu/gorm"
	"time"
)

type (
	// Entity 借贷记录
	Entity struct {
		gorm.Model

		// 流入流出账户
		AccountId uint

		// 转账记录时间戳 关联 记录表数据
		TransferId uint

		// 对方
		UserId uint

		// 借入 1 贷出 2
		LoanType int8

		// 记录发生时刻
		RecordTime time.Time
		// 期望还款时刻
		ExceptTime time.Time

		// 备注
		Comment string

		// 单位:分
		Amount int
	}
)

func (Entity) TableName() string {
	return "record_loan"
}
