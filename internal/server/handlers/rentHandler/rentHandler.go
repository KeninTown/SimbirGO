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
	//user
	GetAvalibleTransport(lat, long, radius float64, transportType string) ([]entities.Transport, error)
	GetRent(rentId int, userId uint) (entities.Rent, error)
	GetUserHistory(userId int) []entities.Rent
	GetTransportHistory(userId, transportId int) ([]entities.Rent, error)
	CreateNewRent(userId, transportId int, rentType string) (entities.Rent, error)
	UserEndRent(userId, rentId int, lat, long float64) (entities.Rent, error)

	//admin usecase
	AdminGetRent(id int) (entities.Rent, error)
	AdminGetUserHistory(userId int) ([]entities.Rent, error)
	AdminGetTransportHistory(transportId int) ([]entities.Rent, error)
	AdminCreateRent(rent entities.Rent) (entities.Rent, error)
	AdminEndRent(id int, lat, long float64) (entities.Rent, error)
	AdminUpdateRent(rent entities.Rent) (entities.Rent, error)
	AdminDeleteRent(id int) error
}

type RentHandler struct {
	ru RentUsecase
}

func New(ru RentUsecase) RentHandler {
	return RentHandler{ru: ru}
}

//user handlers

// @Summary Транспорт для аренды
// @Tags RentController
// @Description Получение транспорта доступного для аренды по параметрам его расположения и типу транспорта
// @Produce json
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

// @Summary Получение аренды пользователя
// @Tags RentController
// @Description Получение данных аренды с id = {rentid}. Данные могут получить только арендатор и арендодатель.
// @Security ApiKeyAuth
// @Produce json
// @Param rentid path uint true "Rent id"
// @Success 200 {object} entities.Rent
// @Failure 400 {object} httpUtil.ResponseError
// @Failure 401 {object} httpUtil.ResponseError
// @Router /api/Rent/{rentid} [get]
func (rh RentHandler) UserGetRent(ctx *gin.Context) {
	rentIdStr := ctx.Param("id")
	rentId, err := strconv.Atoi(rentIdStr)
	if err != nil || rentId < 0 {
		httpUtil.NewResponseError(ctx, http.StatusBadRequest, err)
		return
	}
	userId := ctx.GetUint("id")
	rent, err := rh.ru.GetRent(rentId, userId)
	if err != nil {
		httpUtil.NewResponseError(ctx, http.StatusBadRequest, err)
		return
	}

	ctx.JSON(200, rent)
}

// @Summary Получение истории аредны пользователя
// @Tags RentController
// @Description Получение истории всех аренд пользователя
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {array} entities.Rent
// @Failure 400 {object} httpUtil.ResponseError
// @Failure 401 {object} httpUtil.ResponseError
// @Router /api/Rent/MyHistory [get]
func (rh RentHandler) UserGetHistory(ctx *gin.Context) {
	userId := ctx.GetInt("id")

	rent := rh.ru.GetUserHistory(userId)

	ctx.JSON(200, rent)
}

// @Summary Получение истории аредны транспорта
// @Tags RentController
// @Description Получение истории всех аренд транспорта с id = {transportId}.
// @Description Данные может получить только владелец транспорта.
// @Security ApiKeyAuth
// @Produce json
// @Param transportId path uint true "Transport id"
// @Success 200 {array} entities.Rent
// @Failure 400 {object} httpUtil.ResponseError
// @Failure 401 {object} httpUtil.ResponseError
// @Router /api/Rent/TransportHistory/{transportId} [get]
func (rh RentHandler) UserGetTransportHistory(ctx *gin.Context) {
	userId := ctx.GetInt("id")
	transportIdStr := ctx.Param("id")
	transportId, err := strconv.Atoi(transportIdStr)
	if err != nil || transportId < 0 {
		httpUtil.NewResponseError(ctx, http.StatusBadRequest, err)
		return
	}

	rents, err := rh.ru.GetTransportHistory(userId, transportId)
	if err != nil {
		httpUtil.NewResponseError(ctx, 400, err)
		return
	}

	ctx.JSON(200, rents)
}

// @Summary Новая аренда транспорта
// @Tags RentController
// @Description Создание новой аредны транспорта с id = {transportid}.
// @Description В параметра rentType указывается тип аренды: [Minutes, Days].
// @Security ApiKeyAuth
// @Produce json
// @Param transportId path uint true "Transport id"
// @Param rentType query string true "Rent type: [Minutes, Days]"
// @Success 201 {object} entities.Rent
// @Failure 400 {object} httpUtil.ResponseError
// @Failure 401 {object} httpUtil.ResponseError
// @Router /api/Rent/New/{transportId} [post]
func (rh RentHandler) UserCreateNewRent(ctx *gin.Context) {
	transportIdStr := ctx.Param("id")
	transportId, err := strconv.Atoi(transportIdStr)
	if err != nil || transportId < 0 {
		httpUtil.NewResponseError(ctx, http.StatusBadRequest, err)
		return
	}

	rentType, ok := ctx.GetQuery("rentType")
	if !ok || (rentType != "Minutes" && rentType != "Days") {
		httpUtil.NewResponseError(ctx, 400, fmt.Errorf("invalid value of transport type query param"))
		return
	}

	userId := ctx.GetInt("id")

	rent, err := rh.ru.CreateNewRent(userId, transportId, rentType)

	if err != nil {
		httpUtil.NewResponseError(ctx, 400, err)
	}

	ctx.JSON(201, rent)
}

