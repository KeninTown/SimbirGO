package usecase

import (
	"fmt"
	"simbirGo/internal/database/models"
	"simbirGo/internal/dto"
	"simbirGo/internal/entities"
)

type AuthRepository interface {
	FindUser(username string) models.User
	CreateUser(user models.User) models.User
}

type AuthUsecase struct {
	r AuthRepository
}

func New(r AuthRepository) AuthUsecase {
	return AuthUsecase{r: r}
}

func (au AuthUsecase) SignUp(user entities.User) (string, error) {
	candidate := au.r.FindUser(user.Username)
	if candidate.Id != 0 {
		return "", fmt.Errorf("username is already exist")
	}
	userModel := dto.UserEntitieToModels(user)
	userModel = au.r.CreateUser(userModel)
	userEntite := dto.UserModelToEntitie(userModel)
	token, err := GenerateNewJwt(userEntite)
	if err != nil {
		return "", err
	}
	return token, nil
}
