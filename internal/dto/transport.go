package dto

import (
	"simbirGo/internal/database/models"
	"simbirGo/internal/entities"
)

func TransporEntitieToModel(transport entities.Transport, typeId uint) models.Transport {
	return models.Transport{
		Id:          transport.Id,
		OwnerId:     transport.OwnerId,
		TypeId:      typeId,
		CanBeRented: transport.CanBeRented,
		Model:       transport.Model,
		Color:       transport.Color,
		Identifier:  transport.Identifier,
		Description: transport.Description,
		Latitude:    transport.Latitude,
		Longitude:   transport.Longitude,
		MinutePrice: transport.MinutePrice,
		DayPrice:    transport.DayPrice,
	}
}

func TransportModelToEntite(transport models.Transport, typeStr string) entities.Transport {
	return entities.Transport{
		Id:            transport.Id,
		OwnerId:       transport.OwnerId,
		TransportType: typeStr,
		CanBeRented:   transport.CanBeRented,
		Model:         transport.Model,
		Color:         transport.Color,
		Identifier:    transport.Identifier,
		Description:   transport.Description,
		Latitude:      transport.Latitude,
		Longitude:     transport.Longitude,
		MinutePrice:   transport.MinutePrice,
		DayPrice:      transport.DayPrice,
	}
}
