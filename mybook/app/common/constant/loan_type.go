package constant

const (
	LoanBorrow int8 = 1 // 借入
	LoanLend   int8 = 2 // 贷出

	LoanBorrowRe int8 = 3 // 贷出-负债 借入的归还
	LoanLendRe   int8 = 4 // 借入-负债 贷出的归还
)
