package record

type (
	RecordCreateParamVO struct {
		AccountId       int      `json:"accountId"`
		TargetAccountId int      `json:"targetAccountId"`
		Amount          string   `json:"amount"`
		CategoryId      int      `json:"categoryId"`
		TypeId          int8     `json:"typeId"` // TypeId 含义为 categoryTypeId
		Date            []string `json:"date"`
		Comment         string   `json:"comment"`
	}
	QueryRecordParam struct {
		AccountId  string `form:"accountId"`
		CategoryId string `form:"categoryId"`
		TypeId     string `form:"typeId"` // record_type
		StartDate  string `form:"startDate"`
		EndDate    string `form:"endDate"`
	}
)
