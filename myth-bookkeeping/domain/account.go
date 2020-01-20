package domain

import "github.com/jinzhu/gorm"

type (
	Account struct {
		gorm.Model

		Name       string
		InitAmount int32
		TypeId     int8
	}
)

func (Account) TableName() string {
	return "account"
}
