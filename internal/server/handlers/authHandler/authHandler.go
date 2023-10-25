package authHandler

import (
	"fmt"
	"net/http"
	"simbirGo/internal/entities"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AuthUsecase interface {
	MyAccount(id uint) (entities.User, error)
	SignIn(user entities.User) (string, error)
	SignUp(user entities.User) (string, error)
	Update(user entities.User) (entities.User, error)

	//admin's func
	GetUsers(start, end uint) []entities.User
	CreateUser(user entities.User) (entities.User, error)
	UpdateUser(user entities.User) (entities.User, error)
	DeleteUser(id uint)
}

type AuthHandlers struct {
	uc AuthUsecase
}

func New(uc AuthUsecase) AuthHandlers {
	return AuthHandlers{uc: uc}
}

func (ah AuthHandlers) MyAccount(ctx *gin.Context) {
	id := ctx.GetUint("id")
	fmt.Println(id)
	user, err := ah.uc.MyAccount(id)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, user)
}

func (ah AuthHandlers) SignIn(ctx *gin.Context) {
	var user entities.User
	if err := ctx.BindJSON(&user); err != nil {
		ctx.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}
	token, err := ah.uc.SignIn(user)
	if err != nil {
		ctx.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}
	ctx.SetCookie("access_token", token, 3600, "/", "localhost", false, true)
	ctx.JSON(201, gin.H{"token": token})
}

func (ah AuthHandlers) SignUp(ctx *gin.Context) {
	var user entities.User
	if err := ctx.BindJSON(&user); err != nil {
		ctx.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}
	user.Balance = 0
	token, err := ah.uc.SignUp(user)
	if err != nil {
		ctx.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}
	ctx.SetCookie("access_token", token, 3600, "/", "localhost", false, true)
	ctx.JSON(201, gin.H{"token": token})
}

func (ah AuthHandlers) SignOut(ctx *gin.Context) {
	ctx.SetCookie("access_token", "", 1, "/", "localhost", false, true)
}

func (ah AuthHandlers) Update(ctx *gin.Context) {
	id := ctx.GetUint("id")
	var user entities.User
	if err := ctx.BindJSON(&user); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": "invalid data"})
		return
	}
	user.Id = id
	user, err := ah.uc.Update(user)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, user)
}

func (ah AuthHandlers) GetUsers(ctx *gin.Context) {
	startStr := ctx.Query("start")
	countStr := ctx.Query("count")
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

	users := ah.uc.GetUsers(uint(start), uint(count))

	ctx.JSON(http.StatusOK, users)
}

func (ah AuthHandlers) GetUser(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 0 {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": "invalid value of id param"})
		return
	}
	user, err := ah.uc.MyAccount(uint(id))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, user)
}

func (ah AuthHandlers) CreateUser(ctx *gin.Context) {
	var user entities.User
	if err := ctx.BindJSON(&user); err != nil {
		ctx.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}
	user, err := ah.uc.CreateUser(user)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, user)
}

func (ah AuthHandlers) UpdateUser(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 0 {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": "invalid value of id param"})
		return
	}
	var user entities.User
	if err := ctx.BindJSON(&user); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": "invalid data"})
		return
	}
	user.Id = uint(id)
	user, err = ah.uc.UpdateUser(user)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, user)
}

func (ah AuthHandlers) DeleteUser(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 0 {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": "invalid value of id param"})
		return
	}
	ah.uc.DeleteUser(uint(id))

	ctx.Status(http.StatusOK)
}
