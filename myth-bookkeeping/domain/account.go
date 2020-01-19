package domain

type (
	Account struct {
		Id         int16
		Name       string
		InitAmount int32
		Type       int8

		CreateTime int64
		UpdateTime int64
		IsDeleted  int8
	}
)
