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

func GetCategoryTypeMap() map[int8]*BaseEnum {
	if categoryTypeMap == nil {
		categoryTypeMap = MakeMap(E_CATEGORY_EXPENSE, E_CATEGORY_INCOME, E_CATEGORY_TRANSFER)
	}
	return categoryTypeMap
}

func GetCategoryTypeByIndex(index int8) *BaseEnum {
	return GetCategoryTypeMap()[index]
}

func IsValidCategoryType(typeValue int8) bool {
	return GetCategoryTypeByIndex(typeValue) != nil
}
