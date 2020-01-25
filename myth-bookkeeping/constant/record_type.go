package constant

const (
	EXPENSE      int8 = 1
	INCOME       int8 = 2
	TRANSFER_IN  int8 = 3
	TRANSFER_OUT int8 = 4
	BORROW       int8 = 5
	REVERT       int8 = 6
)

func IsValidRecordType(typeValue int8) bool {
	if typeValue >= EXPENSE && typeValue <= REVERT {
		return true
	}
	return false
}