// @Summary Окончание аренды
// @Tags RentController
// @Description Окончание аренды транспорта с id = {transportid}.
// @Description Происходит рассчет итоговой суммы аренды и если она оказывается больше,
// @Description чем сумма на счете пользователя, то в завершить аренду нельзя.
// @Security ApiKeyAuth
// @Produce json
// @Param rentId path uint true "Transport id"
// @Param lat query float64 true "lat"
// @Param long query float64 true "long"
// @Success 201 {object} entities.Rent
// @Failure 400 {object} httpUtil.ResponseError
// @Failure 401 {object} httpUtil.ResponseError
// @Router /api/Rent/End/{rentId} [post]
func (rh RentHandler) UserEndRent(ctx *gin.Context) {
	rentIdStr := ctx.Param("id")
	rentId, err := strconv.Atoi(rentIdStr)
	if err != nil || rentId < 0 {
		httpUtil.NewResponseError(ctx, http.StatusBadRequest, err)
		return
	}

	latStr := ctx.Query("lat")
	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil || math.Abs(lat) > 90 {
		httpUtil.NewResponseError(ctx, 400, fmt.Errorf("invalid value of lat query param"))
		return
	}

	longStr := ctx.Query("long")
	long, err := strconv.ParseFloat(longStr, 64)
	if err != nil || math.Abs(long) > 180 {
		httpUtil.NewResponseError(ctx, 400, fmt.Errorf("invalid value of long query param"))
		return
	}

	userId := ctx.GetInt("id")

	rent, err := rh.ru.UserEndRent(userId, rentId, lat, long)
	if err != nil {
		httpUtil.NewResponseError(ctx, 400, err)
		return
	}

	ctx.JSON(201, rent)
}

// admin handlers

// @Summary Получение аренды
// @Tags RentControllerAdmin
// @Description Получение аренды с id = {rentid}
// @Security ApiKeyAuth
// @Produce  json
// @Param rentid path uint true "Rent id"
// @Success 200 {object} entities.Rent
// @Failure 400 {object} httpUtil.ResponseError
// @Failure 401 {object} httpUtil.ResponseError
// @Failure 403 {object} httpUtil.ResponseError
// @Router /api/Admin/Rent/{rentid} [get]
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

// @Summary Получение аренды
// @Tags RentControllerAdmin
// @Description Получение истории всех аренд пользователем с id = {userid}
// @Security ApiKeyAuth
// @Produce  json
// @Param userid path uint true "User id"
// @Success 200 {array} entities.Rent
// @Failure 400 {object} httpUtil.ResponseError
// @Failure 401 {object} httpUtil.ResponseError
// @Failure 403 {object} httpUtil.ResponseError
// @Router /api/Admin/UserHistory/{userid} [get]
func (rh RentHandler) AdminGetUserHistory(ctx *gin.Context) {
	userIdStr := ctx.Param("id")
	userId, err := strconv.Atoi(userIdStr)
	if err != nil || userId < 0 {
		httpUtil.NewResponseError(ctx, http.StatusBadRequest, err)
		return
	}
	rents, err := rh.ru.AdminGetUserHistory(userId)
	if err != nil {
		httpUtil.NewResponseError(ctx, 400, err)
		return
	}

	ctx.JSON(200, rents)
}

// @Summary Получение аренды
// @Tags RentControllerAdmin
// @Description Получение истории всех аренд пользователем с id = {transportId}
// @Security ApiKeyAuth
// @Produce  json
// @Param transportId path uint true "Transport id"
// @Success 200 {array} entities.Rent
// @Failure 400 {object} httpUtil.ResponseError
// @Failure 401 {object} httpUtil.ResponseError
// @Failure 403 {object} httpUtil.ResponseError
// @Router /api/Admin/TransportHistory/{transportId} [get]
func (rh RentHandler) AdminGetTransportHistory(ctx *gin.Context) {
	transportIdStr := ctx.Param("id")
	transportId, err := strconv.Atoi(transportIdStr)
	if err != nil || transportId < 0 {
		httpUtil.NewResponseError(ctx, http.StatusBadRequest, err)
		return
	}

	rents, err := rh.ru.AdminGetTransportHistory(transportId)

	if err != nil {
		httpUtil.NewResponseError(ctx, 400, err)
		return
	}

	ctx.JSON(200, rents)
}

