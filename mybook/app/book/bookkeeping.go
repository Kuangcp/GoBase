package book

import "github.com/jinzhu/gorm"

type (
	BookKeeping struct {
		gorm.Model

		Name    string
		Comment string
	}
)

func (BookKeeping) TableName() string {
	return "book_keeping"
}
