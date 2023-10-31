package httpUtil

import "github.com/gin-gonic/gin"

type ResponseError struct {
	Error string `json:"err" example:"error occures"`
}

func NewResponseError(ctx *gin.Context, code int, msg string) {
	responseErr := ResponseError{Error: msg}
	ctx.AbortWithStatusJSON(code, responseErr)
}
