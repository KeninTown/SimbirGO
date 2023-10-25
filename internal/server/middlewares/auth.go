package middlewares

import (
	"fmt"
	"net/http"
	"simbirGo/internal/tokens"

	"github.com/gin-gonic/gin"
)

func CheckAuthification() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token, err := ctx.Cookie("access_token")
		if err != nil || token == "" {
			ctx.AbortWithStatusJSON(401, gin.H{"err": "access_token cookie must be provided"})
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
