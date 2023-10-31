package rentUsecase

import (
	"fmt"
	"math"
	"simbirGo/internal/database/models"
	"simbirGo/internal/dto"
	"simbirGo/internal/entities"
	"time"
)

type RentRepository interface {
	FindTypeByName(typeName string) uint
	FindTypeById(id uint) string
	FindAvalibleTransports(lat, long, radius float64, typeId uint) []models.Transport
	FindUserById(id uint) models.User
	FindTranspot(id uint) models.Transport
	FindRentById(id int) models.Rent
	SaveTransport(transport models.Transport)
	CreateRent(rent models.Rent) models.Rent
	SaveUser(user models.User)
	FindUserRents(id int) []models.Rent
	FindTransportRents(id int) []models.Rent
	SaveRent(rent models.Rent)
	DeleteRent(id int)
	FindRentTypeById(id uint) string
	FindRentTypeByName(typeName string) uint
}

const (
	minuteUnix float64 = 60
	dayUnix    float64 = 86400
)

type RentUsecase struct {
	r RentRepository
}

func New(r RentRepository) RentUsecase {
	return RentUsecase{r: r}
}

// user's usecase
func (ru RentUsecase) GetAvalibleTransport(lat, long, radius float64, transportType string) ([]entities.Transport, error) {
	typeId := ru.r.FindTypeByName(transportType)
	if typeId == 0 && transportType != "All" {
		return nil, fmt.Errorf("invalid transport type: %s", transportType)
	}
	transportModels := ru.r.FindAvalibleTransports(lat, long, radius, typeId)

	transportEntites := make([]entities.Transport, 0, len(transportModels))
	for _, transport := range transportModels {
		typeName := ru.r.FindTypeById(transport.TypeId)
		transportEntites = append(transportEntites, dto.TransportModelToEntite(transport, typeName))
	}

	return transportEntites, nil
}

func (ru RentUsecase) GetRent(rentId int, userId uint) (entities.Rent, error) {
	rentModel := ru.r.FindRentById(rentId)
	if rentModel.Id == 0 {
		return entities.Rent{}, fmt.Errorf("rent is not exist")
	}

	transport := ru.r.FindTranspot(rentModel.TransportId)

	if userId != rentModel.UserId && userId != transport.OwnerId {
		return entities.Rent{}, fmt.Errorf("rent is not exist")
	}
	rentType := ru.r.FindRentTypeById(rentModel.RentTypeId)

	return dto.RentModelToEntitie(rentModel, rentType), nil
}

func (ru RentUsecase) GetUserHistory(userId uint) []entities.Rent {
	rentModels := ru.r.FindUserRents(int(userId))

	rentEntites := make([]entities.Rent, 0, len(rentModels))
	for _, rent := range rentModels {
		rentType := ru.r.FindRentTypeById(rent.RentTypeId)
		rentEntites = append(rentEntites, dto.RentModelToEntitie(rent, rentType))
	}

	return rentEntites
}

func (ru RentUsecase) GetTransportHistory(userId, transportId int) ([]entities.Rent, error) {
	transport := ru.r.FindTranspot(uint(transportId))
	if transport.Id == 0 {
		return nil, fmt.Errorf("transport is not exist")
	}
	if transport.OwnerId != uint(userId) {
		return nil, fmt.Errorf("transport is not exist")
	}

	rentModels := ru.r.FindTransportRents(transportId)

	rentEntites := make([]entities.Rent, 0, len(rentModels))
	for _, rent := range rentModels {
		rentType := ru.r.FindRentTypeById(rent.RentTypeId)
		rentEntites = append(rentEntites, dto.RentModelToEntitie(rent, rentType))
	}

	return rentEntites, nil
}

