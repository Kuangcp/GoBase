package constant

type (
	RecordTypeEnum struct {
		*BaseEnum
		Color string
	}
)

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
	ReportRecordOverview int = 9
	// 支出 父分类报表
	ReportExCategoryOverview int = 10
	// 收入 父分类报表
	ReportInCategoryOverview int = 11
)

var ERecordExpense = RecordTypeEnum{NewBaseEnum(RecordExpense, "支出"), "#D87A80"}
var ERecordIncome = RecordTypeEnum{NewBaseEnum(RecordIncome, "收入"), "#47ABF2"}
var ERecordTransferOut = RecordTypeEnum{NewBaseEnum(RecordTransferOut, "转出"), ""}
var ERecordTransferIn = RecordTypeEnum{NewBaseEnum(RecordTransferIn, "转入"), ""}

var recordTypeMap map[int8]Enum
var recordTypeList []Enum

func GetRecordTypeMap() (map[int8]Enum, []Enum) {
	if recordTypeMap == nil {
		recordTypeMap, recordTypeList = MakeMap(ERecordExpense, ERecordIncome,
			ERecordTransferOut, ERecordTransferIn)
	}
	return recordTypeMap, recordTypeList
}

func GetRecordTypeByIndex(index int8) Enum {
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
