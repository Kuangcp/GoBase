package domain

type (
	Record struct {
		Id         int32
		AccountId  int16
		CategoryId int16

		//Type 支出 收入 转出 转入
		Type int8

		CreateTime int64
		UpdateTime int64
		IsDeleted  int8
	}
)
