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

var E_RECORD_EXPENSE = NewBaseEnum(RECORD_EXPENSE, "支出")
var E_RECORD_INCOME = NewBaseEnum(RECORD_INCOME, "收入")
var E_RECORD_TRANSFER_OUT = NewBaseEnum(RECORD_TRANSFER_OUT, "转出")
var E_RECORD_TRANSFER_IN = NewBaseEnum(RECORD_TRANSFER_IN, "转入")

var recordTypeMap map[int8]*BaseEnum
var recordTypeList []*BaseEnum

func GetRecordTypeMap() (map[int8]*BaseEnum, []*BaseEnum) {
	if recordTypeMap == nil {
		recordTypeMap, recordTypeList = MakeMap(E_RECORD_EXPENSE, E_RECORD_INCOME, E_RECORD_TRANSFER_OUT, E_RECORD_TRANSFER_IN)
	}
	return recordTypeMap, recordTypeList
}

func GetRecordTypeByIndex(index int8) *BaseEnum {
	maps, _ := GetRecordTypeMap()
	return maps[index]
}

func IsValidRecordType(typeValue int8) bool {
	return GetRecordTypeByIndex(typeValue) != nil
}

func IsTransferRecordType(typeValue int8) bool {
	typeEnum := GetRecordTypeByIndex(typeValue)
	return typeEnum != nil && (typeEnum == E_RECORD_TRANSFER_OUT || typeEnum == E_RECORD_TRANSFER_IN)
}
