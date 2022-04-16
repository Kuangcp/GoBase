package user

import (
	"github.com/gin-gonic/gin"
	"github.com/kuangcp/gobase/pkg/ghelp"
	"mybook/app/common/util"
)

func ListUser(c *gin.Context) {
	accounts := ListUsers()
	result := util.Copy(accounts, new([]UserDTO)).(*[]UserDTO)
	ghelp.GinSuccessWith(c, result)
}
