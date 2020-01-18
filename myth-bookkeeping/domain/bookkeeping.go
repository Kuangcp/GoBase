package domain

type (
	BookKeeping struct {
		Id   int16
		Name string

		CreateTime int64
		UpdateTime int64
		IsDeleted  int8
	}
)
