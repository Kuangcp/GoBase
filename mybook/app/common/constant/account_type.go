package constant

//TODO 社保 医保 住房公积金 网贷利息计算

const (
	AccountCash    int8 = 1
	AccountDeposit int8 = 2
	AccountCredit  int8 = 3
	AccountOnline  int8 = 4
	AccountFinance int8 = 5
	AccountAR      int8 = 6
	AccountAP      int8 = 7
)

const (
	AccountARId = 10001 // 应收id
	AccountAPId = 10002 // 应付id
)

var EAccountCash = BaseEnum{Index: AccountCash, Name: "现金"}
var EAccountDeposit = BaseEnum{Index: AccountDeposit, Name: "储蓄卡"}
var EAccountCredit = BaseEnum{Index: AccountCredit, Name: "信用卡"}
var EAccountOnline = BaseEnum{Index: AccountOnline, Name: "在线支付"}
var EAccountFinance = BaseEnum{Index: AccountFinance, Name: "理财"}
var EAccountAR = BaseEnum{Index: AccountAR, Name: "应收"}
var EAccountAP = BaseEnum{Index: AccountAP, Name: "应付"}

var accountTypeMap map[int8]Enum
var accountTypeList []Enum

func GetAccountTypeByIndex(index int8) Enum {
	if accountTypeMap == nil {
		accountTypeMap, accountTypeList = MakeMap(EAccountCash, EAccountDeposit, EAccountCredit, EAccountOnline, EAccountFinance, EAccountAP, EAccountAR)
	}
	return accountTypeMap[index]
}