func (ru RentUsecase) CreateNewRent(userId uint, transportId int, rentType string) (entities.Rent, error) {
	rentTypeId := ru.r.FindRentTypeByName(rentType)
	if rentTypeId == 0 {
		return entities.Rent{}, fmt.Errorf("type id is not exist")
	}
	transport := ru.r.FindTranspot(uint(transportId))
	if transport.Id == 0 {
		return entities.Rent{}, fmt.Errorf("transport is not exist")
	}

	if !transport.CanBeRented {
		return entities.Rent{}, fmt.Errorf("transport can not be rented")
	}

	if userId == transport.OwnerId {
		return entities.Rent{}, fmt.Errorf("you can not rent own transport")
	}

	if transport.MinutePrice == 0 && transport.DayPrice == 0 {
		return entities.Rent{}, fmt.Errorf("rental price for transport is not indicated")
	}

	var priceOfUnit float64
	switch rentType {
	case "Minutes":
		if transport.MinutePrice == 0 {
			return entities.Rent{}, fmt.Errorf("rental price per minute of transport is not indicated")
		}
		priceOfUnit = transport.MinutePrice
	case "Days":
		if transport.DayPrice == 0 {
			return entities.Rent{}, fmt.Errorf("rental price per day of transport is not indicated")
		}
		priceOfUnit = transport.DayPrice
	}

	rent := models.Rent{
		UserId:      userId,
		TransportId: uint(transportId),
		TimeStart:   time.Now(),
		PriceOfUnit: priceOfUnit,
		RentTypeId:  rentTypeId,
	}
	transport.CanBeRented = false
	ru.r.SaveTransport(transport)
	rent = ru.r.CreateRent(rent)

	fmt.Printf("%+v", rent)
	fmt.Printf("%+v", dto.RentModelToEntitie(rent, rentType))
	return dto.RentModelToEntitie(rent, rentType), nil
}

func (ru RentUsecase) UserEndRent(userId uint, rentId int, lat, long float64) (entities.Rent, error) {
	rentModel := ru.r.FindRentById(rentId)
	if rentModel.Id == 0 || rentModel.UserId != userId {
		return entities.Rent{}, fmt.Errorf("rent is not exist")
	}

	transport := ru.r.FindTranspot(rentModel.TransportId)

	transport.CanBeRented = true
	transport.Latitude = lat
	transport.Longitude = long

	t := time.Now()
	rentModel.TimeEnd = &t

	rentType := ru.r.FindRentTypeById(rentModel.RentTypeId)
	switch rentType {
	case "Minutes":
		rentModel.FinalPrice = ru.calculateRentPrice(rentModel.TimeStart.Unix(),
			rentModel.TimeEnd.Unix(), minuteUnix, rentModel.PriceOfUnit)
	case "Days":
		rentModel.FinalPrice = ru.calculateRentPrice(rentModel.TimeStart.Unix(),
			rentModel.TimeEnd.Unix(), dayUnix, rentModel.PriceOfUnit)
	}

	user := ru.r.FindUserById(rentModel.UserId)
	if rentModel.FinalPrice > user.Balance {
		return entities.Rent{}, fmt.Errorf("not enough money in user's balance")
	}
	user.Balance -= rentModel.FinalPrice

	ru.r.SaveUser(user)
	ru.r.SaveTransport(transport)
	ru.r.SaveRent(rentModel)

	return dto.RentModelToEntitie(rentModel, rentType), nil
}

// admin's usecase
func (ru RentUsecase) AdminGetRent(id int) (entities.Rent, error) {
	rent := ru.r.FindRentById(id)
	if rent.Id == 0 {
		return entities.Rent{}, fmt.Errorf("rent is not exist")
	}
	rentType := ru.r.FindRentTypeById(rent.RentTypeId)

	return dto.RentModelToEntitie(rent, rentType), nil
}

func (ru RentUsecase) AdminGetUserHistory(userId int) ([]entities.Rent, error) {
	user := ru.r.FindUserById(uint(userId))
	if user.Id == 0 {
		return nil, fmt.Errorf("user is not exist")
	}
	rentEntites := ru.GetUserHistory(uint(userId))

	return rentEntites, nil
}

func (ru RentUsecase) AdminGetTransportHistory(transportId int) ([]entities.Rent, error) {
	transport := ru.r.FindTranspot(uint(transportId))

	if transport.Id == 0 {
		return nil, fmt.Errorf("user is not exist")
	}

	rentModels := ru.r.FindUserRents(transportId)

	rentEntites := make([]entities.Rent, 0, len(rentModels))
	for _, rent := range rentModels {
		rentType := ru.r.FindRentTypeById(rent.RentTypeId)
		rentEntites = append(rentEntites, dto.RentModelToEntitie(rent, rentType))
	}

	return rentEntites, nil
}

