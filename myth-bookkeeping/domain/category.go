package domain

type (
	Category struct {
		// id:id 构成的绝对路径
		AbsoluteHierarchy string
		ParentId          int16
		Name              string
	}
)

func (Category) TableName() string {
	return "category"
}