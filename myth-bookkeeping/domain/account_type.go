package domain

import "github.com/jinzhu/gorm"

type (
	AccountType struct {
		gorm.Model

		Name string
	}
)

func (AccountType) TableName() string {
	return "account_type"
}
