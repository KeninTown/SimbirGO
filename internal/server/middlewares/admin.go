package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func CheckAdminStatus() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		isAdmin := ctx.GetBool("isAdmin")
		if !isAdmin {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": "access admin only"})
			return
		}
		ctx.Next()
	}
}
