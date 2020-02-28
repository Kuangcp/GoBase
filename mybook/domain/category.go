package domain

import "github.com/jinzhu/gorm"

type (
	Category struct {
		gorm.Model
		// id:id 构成的绝对路径
		AbsoluteHierarchy string
		ParentId          uint
		// 叶节点才参与记账
		Leaf   bool
		Name   string
		// CategoryType
		TypeId int8
	}
)

func (Category) TableName() string {
	return "category"
}
