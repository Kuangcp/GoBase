package service

import (
	"github.com/kuangcp/gobase/myth-bookkeeping/dal"
	"github.com/kuangcp/gobase/myth-bookkeeping/domain"
	"log"
)

func AddCategory(entity *domain.Category) {
	db := dal.GetDB()
	db.Create(entity)
}

func SetParentId(name string, id uint) {
	db := dal.GetDB()
	var lists []domain.Category
	db.Select("name=" + name).Find(&lists)
	log.Println(lists)

}

func FindCategoryById(id uint) *domain.Category {
	db := dal.GetDB()
	var result domain.Category
	db.Where("id = ?", id).First(&result)
	return &result
}
