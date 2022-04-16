package user

import "github.com/jinzhu/gorm"

type (
	// User 借贷发生方
	User struct {
		gorm.Model

		Name string
	}
)

func (User) TableName() string {
	return "user"
}
