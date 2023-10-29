package rentUsecase

import (
	"fmt"
	"simbirGo/internal/database/models"
	"simbirGo/internal/dto"
	"simbirGo/internal/entities"
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

func (ru RentUsecase) AdminGetRent(id int) (entities.Rent, error) {
	rent := ru.r.FindRentById(id)
	if rent.Id == 0 {
		return entities.Rent{}, fmt.Errorf("rent is not exist")
	}
	return dto.RentModelToEntitie(rent), nil
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
	ru.r.SaveTransport(transport)

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
	fmt.Println(rent.FinalPrice)
	rentModel := ru.r.CreateRent(dto.RentEntitieToModel(rent))
	rentEntite := dto.RentModelToEntitie(rentModel)

	return rentEntite, nil
}

func (ru RentUsecase) calculateRentPrice(startTime, endTime int64, timeUnit float64, priceOfUnit float64) float64 {
	diff := endTime - startTime
	fmt.Println("diff = ", diff)

	price := (float64(endTime-startTime) / timeUnit) * priceOfUnit
	return price
}
