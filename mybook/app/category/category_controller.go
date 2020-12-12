package category

import (
	"mybook/app/common/constant"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kuangcp/gobase/pkg/ghelp"
)

func ListCategoryTree(c *gin.Context) {
	_, list := constant.GetCategoryTypeMap()

	var result []*CategoryTree
	categories := ListCategories()
	for _, enum := range list {
		var temp []Category
		for _, entity := range categories {
			if entity.TypeId == enum.GetIndex() {
				temp = append(temp, entity)
			}
		}
		child := buildTreeRoot(temp)
		result = append(result, &CategoryTree{
			ID:       uint(enum.GetIndex()),
			Name:     enum.GetName(),
			Children: child,
		})
	}

	ghelp.GinSuccessWith(c, result)
}

func buildTreeRoot(categories []Category) []*CategoryTree {
	var result []*CategoryTree
	if len(categories) == 0 {
		return result
	}
	var exist = make(map[uint]string)
	for _, category := range categories {
		if category.ParentId == 0 {
			result = append(result, &CategoryTree{
				ID:   category.ID,
				Name: category.Name,
			})
			exist[category.ID] = ""
		}
	}
	categories = removeHandled(categories, exist)
	for len(categories) > 0 {
		exist = make(map[uint]string)
		for _, category := range categories {
			appendResult := appendChild(result, category)
			if appendResult {
				exist[category.ID] = ""
			}
			//logger.Info(category.Name, appendResult)
		}
		categories = removeHandled(categories, exist)
	}
	return result
}

func removeHandled(data []Category, exist map[uint]string) []Category {
	var temp []Category
	for _, category := range data {
		_, ok := exist[category.ID]
		if ok {
			continue
		}
		temp = append(temp, category)
	}
	return temp
}

func appendChild(tree []*CategoryTree, node Category) bool {
	for _, categoryTree := range tree {
		if categoryTree.ID == node.ParentId {
			categoryTree.Children = append(categoryTree.Children, &CategoryTree{
				ID:   node.ID,
				Name: node.Name,
			})
			return true
		}
		if len(categoryTree.Children) != 0 {
			appendResult := appendChild(categoryTree.Children, node)
			if appendResult {
				return true
			}
		}
	}
	return false
}

func ListCategory(c *gin.Context) {
	recordType := c.Query("recordType")
	if recordType == "" {
		ghelp.GinSuccessWith(c, ListCategories())
		return
	}

	i, _ := strconv.Atoi(recordType)
	typeEnum := constant.GetCategoryTypeByRecordTypeIndex(int8(i))
	if typeEnum != nil {
		list := FindLeafCategoryByTypeId(typeEnum.Index)
		ghelp.GinSuccessWith(c, list)
	} else {
		ghelp.GinFailed(c)
	}
}
