package dto

import (
	"simbirGo/internal/database/models"
	"simbirGo/internal/entities"
)

func RentEntitieToModel(rent entities.Rent, rentType uint) models.Rent {
	return models.Rent{
		Id:          rent.Id,
		TransportId: rent.TransportId,
		UserId:      rent.UserId,
		TimeStart:   rent.TimeStart,
		TimeEnd:     rent.TimeEnd,
		PriceOfUnit: rent.PriceOfUnit,
		PriceTypeId: rentType,
		FinalPrice:  rent.FinalPrice,
	}
}

func RentModelToEntitie(rent models.Rent, rentType string) entities.Rent {
	return entities.Rent{
		Id:          rent.Id,
		TransportId: rent.TransportId,
		UserId:      rent.UserId,
		TimeStart:   rent.TimeStart,
		TimeEnd:     rent.TimeEnd,
		PriceOfUnit: rent.PriceOfUnit,
		PriceType:   rentType,
		FinalPrice:  rent.FinalPrice,
	}
}
