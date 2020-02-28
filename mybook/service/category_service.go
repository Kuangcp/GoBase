package service

import (
	"fmt"
	"github.com/kuangcp/gobase/mybook/constant"
	"github.com/kuangcp/gobase/mybook/dal"
	"github.com/kuangcp/gobase/mybook/domain"
	"log"
	"strconv"
)

func AddCategory(entity *domain.Category) {
	db := dal.GetDB()
	db.Create(entity)
}

func SetParentId(name string, id uint) {
	parent := FindCategoryById(id)
	if parent == nil {
		log.Println("parent id not exist")
		return
	}

	current := FindCategoryByName(name)
	if current == nil {
		log.Println("current not exist")
		return
	}

	db := dal.GetDB()
	current.ParentId = id
	// TODO update
	db.Update(current)
}

func FindCategoryByName(name string) *domain.Category {
	db := dal.GetDB()
	var result domain.Category
	db.Where("name = ?", name).First(&result)
	return &result
}

func FindCategoryById(id uint) *domain.Category {
	db := dal.GetDB()
	var result domain.Category
	db.Where("id = ?", id).First(&result)
	return &result
}

func PrintCategory(_ []string) {
	db := dal.GetDB()
	var lists []domain.Category
	db.Where("1=1").Find(&lists)

	resultMap := make(map[int8][]domain.Category)
	for i := range lists {
		category := lists[i]

		categories := resultMap[category.TypeId]
		if categories == nil {
			resultMap[category.TypeId] = []domain.Category{}
			categories = resultMap[category.TypeId]
		}

		resultMap[category.TypeId] = append(categories, category)
	}

	_, categories := constant.GetCategoryTypeMap()
	for i := range categories {
		enum := categories[i]
		value := ""
		for i := range resultMap[enum.Index] {
			category := resultMap[enum.Index][i]
			nameLen := len(category.Name)

			format := "%2d %s%" + strconv.Itoa(12-nameLen/3*2) + "s"
			value += fmt.Sprintf(format, category.ID, category.Name, "")
			if i%10 == 9 {
				value += "\n"
			}
		}
		if len(resultMap[enum.Index]) != 0 {
			fmt.Printf("\n------------------- %v ------------------- \n%v\n", enum.Name, value)
		}
	}
}
