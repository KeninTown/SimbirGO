package usecase

import (
	"fmt"
	"simbirGo/internal/database/models"
	"simbirGo/internal/dto"
	"simbirGo/internal/entities"
	"simbirGo/internal/tokens"
)

type AuthRepository interface {
	FindUserByUsername(username string) models.User
	FindUserById(id uint) models.User
	CreateUser(user models.User) models.User
	SaveUser(user models.User)
	GetUsers(start uint, count int) []models.User
	DeleteUser(id uint)
}

type AuthUsecase struct {
	r AuthRepository
}

func New(r AuthRepository) AuthUsecase {
	return AuthUsecase{r: r}
}

func (au AuthUsecase) MyAccount(id uint) (entities.User, error) {
	user := au.r.FindUserById(id)
	if user.Id == 0 {
		return entities.User{}, fmt.Errorf("user is not exist")
	}
	return dto.UserModelToEntitie(user), nil
}

func (au AuthUsecase) SignIn(user entities.User) (string, error) {
	userModel := au.r.FindUserByUsername(user.Username)
	if userModel.Id == 0 {
		return "", fmt.Errorf("username is not exist")
	}

	if user.Password != userModel.Password {
		return "", fmt.Errorf("invalid password")
	}

	userEntite := dto.UserModelToEntitie(userModel)
	token, err := tokens.GenerateNewJwt(userEntite)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (au AuthUsecase) SignUp(user entities.User) (string, error) {
	candidate := au.r.FindUserByUsername(user.Username)
	if candidate.Id != 0 {
		return "", fmt.Errorf("user is already exist")
	}

	userModel := dto.UserEntitieToModels(user)
	userModel = au.r.CreateUser(userModel)
	userEntite := dto.UserModelToEntitie(userModel)
	token, err := tokens.GenerateNewJwt(userEntite)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (au AuthUsecase) Update(user entities.User) (entities.User, error) {
	userModel := au.r.FindUserById(user.Id)
	if userModel.Id == 0 {
		return entities.User{}, fmt.Errorf("user is not exist")
	}
	candidate := au.r.FindUserByUsername(user.Username)
	if candidate.Id != 0 && candidate.Id != userModel.Id {
		return entities.User{}, fmt.Errorf("username is taken")
	}
	userModel.Username = user.Username
	userModel.Password = user.Password
	au.r.SaveUser(userModel)

	return dto.UserModelToEntitie(userModel), nil
}

//adminAuth

func (au AuthUsecase) GetUsers(start, count uint) []entities.User {
	usersModels := au.r.GetUsers(start, int(count))
	usersEntities := make([]entities.User, 0, len(usersModels))
	for _, user := range usersModels {
		usersEntities = append(usersEntities, dto.UserModelToEntitie(user))
	}
	return usersEntities
}

func (au AuthUsecase) CreateUser(user entities.User) (entities.User, error) {
	candidate := au.r.FindUserByUsername(user.Username)
	if candidate.Id != 0 {
		return entities.User{}, fmt.Errorf("user is already exist")
	}

	userModel := dto.UserEntitieToModels(user)
	userModel = au.r.CreateUser(userModel)
	userEntite := dto.UserModelToEntitie(userModel)
	return userEntite, nil
}

func (au AuthUsecase) UpdateUser(user entities.User) (entities.User, error) {
	userModel := au.r.FindUserById(user.Id)
	if userModel.Id == 0 {
		return entities.User{}, fmt.Errorf("user is not exist")
	}
	candidate := au.r.FindUserByUsername(user.Username)
	if candidate.Id != 0 && candidate.Id != userModel.Id {
		return entities.User{}, fmt.Errorf("username is taken")
	}
	userModel.Username = user.Username
	userModel.Password = user.Password
	userModel.IsAdmin = user.IsAdmin
	userModel.Balance = user.Balance
	au.r.SaveUser(userModel)

	return dto.UserModelToEntitie(userModel), nil
}

func (au AuthUsecase) DeleteUser(id uint) {
	au.r.DeleteUser(id)
}
