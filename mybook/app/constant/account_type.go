package constant

//TODO 社保 医保 住房公积金 网贷利息计算

const (
	ACCOUNT_CASH    int8 = 1
	ACCOUNT_DEPOSIT int8 = 2
	ACCOUNT_CREDIT  int8 = 3
	ACCOUNT_ONLINE  int8 = 4
	ACCOUNT_FINANCE int8 = 5
)

var E_ACCOUNT_CASH = &BaseEnum{Index: ACCOUNT_CASH, Name: "现金"}
var E_ACCOUNT_DEPOSIT = &BaseEnum{Index: ACCOUNT_DEPOSIT, Name: "储蓄卡"}
var E_ACCOUNT_CREDIT = &BaseEnum{Index: ACCOUNT_CREDIT, Name: "信用卡"}
var E_ACCOUNT_ONLINE = &BaseEnum{Index: ACCOUNT_ONLINE, Name: "在线支付"}
var E_ACCOUNT_FINANCE = &BaseEnum{Index: ACCOUNT_FINANCE, Name: "理财"}

var accountTypeMap map[int8]*BaseEnum
var accountTypeList []*BaseEnum

func GetAccountTypeByIndex(index int8) *BaseEnum {
	if accountTypeMap == nil {
		accountTypeMap, accountTypeList = MakeMap(E_ACCOUNT_CASH, E_ACCOUNT_DEPOSIT, E_ACCOUNT_CREDIT, E_ACCOUNT_ONLINE, E_ACCOUNT_FINANCE)
	}
	return accountTypeMap[index]
}
