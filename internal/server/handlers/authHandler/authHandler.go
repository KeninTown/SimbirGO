package authHandler

import (
	"net/http"
	"simbirGo/internal/entities"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AuthUsecase interface {
	//user's cases
	MyAccount(id uint) (entities.User, error)
	SignIn(user entities.User) (string, error)
	SignUp(user entities.User) (entities.User, string, error)
	Update(user entities.User) (entities.User, error)

	//admin's cases
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

// @Summary Просмотр данных текущего аккаунта
// @Tags AccountController
// @Description Просмотр информации о текущем авторизованном аккаунте
// @Produce  json
// @Security ApiKeyAuth
// @Success 200 {object} entities.User
// @Failure 400 {object} httpUtil.ResponseError
// @Failure 401 {object} httpUtil.ResponseError
// @Router /api/Account/Me [get]
func (ah AuthHandlers) MyAccount(ctx *gin.Context) {
	id := ctx.GetUint("id")
	user, err := ah.uc.MyAccount(id)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, user)
}

// @Summary Вход в аккаунт
// @Tags AccountController
// @Description Вход в аккаунт и установление в cookie jwt токена авторизации
// @Accept json
// @Produce  json
// @Param request body authHandler.SignIn.userCreadentials true "User credentials"
// @Success 201 {object} string
// @Failure 400 {object} httpUtil.ResponseError
// @Router /api/Account/SignIn [post]
func (ah AuthHandlers) SignIn(ctx *gin.Context) {
	type userCreadentials struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	userCred := userCreadentials{}
	if err := ctx.BindJSON(&userCred); err != nil {
		ctx.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	user := entities.User{Username: userCred.Username, Password: userCred.Password}
	token, err := ah.uc.SignIn(user)
	if err != nil {
		ctx.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}
	ctx.SetCookie("access_token", token, 3600, "/", "localhost", false, true)
	ctx.JSON(201, gin.H{"token": token})
}

// @Summary Регистрация
// @Tags AccountController
// @Description Регистрация и установление jwt токена в cookie access_token
// @Accept json
// @Produce  json
// @Param request body authHandler.SignUp.userData true "User data"
// @Success 201 {object} entities.User
// @Failure 400 {object} httpUtil.ResponseError
// @Router /api/Account/SignUp [post]
func (ah AuthHandlers) SignUp(ctx *gin.Context) {
	type userData struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
		IsAdmin  bool   `json:"isAdmin"`
	}
	var usData userData
	if err := ctx.BindJSON(&usData); err != nil {
		ctx.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}
	user := entities.User{
		Username: usData.Username,
		Password: usData.Password,
		IsAdmin:  usData.IsAdmin,
	}

	user, token, err := ah.uc.SignUp(user)
	if err != nil {
		ctx.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}
	ctx.SetCookie("access_token", token, 360000, "/", "localhost", false, true)
	ctx.JSON(201, user)
}

// @Summary Выход из аккаунта
// @Tags AccountController
// @Description Удаление jwt токена из cookie access_token
// @Security ApiKeyAuth
// @Success 200
// @Failure 400 {object} httpUtil.ResponseError
// @Router /api/Account/SignOut [post]
func (ah AuthHandlers) SignOut(ctx *gin.Context) {
	ctx.SetCookie("access_token", "", 1, "/", "localhost", false, true)
}

// @Summary Обновление данных аккаунта
// @Tags AccountController
// @Description Проверка входящих данных и обновление данных аккаунта
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param request body authHandler.Update.userData true "User data"
// @Success 200 {object} entities.User
// @Failure 400 {object} httpUtil.ResponseError
// @Failure 401 {object} httpUtil.ResponseError
// @Router /api/Account/Update [put]
func (ah AuthHandlers) Update(ctx *gin.Context) {
	type userData struct {
		Username string  `json:"username" binding:"required"`
		Password string  `json:"password" binding:"required"`
		IsAdmin  bool    `json:"isAdmin" binding:"required"`
		Balance  float64 `json:"balance" binding:"required"`
	}
	var data userData

	if err := ctx.BindJSON(&data); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": "invalid data"})
		return
	}
	id := ctx.GetUint("id")
	user := entities.User{
		Id:       id,
		Username: data.Username,
		Password: data.Password,
		IsAdmin:  data.IsAdmin,
		Balance:  data.Balance,
	}
	user, err := ah.uc.Update(user)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, user)
}

