package constant

type (
	AccountType struct {
		TypeId int8
		Name   string
	}
)

//TODO 社保 医保 住房公积金 网贷利息计算

const (
	CASH_TYPE    int8 = 1
	DEPOSIT_TYPE int8 = 2
	CREDIT_TYPE  int8 = 3
	ONLINE_TYPE  int8 = 4
	FINANCE_TYPE int8 = 5
)

var CASH = AccountType{TypeId: CASH_TYPE, Name: "现金"}
var DEPOSIT = AccountType{TypeId: DEPOSIT_TYPE, Name: "储蓄卡"}
var CREDIT = AccountType{TypeId: CREDIT_TYPE, Name: "信用卡"}
var ONLINE = AccountType{TypeId: ONLINE_TYPE, Name: "在线支付"}
var FINANCE = AccountType{TypeId: FINANCE_TYPE, Name: "理财"}

var accountTypeList = []AccountType{CASH, DEPOSIT, CREDIT, ONLINE, FINANCE}
var accountTypeMap map[int8]AccountType

func GetById(id int8) AccountType {
	if accountTypeMap == nil {
		accountTypeMap = make(map[int8]AccountType)
		for i := range accountTypeList {
			value := accountTypeList[i]
			accountTypeMap[value.TypeId] = value
		}
	}
	return accountTypeMap[id]
}
