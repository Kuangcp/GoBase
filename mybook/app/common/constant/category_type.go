package constant

const (
	// 支出
	CATEGORY_EXPENSE int8 = 1

	// 收入
	CATEGORY_INCOME int8 = 2

	// 转账
	CATEGORY_TRANSFER int8 = 3
)

var E_CATEGORY_EXPENSE = NewBaseEnum(CATEGORY_EXPENSE, "支出")
var E_CATEGORY_INCOME = NewBaseEnum(CATEGORY_INCOME, "收入")
var E_CATEGORY_TRANSFER = NewBaseEnum(CATEGORY_TRANSFER, "转账")

var categoryTypeMap map[int8]*BaseEnum
var categoryTypeList []*BaseEnum

func GetCategoryTypeMap() (map[int8]*BaseEnum, []*BaseEnum) {
	if categoryTypeMap == nil {
		categoryTypeMap, categoryTypeList = MakeMap(E_CATEGORY_EXPENSE, E_CATEGORY_INCOME, E_CATEGORY_TRANSFER)
	}
	return categoryTypeMap, categoryTypeList
}

func GetCategoryTypeByIndex(index int8) *BaseEnum {
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
		return E_CATEGORY_EXPENSE
	case ERecordIncome:
		return E_CATEGORY_INCOME
	case ERecordTransferIn:
		return E_CATEGORY_TRANSFER
	case ERecordTransferOut:
		return E_CATEGORY_TRANSFER
	}
	return nil
}
