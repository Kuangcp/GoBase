package common

import (
	"github.com/gin-gonic/gin"
	"github.com/kuangcp/gobase/pkg/ghelp"
	"mybook/app/common/constant"
	"mybook/app/domain"
	"mybook/app/service"
	"mybook/app/vo"
	"strconv"
)

// 简单查询

func HealthCheck(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func ListRecordType(c *gin.Context) {
	_, list := constant.GetRecordTypeMap()
	ghelp.GinSuccessWith(c, list)
}

func ListCategoryType(c *gin.Context) {
	_, list := constant.GetCategoryTypeMap()
	ghelp.GinSuccessWith(c, list)
}

func ListCategoryTree(c *gin.Context) {
	_, list := constant.GetCategoryTypeMap()

	var result []*vo.CategoryTree
	categories := service.ListCategories()
	for _, enum := range list {
		var temp []domain.Category
		for _, category := range categories {
			if category.TypeId == enum.Index {
				temp = append(temp, category)
			}
		}
		child := buildTreeRoot(temp)
		result = append(result, &vo.CategoryTree{
			ID:       uint(enum.Index),
			Name:     enum.Name,
			Children: child,
		})
	}

	ghelp.GinSuccessWith(c, result)
}

func buildTreeRoot(categories []domain.Category) []*vo.CategoryTree {
	var result []*vo.CategoryTree
	if len(categories) == 0 {
		return result
	}
	var exist = make(map[uint]string)
	for _, category := range categories {
		if category.ParentId == 0 {
			result = append(result, &vo.CategoryTree{
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

func removeHandled(data []domain.Category, exist map[uint]string) []domain.Category {
	var temp []domain.Category
	for _, category := range data {
		_, ok := exist[category.ID]
		if ok {
			continue
		}
		temp = append(temp, category)
	}
	return temp
}

func appendChild(tree []*vo.CategoryTree, node domain.Category) bool {
	for _, categoryTree := range tree {
		if categoryTree.ID == node.ParentId {
			categoryTree.Children = append(categoryTree.Children, &vo.CategoryTree{
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
		ghelp.GinSuccessWith(c, service.ListCategories())
		return
	}

	i, _ := strconv.Atoi(recordType)
	typeEnum := constant.GetCategoryTypeByRecordTypeIndex(int8(i))
	if typeEnum != nil {
		list := service.FindLeafCategoryByTypeId(typeEnum.Index)
		ghelp.GinSuccessWith(c, list)
	} else {
		ghelp.GinFailed(c)
	}
}