func (ru RentUsecase) AdminCreateRent(rent entities.Rent) (entities.Rent, error) {
	user := ru.r.FindUserById(rent.UserId)
	if user.Id == 0 {
		return entities.Rent{}, fmt.Errorf("user is not exist")
	}

	transport := ru.r.FindTranspot(rent.TransportId)
	if transport.Id == 0 {
		return entities.Rent{}, fmt.Errorf("transport is not exist")
	}

	if !transport.CanBeRented {
		return entities.Rent{}, fmt.Errorf("transport can not be rented")
	}
	transport.CanBeRented = false

	if rent.TimeEnd != nil {
		switch rent.PriceType {
		case "Minutes":
			rent.FinalPrice = ru.calculateRentPrice(rent.TimeStart.Unix(),
				rent.TimeEnd.Unix(), minuteUnix, rent.PriceOfUnit)
		case "Days":
			rent.FinalPrice = ru.calculateRentPrice(rent.TimeStart.Unix(),
				rent.TimeEnd.Unix(), dayUnix, rent.PriceOfUnit)
		}
	}

	rentTypeId := ru.r.FindRentTypeByName(rent.PriceType)
	if rentTypeId == 0 {
		return entities.Rent{}, fmt.Errorf("invalid price type")
	}

	rentM := dto.RentEntitieToModel(rent, rentTypeId)
	fmt.Printf("%+v\n", rentM)

	rentModel := ru.r.CreateRent(dto.RentEntitieToModel(rent, rentTypeId))
	ru.r.SaveTransport(transport)
	rentEntite := dto.RentModelToEntitie(rentModel, rent.PriceType)

	return rentEntite, nil
}

func (ru RentUsecase) AdminEndRent(id int, lat, long float64) (entities.Rent, error) {
	rentModel := ru.r.FindRentById(id)
	if rentModel.Id == 0 {
		return entities.Rent{}, fmt.Errorf("rent is not exist")
	}

	transport := ru.r.FindTranspot(rentModel.TransportId)
	transport.CanBeRented = true
	transport.Latitude = lat
	transport.Longitude = long

	t := time.Now()
	rentModel.TimeEnd = &t
	rentType := ru.r.FindRentTypeById(rentModel.RentTypeId)
	switch rentType {
	case "Minutes":
		rentModel.FinalPrice = ru.calculateRentPrice(rentModel.TimeStart.Unix(),
			rentModel.TimeEnd.Unix(), minuteUnix, rentModel.PriceOfUnit)
	case "Days":
		rentModel.FinalPrice = ru.calculateRentPrice(rentModel.TimeStart.Unix(),
			rentModel.TimeEnd.Unix(), dayUnix, rentModel.PriceOfUnit)
	}

	user := ru.r.FindUserById(rentModel.UserId)
	if rentModel.FinalPrice > user.Balance {
		return entities.Rent{}, fmt.Errorf("not enough money in user's balance")
	}

	user.Balance -= rentModel.FinalPrice

	ru.r.SaveUser(user)
	ru.r.SaveTransport(transport)
	ru.r.SaveRent(rentModel)

	return dto.RentModelToEntitie(rentModel, rentType), nil
}

func (ru RentUsecase) AdminUpdateRent(rent entities.Rent) (entities.Rent, error) {
	rentModel := ru.r.FindRentById(int(rent.Id))
	if rentModel.Id == 0 {
		return entities.Rent{}, fmt.Errorf("rent is not exist")
	}
	rentTypeId := ru.r.FindRentTypeByName(rent.PriceType)
	if rentTypeId == 0 {
		return entities.Rent{}, fmt.Errorf("price type is not exist")
	}
	rentModel.TransportId = rent.TransportId
	rentModel.UserId = rent.UserId
	rentModel.TimeStart = rent.TimeStart
	rentModel.TimeEnd = rent.TimeEnd
	rentModel.PriceOfUnit = rent.PriceOfUnit
	rentModel.RentTypeId = rentTypeId

	if rentModel.TimeEnd != nil {
		switch rent.PriceType {
		case "Minutes":
			rentModel.FinalPrice = ru.calculateRentPrice(rent.TimeStart.Unix(),
				rent.TimeEnd.Unix(), minuteUnix, rent.PriceOfUnit)
		case "Days":
			rentModel.FinalPrice = ru.calculateRentPrice(rent.TimeStart.Unix(),
				rent.TimeEnd.Unix(), dayUnix, rent.PriceOfUnit)
		}
	}

	ru.r.SaveRent(rentModel)
	return dto.RentModelToEntitie(rentModel, rent.PriceType), nil
}

func (ru RentUsecase) AdminDeleteRent(id int) error {
	rent := ru.r.FindRentById(id)
	if rent.Id == 0 {
		return fmt.Errorf("rent is not exist")
	}

	ru.r.DeleteRent(id)
	return nil
}

func (ru RentUsecase) calculateRentPrice(startTime, endTime int64, timeUnit float64, priceOfUnit float64) float64 {
	return math.Ceil((float64(endTime-startTime) / timeUnit)) * priceOfUnit
}
