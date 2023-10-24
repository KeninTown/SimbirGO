package dto

import (
	"simbirGo/internal/database/models"
	"simbirGo/internal/entities"
)

func UserEntitieToModels(user entities.User) models.User {
	return models.User{
		Id:       user.Id,
		Username: user.Username,
		Password: user.Password,
		IsAdmin:  user.IsAdmin,
	}
}

func UserModelToEntitie(user models.User) entities.User {
	return entities.User{
		Id:       user.Id,
		Username: user.Username,
		Password: user.Password,
		IsAdmin:  user.IsAdmin,
	}
}
