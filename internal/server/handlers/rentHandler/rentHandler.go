package rentHandler

import (
	"fmt"
	"math"
	"net/http"
	"simbirGo/internal/entities"
	httpUtil "simbirGo/internal/httputil"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type RentUsecase interface {
	GetAvalibleTransport(lat, long, radius float64, transportType string) ([]entities.Transport, error)
	AdminGetRent(id int) (entities.Rent, error)
	AdminCreateRent(rent entities.Rent) (entities.Rent, error)
}

type RentHandler struct {
	ru RentUsecase
}

func New(ru RentUsecase) RentHandler {
	return RentHandler{ru: ru}
}

// @Summary Транспорт для аренды
// @Tags RentController
// @Description Получение транспорта доступного для аренды по параметрам его расположения и типу транспорта
// @Security ApiKeyAuth
// @Produce  json
// @Param lat query float64 true "lat"
// @Param radius query float64 true "radius"
// @Param long query float64 true "long"
// @Param transportType query string true "transportType"
// @Success 200 {array} entities.Transport
// @Failure 400 {object} httpUtil.ResponseError
// @Router /api/Rent/Transport [get]
func (rh RentHandler) GetAvalibleTransport(ctx *gin.Context) {
	latStr := ctx.Query("lat")
	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil || math.Abs(lat) > 90 {
		httpUtil.NewResponseError(ctx, 400, fmt.Errorf("invalid value of lat query param"))
	}

	longStr := ctx.Query("long")
	long, err := strconv.ParseFloat(longStr, 64)
	if err != nil || math.Abs(long) > 180 {
		httpUtil.NewResponseError(ctx, 400, fmt.Errorf("invalid value of long query param"))
	}

	radiusStr := ctx.Query("radius")
	radius, err := strconv.ParseFloat(radiusStr, 64)
	if err != nil || radius < 0 {
		httpUtil.NewResponseError(ctx, 400, fmt.Errorf("invalid value of long radius query param"))
	}
	transportType, ok := ctx.GetQuery("transportType")
	if !ok {
		httpUtil.NewResponseError(ctx, 400, fmt.Errorf("invalid value of transport type query param"))
	}

	transports, err := rh.ru.GetAvalibleTransport(lat, long, radius, transportType)
	if err != nil {
		httpUtil.NewResponseError(ctx, 400, err)
		return
	}

	type transportData struct {
		TransportType string  `json:"transportType" binding:"required"`
		CanBeRented   bool    `json:"canBeRented"`
		Model         string  `json:"model" binding:"required"`
		Color         string  `json:"color" binding:"required"`
		Identifier    string  `json:"identifier" binding:"required"`
		Description   string  `json:"description"`
		Latitude      float64 `json:"latitude" binding:"required"`
		Longitude     float64 `json:"longitude" binding:"required"`
		MinutePrice   float64 `json:"minutePrice"`
		DayPrice      float64 `json:"dayPrice"`
	}

	//create transportDomainDto or smth
	tData := make([]transportData, 0, len(transports))
	for _, tr := range transports {
		tData = append(tData, transportData{
			TransportType: tr.TransportType,
			CanBeRented:   tr.CanBeRented,
			Model:         tr.Model,
			Color:         tr.Color,
			Identifier:    tr.Identifier,
			Description:   tr.Description,
			Latitude:      tr.Latitude,
			Longitude:     tr.Longitude,
			MinutePrice:   tr.MinutePrice,
			DayPrice:      tr.DayPrice,
		})
	}

	ctx.JSON(200, tData)
}

// admin routes

// @Summary Получение аренды
// @Tags AdminRentController
// @Description Получение аренды с id = {id}
// @Security ApiKeyAuth
// @Accept json
// @Produce  json
// @Param id path uint true "Transport id"
// @Success 200 {object} entities.Rent
// @Failure 400 {object} httpUtil.ResponseError
// @Failure 401 {object} httpUtil.ResponseError
// @Failure 403 {object} httpUtil.ResponseError
// @Router /api/Admin/Rent/{id} [get]
func (rh RentHandler) AdminGetRent(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 0 {
		httpUtil.NewResponseError(ctx, http.StatusBadRequest, err)
		return
	}
	rent, err := rh.ru.AdminGetRent(id)
	if err != nil {
		httpUtil.NewResponseError(ctx, 400, err)
		return
	}

	ctx.JSON(200, rent)
}

// @Summary Создание аренды
// @Tags AdminRentController
// @Description Аредна транспорта с id = transportId пользователем с id = userId
// @Security ApiKeyAuth
// @Accept json
// @Produce  json
// @Param request body rentHandler.AdminCreateRent.rentData true "rent data"
// @Success 201 {object} entities.Rent
// @Failure 400 {object} httpUtil.ResponseError
// @Failure 401 {object} httpUtil.ResponseError
// @Failure 403 {object} httpUtil.ResponseError
// @Router /api/Admin/Rent [post]
func (rh RentHandler) AdminCreateRent(ctx *gin.Context) {
	type rentData struct {
		TransportId uint    `json:"transportId" binding:"required"`
		UserId      uint    `json:"userId" binding:"required"`
		TimeStart   string  `json:"timeStart" binding:"required"`
		TimeEnd     string  `json:"timeEnd"`
		PriceOfUnit float64 `json:"priceOfUnit" binding:"required"`
		PriceType   string  `json:"priceType" binding:"required" enum:"Days, Minutes"`
	}
	var rData rentData
	if err := ctx.BindJSON(&rData); err != nil {
		httpUtil.NewResponseError(ctx, 400, err)
		return
	}

	if rData.PriceType != "Days" && rData.PriceType != "Minutes" {
		httpUtil.NewResponseError(ctx, 400, fmt.Errorf("invalid priceType"))
		return
	}

	timeStart, err := time.Parse(time.RFC3339, rData.TimeStart)
	if err != nil {
		httpUtil.NewResponseError(ctx, 400, fmt.Errorf("invalid end time value, should be : yyyy-mm-ddThh:mm:ss±hh:mm"))
		return
	}

	var timeEnd *time.Time
	if rData.TimeEnd != "" {
		t, err := time.Parse(time.RFC3339, rData.TimeEnd)
		if err != nil && rData.TimeEnd != "" {
			httpUtil.NewResponseError(ctx, 400, fmt.Errorf("invalid end time value, should be : yyyy-mm-ddThh:mm:ss±hh:mm"))
			return
		}

		if t.Unix()-timeStart.Unix() <= 0 {
			httpUtil.NewResponseError(ctx, 400, fmt.Errorf("invalid end time value: end time must be after start time"))
			return
		}

		timeEnd = &t
	}

	rent := entities.Rent{
		TransportId: rData.TransportId,
		UserId:      rData.UserId,
		TimeStart:   timeStart,
		TimeEnd:     timeEnd,
		PriceOfUnit: rData.PriceOfUnit,
		PriceType:   rData.PriceType,
	}

	rent, err = rh.ru.AdminCreateRent(rent)
	if err != nil {
		httpUtil.NewResponseError(ctx, 400, err)
		return
	}

	ctx.JSON(201, rent)
}
