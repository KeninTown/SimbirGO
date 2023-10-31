package middlewares

import (
	httpUtil "simbirGo/internal/httputil"
	"simbirGo/internal/tokens"
	"strings"

	"github.com/gin-gonic/gin"
)

func CheckAuthification() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		authHeaderArray := strings.Split(authHeader, " ")
		if len(authHeaderArray) != 2 {
			httpUtil.NewResponseError(ctx, 401, "invalid authorization header")
			return
		}
		if authHeaderArray[1] == "" {
			httpUtil.NewResponseError(ctx, 401, "invalid jwt token")
		}
		token := authHeaderArray[1]
		if tokens.IsInBlackList(token) {
			httpUtil.NewResponseError(ctx, 401, "unauthorized")
			return
		}

		tokenData, err := tokens.ParseToken(token)
		if err != nil {
			httpUtil.NewResponseError(ctx, 401, err.Error())
			return
		}
		ctx.Set("id", tokenData.Id)
		ctx.Set("isAdmin", tokenData.IsAdmin)
		ctx.Next()
	}
}
