package constant

const (
	// 支出
	RECORD_EXPENSE int8 = 1
	// 收入
	RECORD_INCOME int8 = 2
	// 转出
	RECORD_TRANSFER_OUT int8 = 3
	// 转入
	RECORD_TRANSFER_IN int8 = 4
)

var E_RECORD_EXPENSE  = NewBaseEnum(RECORD_EXPENSE ,"支出")
var E_RECORD_INCOME  = NewBaseEnum(RECORD_INCOME ,"收入")
var E_RECORD_TRANSFER_OUT = NewBaseEnum(RECORD_TRANSFER_OUT,"转出")
var E_RECORD_TRANSFER_IN = NewBaseEnum(RECORD_TRANSFER_IN,"转入")

var recordTypeMap map[int8]*BaseEnum

func GetRecordTypeByIndex(index int8) *BaseEnum {
	if recordTypeMap == nil {
		recordTypeMap = MakeMap(E_RECORD_EXPENSE, E_RECORD_INCOME, E_RECORD_TRANSFER_OUT, E_RECORD_TRANSFER_IN)
	}
	return recordTypeMap[index]
}

func IsValidRecordType(typeValue int8) bool {
	return GetRecordTypeByIndex(typeValue) != nil
}
