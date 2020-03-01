package service

import (
	"fmt"
	"github.com/kuangcp/gobase/cuibase"
	"github.com/kuangcp/gobase/mybook/constant"
	"github.com/kuangcp/gobase/mybook/dal"
	"github.com/kuangcp/gobase/mybook/domain"
	"github.com/kuangcp/gobase/mybook/util"
	"github.com/wonderivan/logger"
)

func AddCategory(entity *domain.Category) {
	db := dal.GetDB()
	db.Create(entity)
}

func SetParentId(name string, id uint) {
	parent := FindCategoryById(id)
	if parent == nil {
		logger.Error("parent id not exist")
		return
	}

	current := FindCategoryByName(name)
	if current == nil {
		logger.Error("current not exist")
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

func FindCategoryByTypeId(typeId int8) *[]domain.Category {
	db := dal.GetDB()

	var lists []domain.Category
	db.Where("type_id=?", typeId).Find(&lists)
	return &lists
}

func FindCategoryById(id uint) *domain.Category {
	db := dal.GetDB()
	var result domain.Category
	db.Where("id = ?", id).First(&result)
	return &result
}

func FindAllCategory() *[]domain.Category {
	db := dal.GetDB()
	var lists []domain.Category
	db.Where("1=1").Find(&lists)
	return &lists
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

			chFormat := util.BuildCHCharFormat(10, category.Name)
			value += fmt.Sprintf("%3d %s"+chFormat, category.ID, category.Name, "")
			if i%10 == 9 {
				value += "\n"
			}
		}
		if len(resultMap[enum.Index]) != 0 {
			fmt.Printf(cuibase.Cyan+"> %v  "+cuibase.End+"\n%v\n\n", enum.Name, value)
		}
	}
}
