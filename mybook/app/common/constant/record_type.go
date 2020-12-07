package constant

const (
	// 支出
	RecordExpense int8 = 1
	// 收入
	RecordIncome int8 = 2
	// 转出
	RecordTransferOut int8 = 3
	// 转入
	RecordTransferIn int8 = 4

	// 收支 仅用于报表
	RecordOverview int8 = 9
)

var ERecordExpense = NewBaseEnum(RecordExpense, "支出")
var ERecordIncome = NewBaseEnum(RecordIncome, "收入")
var ERecordTransferOut = NewBaseEnum(RecordTransferOut, "转出")
var ERecordTransferIn = NewBaseEnum(RecordTransferIn, "转入")

var recordTypeMap map[int8]*BaseEnum
var recordTypeList []*BaseEnum

func GetRecordTypeMap() (map[int8]*BaseEnum, []*BaseEnum) {
	if recordTypeMap == nil {
		recordTypeMap, recordTypeList = MakeMap(ERecordExpense, ERecordIncome,
			ERecordTransferOut, ERecordTransferIn)
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
	return typeEnum != nil && (typeEnum == ERecordTransferOut || typeEnum == ERecordTransferIn)
}

func IsExpense(index int8) bool {
	return index == RecordExpense || index == RecordTransferOut
}
