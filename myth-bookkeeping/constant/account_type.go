package constant

type (
	AccountType struct {
		TypeId int8
		Name   string
	}
)

//TODO 社保 医保 住房公积金

var CASH_TYPE int8  = 1
var DEPOSIT_TYPE int8 = 2
var CREDIT_TYPE int8 = 3
var ONLINE_TYPE int8 = 4
var FINANCE_TYPE int8 = 5

var CASH = AccountType{TypeId: CASH_TYPE, Name: "现金"}
var DEPOSIT = AccountType{TypeId: DEPOSIT_TYPE, Name: "储蓄卡"}
var CREDIT = AccountType{TypeId: CREDIT_TYPE, Name: "信用卡"}
var ONLINE = AccountType{TypeId: ONLINE_TYPE, Name: "在线支付"}
var FINANCE = AccountType{TypeId: FINANCE_TYPE, Name: "理财"}

var accountTypeMap map[int8]AccountType

func GetById(id int8) AccountType {
	if accountTypeMap == nil {
		accountTypeMap = make(map[int8]AccountType)
		accountTypeMap[CASH_TYPE] = CASH
		accountTypeMap[DEPOSIT_TYPE] = DEPOSIT
		accountTypeMap[CREDIT_TYPE] = CREDIT
		accountTypeMap[ONLINE_TYPE] = ONLINE
		accountTypeMap[FINANCE_TYPE] = FINANCE
	}
	return accountTypeMap[id]
}
