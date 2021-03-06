package account

type (
	AccountDTO struct {
		ID            uint
		Name          string
		InitAmount    int32
		CurrentAmount int32
		TypeId        int8
		// 账本id
		BookId int

		// 信用卡 最大额度
		MaxAmount int32

		// 信用卡 账单日
		BillDay int8

		// 信用卡 还款日 负数表示账单日的相对天数
		RepaymentDay int8
	}
)
