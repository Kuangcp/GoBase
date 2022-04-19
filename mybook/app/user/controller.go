package user

import (
	"mybook/app/common/dal"
	"mybook/app/common/util"

	"github.com/gin-gonic/gin"
	"github.com/kuangcp/gobase/pkg/ghelp"
)

func ListUser(c *gin.Context) {
	accounts := ListUsers()
	result := util.Copy(accounts, new([]UserDTO)).(*[]UserDTO)
	ghelp.GinSuccessWith(c, result)
}

func AddUser(c *gin.Context) {
	name := c.Query("name")
	if name == "" {
		ghelp.GinFailedWithMsg(c, "name 为空")
		return
	}

	users := ListUsers()
	for _, user := range users {
		if user.Name == name {
			ghelp.GinFailedWithMsg(c, "name "+name+" 已存在")
			return
		}
	}

	db := dal.GetDB()
	db.Save(&User{Name: name})
	ghelp.GinSuccess(c)
}
