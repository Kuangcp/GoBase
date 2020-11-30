package param

type (
	RecordCreateParamVO struct {
		AccountId       int    `json:"accountId"`
		TargetAccountId int    `json:"targetAccountId"`
		Amount          int     `json:"amount"`
		CategoryId      int    `json:"categoryId"`
		TypeId          int8    `json:"typeId"` // TypeId 含义为 categoryTypeId
		Date            []string `json:"date"`
		Comment         string  `json:"comment"`
	}
)
