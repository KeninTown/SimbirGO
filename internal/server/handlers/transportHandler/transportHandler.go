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
}

type TransportHandler struct {
	tu TransportUsecase
}

func New(tu TransportUsecase) TransportHandler {
	return TransportHandler{tu: tu}
}

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

	ctx.JSON(200, transport)
}

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

	ctx.JSON(201, transport)
}

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

	ctx.JSON(200, transport)
}

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
