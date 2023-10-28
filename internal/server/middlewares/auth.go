package middlewares

import (
	"fmt"
	"net/http"
	httpUtil "simbirGo/internal/httputil"
	"simbirGo/internal/tokens"
	"strings"

	"github.com/gin-gonic/gin"
)

func CheckAuthification() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		authHeaderArray := strings.Split(authHeader, " ")
		fmt.Println("header = ", authHeader)
		if len(authHeaderArray) != 2 {
			httpUtil.NewResponseError(ctx, 401, fmt.Errorf("invalid authorization header"))
			return
		}
		if authHeaderArray[1] == "" {
			httpUtil.NewResponseError(ctx, 401, fmt.Errorf("invalid jwt token"))
		}
		token := authHeaderArray[1]
		if tokens.IsInBlackList(token) {
			httpUtil.NewResponseError(ctx, 401, fmt.Errorf("Unauthorized"))
			return
		}

		tokenData, err := tokens.ParseToken(token)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"err": err.Error()})
			return
		}
		fmt.Println("pivo")
		ctx.Set("id", tokenData.Id)
		ctx.Set("isAdmin", tokenData.IsAdmin)
		ctx.Next()
	}
}