// @Summary Создание новой аренды
// @Tags RentControllerAdmin
// @Description Создание аренды транспорта с id = transportId пользователем с id = userId
// @Security ApiKeyAuth
// @Accept json
// @Produce  json
// @Param request body rentHandler.AdminCreateRent.rentData true "Rent data"
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
		tEnd, err := time.Parse(time.RFC3339, rData.TimeEnd)
		if err != nil && rData.TimeEnd != "" {
			httpUtil.NewResponseError(ctx, 400, fmt.Errorf("invalid end time value, should be : yyyy-mm-ddThh:mm:ss±hh:mm"))
			return
		}

		if tEnd.Unix()-timeStart.Unix() <= 0 {
			httpUtil.NewResponseError(ctx, 400, fmt.Errorf("invalid end time value: end time must be after start time"))
			return
		}

		timeEnd = &tEnd
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

// @Summary Завершение аренды
// @Tags RentControllerAdmin
// @Description Завершение аренды транспорта с id = {rentId}.
// @Description Происходит рассчет итоговой суммы аренды и если она оказывается больше,
// @Description чем сумма на счете пользователя, то в завершить аренду нельзя.
// @Security ApiKeyAuth
// @Produce  json
// @Param rentId path uint true "Rent id"
// @Success 201 {object} entities.Rent
// @Failure 400 {object} httpUtil.ResponseError
// @Failure 401 {object} httpUtil.ResponseError
// @Failure 403 {object} httpUtil.ResponseError
// @Router /api/Admin/Rent/End/{rentId} [post]
func (rh RentHandler) AdminEndRent(ctx *gin.Context) {
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

	rentIdStr := ctx.Param("id")
	rentId, err := strconv.Atoi(rentIdStr)
	if err != nil || rentId < 0 {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": "invalid value of id param"})
		return
	}

	rent, err := rh.ru.AdminEndRent(rentId, lat, long)
	if err != nil {
		httpUtil.NewResponseError(ctx, 400, err)
		return
	}

	ctx.JSON(201, rent)
}

// @Summary Обновление аренды
// @Tags RentControllerAdmin
// @Description Обновление аренды с id = {rentId}
// @Description Если в обновлении аренды указывается дата ее окончания, то аренда считается завершенной.
// @Description Происходит рассчет итоговой суммы аренды и если она оказывается больше,
// @Description чем сумма на счете пользователя, то в обновить аренду нельзя.
// @Security ApiKeyAuth
// @Accept json
// @Produce  json
// @Param request body rentHandler.AdminUpdateRent.rentData true "Rent data"
// @Param rentId path uint true "Rent id"
// @Success 201 {object} entities.Rent
// @Failure 400 {object} httpUtil.ResponseError
// @Failure 401 {object} httpUtil.ResponseError
// @Failure 403 {object} httpUtil.ResponseError
// @Router /api/Admin/Rent/{rentId} [put]
func (rh RentHandler) AdminUpdateRent(ctx *gin.Context) {
	rentIdStr := ctx.Param("id")
	rentId, err := strconv.Atoi(rentIdStr)
	if err != nil || rentId < 0 {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": "invalid value of id param"})
		return
	}

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
		tEnd, err := time.Parse(time.RFC3339, rData.TimeEnd)
		if err != nil && rData.TimeEnd != "" {
			httpUtil.NewResponseError(ctx, 400, fmt.Errorf("invalid end time value, should be : yyyy-mm-ddThh:mm:ss±hh:mm"))
			return
		}

		if tEnd.Unix()-timeStart.Unix() <= 0 {
			httpUtil.NewResponseError(ctx, 400, fmt.Errorf("invalid end time value: end time must be after start time"))
			return
		}

		timeEnd = &tEnd
	}

	rent := entities.Rent{
		TransportId: rData.TransportId,
		UserId:      rData.UserId,
		TimeStart:   timeStart,
		TimeEnd:     timeEnd,
		PriceOfUnit: rData.PriceOfUnit,
		PriceType:   rData.PriceType,
	}

	rent, err = rh.ru.AdminUpdateRent(rent)
	if err != nil {
		httpUtil.NewResponseError(ctx, 400, err)
		return
	}

	ctx.JSON(201, rent)
}

// @Summary Удаление аренды
// @Tags RentControllerAdmin
// @Description Удаление аренды с id = {rentId}
// @Security ApiKeyAuth
// @Produce json
// @Param rentId path uint true "Rent id"
// @Success 200
// @Failure 400 {object} httpUtil.ResponseError
// @Failure 401 {object} httpUtil.ResponseError
// @Failure 403 {object} httpUtil.ResponseError
// @Router /api/Admin/Rent/{rentId} [delete]
func (rh RentHandler) AdminDeleteRent(ctx *gin.Context) {
	rentIdStr := ctx.Param("id")
	rentId, err := strconv.Atoi(rentIdStr)
	if err != nil || rentId < 0 {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": "invalid value of id param"})
		return
	}

	err = rh.ru.AdminDeleteRent(rentId)
	if err != nil {
		httpUtil.NewResponseError(ctx, 400, err)
		return
	}

	ctx.Status(200)
}
