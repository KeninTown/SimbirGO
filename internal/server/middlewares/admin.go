package middlewares

import (
	httpUtil "simbirGo/internal/httputil"

	"github.com/gin-gonic/gin"
)

func CheckAdminStatus() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		isAdmin := ctx.GetBool("isAdmin")
		if !isAdmin {
			httpUtil.NewResponseError(ctx, 403, "admin access only")
			return
		}
		ctx.Next()
	}
}
