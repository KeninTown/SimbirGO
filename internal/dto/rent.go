package dto

import (
	"simbirGo/internal/database/models"
	"simbirGo/internal/entities"
)

func RentEntitieToModel(rent entities.Rent) models.Rent {
	return models.Rent{
		Id:          rent.Id,
		TransportId: rent.TransportId,
		UserId:      rent.UserId,
		TimeStart:   rent.TimeStart,
		TimeEnd:     rent.TimeEnd,
		PriceOfUnit: rent.PriceOfUnit,
		PriceType:   rent.PriceType,
		FinalPrice:  rent.FinalPrice,
	}
}

func RentModelToEntitie(rent models.Rent) entities.Rent {
	return entities.Rent{
		Id:          rent.Id,
		TransportId: rent.TransportId,
		UserId:      rent.UserId,
		TimeStart:   rent.TimeStart,
		TimeEnd:     rent.TimeEnd,
		PriceOfUnit: rent.PriceOfUnit,
		PriceType:   rent.PriceType,
		FinalPrice:  rent.FinalPrice,
	}
}
