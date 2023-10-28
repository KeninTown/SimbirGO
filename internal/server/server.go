package server

import (
	"context"
	"log"
	"net/http"
	"simbirGo/internal/server/handlers/authHandler"
	"simbirGo/internal/server/handlers/paymentHandler"
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

func (s *Server) Run(ctx context.Context, uc authHandler.AuthUsecase, pu paymentHandler.PaymentUsecase, tu transportHandler.TransportUsecase) {
	//swagger route
	s.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	//auth routes
	ah := authHandler.New(uc)

	authRouts := s.router.Group("/", middleware.CheckAuthification())
	authRouts.GET("/api/Account/Me", ah.MyAccount)
	s.router.POST("/api/Account/SignIn", ah.SignIn)
	s.router.POST("/api/Account/SignUp", ah.SignUp)
	authRouts.POST("/api/Account/SignOut", ah.SignOut)
	authRouts.PUT("/api/Account/Update", ah.Update)

	//admin auth routes
	adminAuthRouts := s.router.Group("/api/Admin/Account", middleware.CheckAuthification(), middleware.CheckAdminStatus())
	adminAuthRouts.GET("/", ah.GetUsers)
	adminAuthRouts.GET("/:id", ah.GetUser)
	adminAuthRouts.POST("/", ah.CreateUser)
	adminAuthRouts.PUT("/:id", ah.UpdateUser)
	adminAuthRouts.DELETE("/:id", ah.DeleteUser)

	//payment rout
	ph := paymentHandler.New(pu)
	s.router.POST("/api/Payment/Hesoyam/:id", middleware.CheckAuthification(), ph.IncreaseBalance)

	//transport routes
	th := transportHandler.New(tu)

	s.router.GET("/api/Transport/:id", th.GetTransport)
	transportAuthRoutes := s.router.Group("/api/Transport", middleware.CheckAuthification())
	transportAuthRoutes.POST("/", th.CreateTransport)
	transportAuthRoutes.PUT("/:id", th.UpdateTransport)
	transportAuthRoutes.DELETE("/:id", th.DeleteUserTransport)

	//transport auth routes
	transportAdminRoutes := s.router.Group("/api/Admin/Transport", middleware.CheckAuthification(), middleware.CheckAdminStatus())
	transportAdminRoutes.GET("/", th.AdminGetTransports)
	transportAdminRoutes.GET("/:id", th.AdminGetTransport)
	transportAdminRoutes.POST("/", th.AdminCreateTransport)
	transportAdminRoutes.PUT("/:id", th.AdminUpdateTransport)
	transportAdminRoutes.DELETE("/:id", th.DeleteTransport)

	srv := http.Server{
		Addr:    s.addr,
		Handler: s.router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println("failed to listen")
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
