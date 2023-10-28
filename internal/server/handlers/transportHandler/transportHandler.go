package transportHandler

import (
	"net/http"
	"simbirGo/internal/entities"
	httpUtil "simbirGo/internal/httputil"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TransportUsecase interface {
	GetTransport(id uint) (entities.Transport, error)
	CreateTransport(transport entities.Transport) (entities.Transport, error)
	UpdateUserTransport(transport entities.Transport) (entities.Transport, error)
	DeleteUserTransport(userId, transportId uint) error
	GetTransports(start, count int, transportType string) ([]entities.Transport, error)
	AdminCreateTransport(transport entities.Transport) (entities.Transport, error)
	AdminUpdateTransport(transport entities.Transport) (entities.Transport, error)
	AdminDeleteTransport(id uint)
}

type TransportHandler struct {
	tu TransportUsecase
}

func New(tu TransportUsecase) TransportHandler {
	return TransportHandler{tu: tu}
}

//user handlers

// @Summary Получение информации о транспотре
// @Tags TransportController
// @Description Просмотр информации о транспорте с id = {id}
// @Produce  json
// @Param id path uint true "Transport id"
// @Success 200 {object} transportHandler.GetTransport.transportData
// @Failure 400 {object} httpUtil.ResponseError
// @Router /api/Transport/{id} [get]
func (th TransportHandler) GetTransport(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 0 {
		httpUtil.NewResponseError(ctx, http.StatusBadRequest, err)
		return
	}
	transport, err := th.tu.GetTransport(uint(id))
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

	ctx.JSON(200, transportData{
		TransportType: transport.TransportType,
		CanBeRented:   transport.CanBeRented,
		Model:         transport.Model,
		Color:         transport.Color,
		Identifier:    transport.Identifier,
		Description:   transport.Description,
		Latitude:      transport.Latitude,
		Longitude:     transport.Longitude,
		MinutePrice:   transport.MinutePrice,
		DayPrice:      transport.DayPrice,
	})
}

// @Summary Создаение транспорта
// @Tags TransportController
// @Description Создает транспорт у текущего авторизованного пользователя
// @Security ApiKeyAuth
// @Accept json
// @Produce  json
// @Param request body transportHandler.CreateTransport.transportData true "Transport data"
// @Success 201 {object} transportHandler.CreateTransport.responseData
// @Failure 400 {object} httpUtil.ResponseError
// @Failure 401 {object} httpUtil.ResponseError
// @Router /api/Transport [post]
func (th TransportHandler) CreateTransport(ctx *gin.Context) {
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

	var tData transportData
	if err := ctx.BindJSON(&tData); err != nil {
		httpUtil.NewResponseError(ctx, 400, err)
		return
	}

	id := ctx.GetUint("id")
	transport := entities.Transport{
		OwnerId:       id,
		TransportType: tData.TransportType,
		CanBeRented:   tData.CanBeRented,
		Model:         tData.Model,
		Color:         tData.Color,
		Identifier:    tData.Identifier,
		Description:   tData.Description,
		Latitude:      tData.Latitude,
		Longitude:     tData.Longitude,
		MinutePrice:   tData.MinutePrice,
		DayPrice:      tData.DayPrice,
	}

	transport, err := th.tu.CreateTransport(transport)
	if err != nil {
		httpUtil.NewResponseError(ctx, 400, err)
		return
	}

	type responseData struct {
		Id            uint
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

	ctx.JSON(200, responseData{
		Id:            transport.Id,
		TransportType: transport.TransportType,
		CanBeRented:   transport.CanBeRented,
		Model:         transport.Model,
		Color:         transport.Color,
		Identifier:    transport.Identifier,
		Description:   transport.Description,
		Latitude:      transport.Latitude,
		Longitude:     transport.Longitude,
		MinutePrice:   transport.MinutePrice,
		DayPrice:      transport.DayPrice,
	})
}

// @Summary Обновление информации о транспотре
// @Tags TransportController
// @Description Обновление информации о транспорте с id = {id}
// @Security ApiKeyAuth
// @Accept json
// @Produce  json
// @Param id path uint true "Transport id"
// @Param request body transportHandler.UpdateTransport.transportData true "Transport data"
// @Success 200 {object} transportHandler.UpdateTransport.responseData
// @Failure 400 {object} httpUtil.ResponseError
// @Failure 401 {object} httpUtil.ResponseError
// @Router /api/Transport/{id} [put]
func (th TransportHandler) UpdateTransport(ctx *gin.Context) {
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

	var tData transportData
	if err := ctx.BindJSON(&tData); err != nil {
		httpUtil.NewResponseError(ctx, 400, err)
		return
	}

	transportIdStr := ctx.Param("id")
	transportId, err := strconv.Atoi(transportIdStr)
	if err != nil || transportId < 0 {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": "invalid value of id param"})
		return
	}
	userId := ctx.GetUint("id")
	transport := entities.Transport{
		Id:            uint(transportId),
		OwnerId:       userId,
		TransportType: tData.TransportType,
		CanBeRented:   tData.CanBeRented,
		Model:         tData.Model,
		Color:         tData.Color,
		Identifier:    tData.Identifier,
		Description:   tData.Description,
		Latitude:      tData.Latitude,
		Longitude:     tData.Longitude,
		MinutePrice:   tData.MinutePrice,
		DayPrice:      tData.DayPrice,
	}
	transport, err = th.tu.UpdateUserTransport(transport)
	if err != nil {
		httpUtil.NewResponseError(ctx, 400, err)
		return
	}

	type responseData struct {
		Id            uint
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
	ctx.JSON(200, responseData{
		Id:            transport.Id,
		TransportType: transport.TransportType,
		CanBeRented:   transport.CanBeRented,
		Model:         transport.Model,
		Color:         transport.Color,
		Identifier:    transport.Identifier,
		Description:   transport.Description,
		Latitude:      transport.Latitude,
		Longitude:     transport.Longitude,
		MinutePrice:   transport.MinutePrice,
		DayPrice:      transport.DayPrice,
	})
}

// @Summary Удаление транспорта
// @Tags TransportController
// @Description Удаление транспорта с id = {id} если данные транспорт принадлежит текущему авторизованному пользователю
// @Security ApiKeyAuth
// @Param id path uint true "Transport id"
// @Success 200
// @Failure 400 {object} httpUtil.ResponseError
// @Failure 401 {object} httpUtil.ResponseError
// @Router /api/Transport/{id} [delete]
func (th TransportHandler) DeleteUserTransport(ctx *gin.Context) {
	transportIdStr := ctx.Param("id")
	transportId, err := strconv.Atoi(transportIdStr)
	if err != nil || transportId < 0 {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": "invalid value of id param"})
		return
	}
	userId := ctx.GetUint("id")
	err = th.tu.DeleteUserTransport(userId, uint(transportId))
	if err != nil {
		httpUtil.NewResponseError(ctx, 400, err)
		return
	}
	ctx.Status(200)
}

// admin handlers

// @Summary Информация о транспортных средствах
// @Tags AdminTransportController
// @Description Получение count транспортных средств с id >= start с типом транспорта transportType
// @Security ApiKeyAuth
// @Produce  json
// @Param start query uint true "start"
// @Param count query uint true "count"
// @Param count query string true "transportType"
// @Success 200 {array} entities.Transport
// @Failure 400 {object} httpUtil.ResponseError
// @Failure 401 {object} httpUtil.ResponseError
// @Failure 403 {object} httpUtil.ResponseError
// @Router /api/Admin/Transport [get]
func (th TransportHandler) AdminGetTransports(ctx *gin.Context) {
	startStr := ctx.Query("start")
	countStr := ctx.Query("count")
	transportType := ctx.Query("transportType")
	start, err := strconv.Atoi(startStr)
	if err != nil || start < 0 {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": "invalid value of start query param"})
		return
	}

	count, err := strconv.Atoi(countStr)
	if err != nil || count < 0 {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": "invalid value of count query param"})
		return
	}

	transports, err := th.tu.GetTransports(start, count, transportType)
	if err != nil {
		httpUtil.NewResponseError(ctx, 400, err)
		return
	}

	ctx.JSON(200, transports)
}

// @Summary Информация о транспортном средстве
// @Tags AdminTransportController
// @Description Получение информации о транспортном средстве с id = {id}
// @Security ApiKeyAuth
// @Produce json
// @Param id path uint true "Transport id"
// @Success 200 {object} entities.Transport
// @Failure 400 {object} httpUtil.ResponseError
// @Failure 401 {object} httpUtil.ResponseError
// @Failure 403 {object} httpUtil.ResponseError
// @Router /api/Admin/Transport/{id} [get]
func (th TransportHandler) AdminGetTransport(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 0 {
		httpUtil.NewResponseError(ctx, http.StatusBadRequest, err)
		return
	}
	transport, err := th.tu.GetTransport(uint(id))
	if err != nil {
		httpUtil.NewResponseError(ctx, 400, err)
		return
	}

	ctx.JSON(200, transport)
}

// @Summary Создание транспортного средства
// @Tags AdminTransportController
// @Description Создание транспортного средства указывая пользователя с id = ownerId
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param request body transportHandler.AdminCreateTransport.transportData true "Transport data"
// @Success 201 {object} entities.Transport
// @Failure 400 {object} httpUtil.ResponseError
// @Failure 401 {object} httpUtil.ResponseError
// @Failure 403 {object} httpUtil.ResponseError
// @Router /api/Admin/Transport [post]
func (th TransportHandler) AdminCreateTransport(ctx *gin.Context) {
	type transportData struct {
		OwnerId       uint    `json:"ownerId" binding:"required"`
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

	var tData transportData
	if err := ctx.BindJSON(&tData); err != nil {
		httpUtil.NewResponseError(ctx, 400, err)
		return
	}

	transport := entities.Transport{
		OwnerId:       tData.OwnerId,
		TransportType: tData.TransportType,
		CanBeRented:   tData.CanBeRented,
		Model:         tData.Model,
		Color:         tData.Color,
		Identifier:    tData.Identifier,
		Description:   tData.Description,
		Latitude:      tData.Latitude,
		Longitude:     tData.Longitude,
		MinutePrice:   tData.MinutePrice,
		DayPrice:      tData.DayPrice,
	}

	transport, err := th.tu.AdminCreateTransport(transport)
	if err != nil {
		httpUtil.NewResponseError(ctx, 400, err)
		return
	}

	ctx.JSON(201, transport)
}

// @Summary Обновление транспортного средства
// @Tags AdminTransportController
// @Description Обновление транспортного средства с id = {id}
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param id path uint true "Transport id"
// @Param request body transportHandler.AdminUpdateTransport.transportData true "Transport data"
// @Success 200 {object} entities.Transport
// @Failure 400 {object} httpUtil.ResponseError
// @Failure 401 {object} httpUtil.ResponseError
// @Failure 403 {object} httpUtil.ResponseError
// @Router /api/Admin/Transport/{id} [put]
func (th TransportHandler) AdminUpdateTransport(ctx *gin.Context) {
	type transportData struct {
		OwnerId       uint    `json:"ownerId" binding:"required"`
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

	var tData transportData
	if err := ctx.BindJSON(&tData); err != nil {
		httpUtil.NewResponseError(ctx, 400, err)
		return
	}

	transportIdStr := ctx.Param("id")
	transportId, err := strconv.Atoi(transportIdStr)
	if err != nil || transportId < 0 {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": "invalid value of id param"})
		return
	}

	transport := entities.Transport{
		Id:            uint(transportId),
		OwnerId:       tData.OwnerId,
		TransportType: tData.TransportType,
		CanBeRented:   tData.CanBeRented,
		Model:         tData.Model,
		Color:         tData.Color,
		Identifier:    tData.Identifier,
		Description:   tData.Description,
		Latitude:      tData.Latitude,
		Longitude:     tData.Longitude,
		MinutePrice:   tData.MinutePrice,
		DayPrice:      tData.DayPrice,
	}

	transport, err = th.tu.AdminUpdateTransport(transport)
	if err != nil {
		httpUtil.NewResponseError(ctx, 400, err)
		return
	}
	ctx.JSON(200, transport)
}

// @Summary Удаление транспортного средства
// @Tags AdminTransportController
// @Description Удаление транспортного средства с id = {id}
// @Security ApiKeyAuth
// @Param id path uint true "Transport id"
// @Success 200
// @Failure 400 {object} httpUtil.ResponseError
// @Failure 401 {object} httpUtil.ResponseError
// @Failure 403 {object} httpUtil.ResponseError
// @Router /api/Admin/Transport/{id} [delete]
func (th TransportHandler) DeleteTransport(ctx *gin.Context) {
	transportIdStr := ctx.Param("id")
	transportId, err := strconv.Atoi(transportIdStr)
	if err != nil || transportId < 0 {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": "invalid value of id param"})
		return
	}

	th.tu.AdminDeleteTransport(uint(transportId))
}
