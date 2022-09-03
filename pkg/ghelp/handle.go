package ghelp

import (
	"github.com/gin-gonic/gin"
)

func WrapVoidFunc(action func()) func(context *gin.Context) {
	return func(context *gin.Context) {
		action()
		GinSuccess(context)
	}
}
