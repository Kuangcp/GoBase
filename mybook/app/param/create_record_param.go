package param

type (
	CreateRecordParam struct {
		AccountId       string `json:"accountId"`
		TargetAccountId string `json:"targetAccountId"`
		Amount          string `json:"amount"`
		CategoryId      string `json:"categoryId"`
		TypeId          string `json:"typeId"` // TypeId 含义为 categoryTypeId
		Date            string `json:"date"`
		Comment         string `json:"comment"`
	}
)
