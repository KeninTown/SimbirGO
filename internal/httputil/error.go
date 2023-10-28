package httpUtil

import "github.com/gin-gonic/gin"

type ResponseError struct {
	Error string `json:"err" example:"error occures"`
}

func NewResponseError(ctx *gin.Context, code int, err error) {
	responseErr := ResponseError{Error: err.Error()}
	ctx.AbortWithStatusJSON(code, responseErr)
}
