package category

import (
	"fmt"

	"mybook/app/common/constant"
	"mybook/app/common/dal"
	"mybook/app/common/util"

	"github.com/kuangcp/gobase/pkg/cuibase"
	"github.com/kuangcp/logger"
)

func AddCategory(entity *Category) {
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
	if current == nil || current.ID == 0 {
		logger.Error("current not exist")
		return
	}

	db := dal.GetDB()
	current.ParentId = id
	// TODO update
	db.Model(&current).Update("parent_id", id)

	news := FindCategoryByName(name)
	logger.Info(current.ID, current.Name, current.ParentId, news.ID, news.ParentId)
}

func FindCategoryByName(name string) *Category {
	db := dal.GetDB()
	var result Category
	db.Where("name = ?", name).First(&result)
	return &result
}

func FindCategoryByTypeId(typeId int8) *[]Category {
	db := dal.GetDB()

	var lists []Category
	db.Where("type_id = ?", typeId).Find(&lists)
	return &lists
}

func FindLeafCategoryByTypeId(typeId int8) *[]Category {
	db := dal.GetDB()

	var lists []Category
	db.Where("type_id = ? AND leaf = true", typeId).Find(&lists)
	return &lists
}

func FindCategoryById(id uint) *Category {
	db := dal.GetDB()
	var result Category
	db.Where("id = ?", id).First(&result)
	return &result
}

func ListCategories() []Category {
	db := dal.GetDB()
	var lists []Category
	db.Where("1=1").Find(&lists)
	return lists
}

func ListCategoryMap() map[uint]Category {
	categories := ListCategories()
	result := make(map[uint]Category)
	for i := range categories {
		category := categories[i]
		result[category.ID] = category
	}
	return result
}

func PrintCategory() {
	db := dal.GetDB()
	var lists []Category
	db.Where("1=1").Find(&lists)

	resultMap := make(map[int8][]Category)
	for i := range lists {
		category := lists[i]

		categories := resultMap[category.TypeId]
		if categories == nil {
			resultMap[category.TypeId] = []Category{}
			categories = resultMap[category.TypeId]
		}

		resultMap[category.TypeId] = append(categories, category)
	}

	_, categories := constant.GetCategoryTypeMap()
	for i := range categories {
		enum := categories[i]
		value := ""
		for i := range resultMap[enum.GetIndex()] {
			category := resultMap[enum.GetIndex()][i]

			chFormat := util.BuildCHCharFormat(10, category.Name)
			value += fmt.Sprintf("%3d %s"+chFormat, category.ID, category.Name, "")
			if i%10 == 9 {
				value += "\n"
			}
		}
		if len(resultMap[enum.GetIndex()]) != 0 {
			fmt.Printf(cuibase.Cyan.Print("> %v  ")+"\n%v\n\n", enum.GetName(), value)
		}
	}
}
