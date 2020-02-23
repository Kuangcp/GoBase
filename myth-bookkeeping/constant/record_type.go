package constant

const (
	// 支出
	EXPENSE      int8 = 1
	// 收入
	INCOME       int8 = 2
	TRANSFER_IN  int8 = 3
	TRANSFER_OUT int8 = 4
	// 借出
	BORROW       int8 = 5
	// 归还
	REVERT       int8 = 6
)

func IsValidRecordType(typeValue int8) bool {
	if typeValue >= EXPENSE && typeValue <= REVERT {
		return true
	}
	return false
}