// admin handlers

// @Summary Получение данных пользователей
// @Tags AdminAccountController
// @Description Получение данных count пользователей начиная с id = start
// @Security ApiKeyAuth
// @Produce json
// @Param start query uint true "start"
// @Param count query uint true "count"
// @Success 200 {array} entities.User
// @Failure 400 {object} httpUtil.ResponseError
// @Failure 401 {object} httpUtil.ResponseError
// @Failure 403 {object} httpUtil.ResponseError
// @Router /api/Admin/Account [get]
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

// @Summary Получение пользователя
// @Tags AdminAccountController
// @Description Получение пользователя с id = {id}
// @Security ApiKeyAuth
// @Produce json
// @Param id path uint true "Account id"
// @Success 200 {object} entities.User
// @Failure 400 {object} httpUtil.ResponseError
// @Failure 401 {object} httpUtil.ResponseError
// @Failure 403 {object} httpUtil.ResponseError
// @Router /api/Admin/Account/{id} [get]
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

// @Summary Создание нового пользователя
// @Tags AdminAccountController
// @Description Создание нового пользователя с указанными данными
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param request body authHandler.CreateUser.userData true "User data"
// @Success 201 {object} entities.User
// @Failure 400 {object} httpUtil.ResponseError
// @Failure 401 {object} httpUtil.ResponseError
// @Failure 403 {object} httpUtil.ResponseError
// @Router /api/Admin/Account [post]
func (ah AuthHandlers) CreateUser(ctx *gin.Context) {
	type userData struct {
		Username string  `json:"username" binding:"required"`
		Password string  `json:"password" binding:"required"`
		IsAdmin  bool    `json:"isAdmin"`
		Balance  float64 `json:"balance"`
	}

	var usrData userData
	if err := ctx.BindJSON(&usrData); err != nil {
		ctx.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	user := entities.User{
		Username: usrData.Username,
		Password: usrData.Password,
		IsAdmin:  usrData.IsAdmin,
		Balance:  usrData.Balance,
	}

	user, err := ah.uc.CreateUser(user)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, user)
}

// @Summary Обновление данных пользователя
// @Tags AdminAccountController
// @Description Обновление данных пользователя с id={id}
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param id path uint true "Account id"
// @Param requset body authHandler.UpdateUser.userData true "User data"
// @Success 200 {object} entities.User
// @Failure 400 {object} httpUtil.ResponseError
// @Failure 401 {object} httpUtil.ResponseError
// @Failure 403 {object} httpUtil.ResponseError
// @Router /api/Admin/Account/{id} [put]
func (ah AuthHandlers) UpdateUser(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 0 {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": "invalid value of id param"})
		return
	}
	type userData struct {
		Username string  `json:"username" binding:"required"`
		Password string  `json:"password" binding:"required"`
		IsAdmin  bool    `json:"isAdmin"`
		Balance  float64 `json:"balance"`
	}
	var usrData userData
	if err := ctx.BindJSON(&usrData); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": "invalid data"})
		return
	}

	user := entities.User{
		Id:       uint(id),
		Username: usrData.Username,
		Password: usrData.Password,
		IsAdmin:  usrData.IsAdmin,
		Balance:  usrData.Balance,
	}

	user, err = ah.uc.UpdateUser(user)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, user)
}

// @Summary Обновление данных пользователя
// @Tags AdminAccountController
// @Description Обновление данных пользователя с id={id}
// @Security ApiKeyAuth
// @Param id path uint true "Account id"
// @Success 200 {object} entities.User
// @Failure 400 {object} httpUtil.ResponseError
// @Failure 401 {object} httpUtil.ResponseError
// @Failure 403 {object} httpUtil.ResponseError
// @Router /api/Admin/Account/{id} [delete]
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
