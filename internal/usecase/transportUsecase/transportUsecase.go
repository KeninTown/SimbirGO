package transportusecase

import (
	"fmt"
	"simbirGo/internal/database/models"
	"simbirGo/internal/dto"
	"simbirGo/internal/entities"
)

type TransportRepository interface {
	FindTypeById(id uint) string
	FindTypeByName(typeName string) uint
	FindTranspot(id uint) models.Transport
	CreateTransport(transport models.Transport) models.Transport
	FindUserTransport(userId, transportId uint) models.Transport
	SaveTransport(transport models.Transport)
	DeleteUserTransport(ownerId, transportId uint)
}

type TransportUsecase struct {
	r TransportRepository
}

func New(r TransportRepository) TransportUsecase {
	return TransportUsecase{r: r}
}

func (tu TransportUsecase) GetTransport(id uint) (entities.Transport, error) {
	transportModel := tu.r.FindTranspot(id)
	if transportModel.Id == 0 {
		return entities.Transport{}, fmt.Errorf("transport is not exist")
	}

	typeStr := tu.r.FindTypeById(transportModel.TypeId)
	transportEntite := dto.TransportModelToEntite(transportModel, typeStr)
	return transportEntite, nil
}

func (tu TransportUsecase) CreateTransport(transport entities.Transport) (entities.Transport, error) {
	typeId := tu.r.FindTypeByName(transport.TransportType)
	if typeId == 0 {
		return entities.Transport{}, fmt.Errorf("invalid transport type")
	}
	transportModel := dto.TransporEntitieToModel(transport, typeId)
	transportModel = tu.r.CreateTransport(transportModel)
	transportEntite := dto.TransportModelToEntite(transportModel, transport.TransportType)
	return transportEntite, nil
}

func (tu TransportUsecase) UpdateUserTransport(transport entities.Transport) (entities.Transport, error) {
	transportModel := tu.r.FindUserTransport(transport.OwnerId, transport.Id)
	if transportModel.Id == 0 {
		return entities.Transport{}, fmt.Errorf("transport is not exist")
	}

	typeId := tu.r.FindTypeByName(transport.TransportType)
	if typeId == 0 {
		return entities.Transport{}, fmt.Errorf("invalid transport type")
	}

	transportModel.TypeId = typeId
	transportModel.CanBeRented = transport.CanBeRented
	transportModel.Model = transport.Model
	transportModel.Color = transport.Color
	transportModel.Identifier = transport.Identifier
	transportModel.Description = transport.Description
	transportModel.Latitude = transport.Latitude
	transportModel.Longitude = transport.Longitude
	transportModel.MinutePrice = transport.MinutePrice
	transportModel.DayPrice = transport.DayPrice

	tu.r.SaveTransport(transportModel)
	transportEntite := dto.TransportModelToEntite(transportModel, transport.TransportType)
	return transportEntite, nil
}

func (tu TransportUsecase) DeleteUserTransport(userId, transportId uint) error {
	transport := tu.r.FindUserTransport(userId, transportId)
	if transport.Id == 0 {
		return fmt.Errorf("transport is not exist")
	}
	tu.r.DeleteUserTransport(userId, transportId)
	return nil
}
