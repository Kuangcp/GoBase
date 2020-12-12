package constant

const (
	// 支出
	CategoryExpense int8 = 1

	// 收入
	CategoryIncome int8 = 2

	// 转账
	CategoryTransfer int8 = 3
)

var ECategoryExpense = NewBaseEnum(CategoryExpense, "支出")
var ECategoryIncome = NewBaseEnum(CategoryIncome, "收入")
var ECategoryTransfer = NewBaseEnum(CategoryTransfer, "转账")

var categoryTypeMap map[int8]Enum
var categoryTypeList []Enum

func GetCategoryTypeMap() (map[int8]Enum, []Enum) {
	if categoryTypeMap == nil {
		categoryTypeMap, categoryTypeList = MakeMap(ECategoryExpense, ECategoryIncome, ECategoryTransfer)
	}
	return categoryTypeMap, categoryTypeList
}

func GetCategoryTypeByIndex(index int8) Enum {
	maps, _ := GetCategoryTypeMap()
	return maps[index]
}

func IsValidCategoryType(typeValue int8) bool {
	return GetCategoryTypeByIndex(typeValue) != nil
}

func GetCategoryTypeByRecordTypeIndex(index int8) *BaseEnum {
	recordType := GetRecordTypeByIndex(index)
	if recordType == nil {
		return nil
	}

	switch recordType {
	case ERecordExpense:
		return ECategoryExpense
	case ERecordIncome:
		return ECategoryIncome
	case ERecordTransferIn:
		return ECategoryTransfer
	case ERecordTransferOut:
		return ECategoryTransfer
	}
	return nil
}
