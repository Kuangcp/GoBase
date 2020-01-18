package domain

type (
	Category struct {
		// id:id 构成的绝对路径
		AbsoluteHierarchy string
		Id                int16
		ParentId          int16
		Name              string

		CreateTime int64
		UpdateTime int64
		IsDeleted  int8
	}
)
