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
	FindUserById(id uint) models.User
	FindTranspots(start, count int, transportId uint) []models.Transport
	DeleteTransport(id uint)
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

	transportModel.OwnerId = transport.OwnerId
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

func (tu TransportUsecase) GetTransports(start, count int, transportType string) ([]entities.Transport, error) {
	transportTypeId := tu.r.FindTypeByName(transportType)
	if transportTypeId == 0 {
		return nil, fmt.Errorf("invalid transport type: %s", transportType)
	}
	transportModels := tu.r.FindTranspots(start, count, transportTypeId)
	transportEntites := make([]entities.Transport, 0, len(transportModels))
	for _, tr := range transportModels {
		transportEntites = append(transportEntites, dto.TransportModelToEntite(tr, transportType))
	}

	return transportEntites, nil
}

func (tu TransportUsecase) AdminCreateTransport(transport entities.Transport) (entities.Transport, error) {
	//find user with ownerId
	owner := tu.r.FindUserById(transport.OwnerId)
	if owner.Id == 0 {
		return entities.Transport{}, fmt.Errorf("user with id = %d is not exist", transport.OwnerId)
	}
	//create
	return tu.CreateTransport(transport)
}

func (tu TransportUsecase) AdminUpdateTransport(transport entities.Transport) (entities.Transport, error) {
	transportModel := tu.r.FindTranspot(transport.Id)
	if transportModel.Id == 0 {
		return entities.Transport{}, fmt.Errorf("transport is not exist")
	}

	owner := tu.r.FindUserById(transport.OwnerId)
	if owner.Id == 0 {
		return entities.Transport{}, fmt.Errorf("user with id = %d is not exist", transport.OwnerId)
	}

	typeId := tu.r.FindTypeByName(transport.TransportType)
	if typeId == 0 {
		return entities.Transport{}, fmt.Errorf("invalid transport type")
	}

	transportModel.OwnerId = transport.OwnerId
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

func (tu TransportUsecase) AdminDeleteTransport(id uint) error {
	transport := tu.r.FindTranspot(id)
	if transport.Id == 0 {
		return fmt.Errorf("transport is not exist")
	}
	tu.r.DeleteTransport(id)
	return nil
}
