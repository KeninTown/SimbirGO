package handlers

import (
	"fmt"
	"simbirGo/internal/entities"

	"github.com/gin-gonic/gin"
)

type AuthUsecase interface {
	SignUp(user entities.User) (string, error)
}

type AuthHandlers struct {
	uc AuthUsecase
}

func New(uc AuthUsecase) AuthHandlers {
	return AuthHandlers{uc: uc}
}

func (ah AuthHandlers) MyAccount(ctx *gin.Context) {

}

func (ah AuthHandlers) SignIn(ctx *gin.Context) {

}

func (ah AuthHandlers) SignUp(ctx *gin.Context) {
	var user entities.User
	if err := ctx.BindJSON(&user); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	fmt.Println(user)
	token, err := ah.uc.SignUp(user)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(201, gin.H{"token": token})
}
