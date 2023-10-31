package server

import (
	"context"
	"log"
	"net/http"
	"simbirGo/internal/server/handlers/authHandler"
	"simbirGo/internal/server/handlers/paymentHandler"
	"simbirGo/internal/server/handlers/rentHandler"
	"simbirGo/internal/server/handlers/transportHandler"
	middleware "simbirGo/internal/server/middlewares"
	"time"

	_ "simbirGo/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Usecase interface {
	authHandler.AuthUsecase
	paymentHandler.PaymentUsecase
	transportHandler.TransportUsecase
}

type Server struct {
	addr   string
	router *gin.Engine
}

func New(addr string) Server {
	return Server{
		addr:   addr,
		router: gin.Default(),
	}
}

func (s *Server) Run(ctx context.Context, uc authHandler.AuthUsecase, pu paymentHandler.PaymentUsecase, tu transportHandler.TransportUsecase, ru rentHandler.RentUsecase) {
	//swagger route
	s.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	//auth routes
	ah := authHandler.New(uc)

	//user auth routes
	authRouts := s.router.Group("/", middleware.CheckAuthification())
	authRouts.GET("/api/Account/Me", ah.UserMyAccount)
	s.router.POST("/api/Account/SignIn", ah.UserSignIn)
	s.router.POST("/api/Account/SignUp", ah.UserSignUp)
	authRouts.POST("/api/Account/SignOut", ah.UserSignOut)
	authRouts.PUT("/api/Account/Update", ah.UserUpdate)

	//admin auth routes
	adminAuthRouts := s.router.Group("/api/Admin/Account", middleware.CheckAuthification(),
		middleware.CheckAdminStatus())
	adminAuthRouts.GET("/", ah.AdminGetUsers)
	adminAuthRouts.GET("/:id", ah.AdminGetUser)
	adminAuthRouts.POST("/", ah.AdminCreateUser)
	adminAuthRouts.PUT("/:id", ah.AdminUpdateUser)
	adminAuthRouts.DELETE("/:id", ah.AdminDeleteUser)

	//payment rout
	ph := paymentHandler.New(pu)
	s.router.POST("/api/Payment/Hesoyam/:id", middleware.CheckAuthification(), ph.IncreaseBalance)

	//transport routes
	th := transportHandler.New(tu)

	//user transport routes
	s.router.GET("/api/Transport/:id", th.UserGetTransport)
	transportAuthRoutes := s.router.Group("/api/Transport",
		middleware.CheckAuthification())
	transportAuthRoutes.POST("/", th.UserCreateTransport)
	transportAuthRoutes.PUT("/:id", th.UserUpdateTransport)
	transportAuthRoutes.DELETE("/:id", th.UserDeleteTransport)

	//admin transport routes
	transportAdminRoutes := s.router.Group("/api/Admin/Transport",
		middleware.CheckAuthification(), middleware.CheckAdminStatus())
	transportAdminRoutes.GET("/", th.AdminGetTransports)
	transportAdminRoutes.GET("/:id", th.AdminGetTransport)
	transportAdminRoutes.POST("/", th.AdminCreateTransport)
	transportAdminRoutes.PUT("/:id", th.AdminUpdateTransport)
	transportAdminRoutes.DELETE("/:id", th.AdminDeleteTransport)

	//rent routes
	rh := rentHandler.New(ru)

	//user rent routes
	s.router.GET("/api/Rent/Transport", rh.GetAvalibleTransport)
	rentRouts := s.router.Group("/api/Rent", middleware.CheckAuthification())
	rentRouts.GET("/:id", rh.UserGetRent)
	rentRouts.GET("/MyHistory", rh.UserGetHistory)
	rentRouts.GET("/TransportHistory/:id", rh.UserGetTransportHistory)
	rentRouts.POST("/New/:id", rh.UserCreateNewRent)
	rentRouts.POST("/End/:id", rh.UserEndRent)

	//admin rent routes
	rentsAdminRoutes := s.router.Group("/api/Admin", middleware.CheckAuthification(),
		middleware.CheckAdminStatus())
	rentsAdminRoutes.GET("/Rent/:id", rh.AdminGetRent)
	rentsAdminRoutes.POST("/Rent", rh.AdminCreateRent)
	rentsAdminRoutes.POST("/Rent/End/:id", rh.AdminEndRent)
	rentsAdminRoutes.GET("/UserHistory/:id", rh.AdminGetUserHistory)
	rentsAdminRoutes.GET("/TransportHistory/:id", rh.AdminGetTransportHistory)
	rentsAdminRoutes.PUT("/Rent/:id", rh.AdminUpdateRent)
	rentsAdminRoutes.DELETE("/Rent/:id", rh.AdminDeleteRent)

	srv := http.Server{
		Addr:    s.addr,
		Handler: s.router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println("failed to listen server")
		}
	}()

	//gracefull shutdown
	<-ctx.Done()
	log.Println("closing server gracefully...")
	ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	if err := srv.Shutdown(ctxTimeout); err != nil {
		log.Println("failed to shutdown server gracefully")
	}
	log.Println("server closed gracefully")
}
