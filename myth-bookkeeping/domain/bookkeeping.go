package domain

import "github.com/jinzhu/gorm"

type (
	BookKeeping struct {
		gorm.Model

		Name string
	}
)
