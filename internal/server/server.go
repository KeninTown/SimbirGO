package server

import (
	"context"
	"log"
	"net/http"
	auth "simbirGo/internal/server/handlers/authHandler"
	middleware "simbirGo/internal/server/middlewares"
	"time"

	_ "simbirGo/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Usecase interface {
	auth.AuthUsecase
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

func (s *Server) Run(ctx context.Context, uc Usecase) {
	s.router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"hello": "world",
		})

		ctx.JSON(201, gin.H{"pivo": "ochen"})
	})

	// s.router.GET("/api/Account/Me", handlers.MyAccount(uc))
	ah := auth.New(uc)

	//auth routes
	gr := s.router.Group("/", middleware.CheckAuthification())
	gr.GET("/api/Account/Me", ah.MyAccount)
	s.router.POST("/api/Account/SignIn", ah.SignIn)
	s.router.POST("/api/Account/SignUp", ah.SignUp)
	gr.POST("/api/Account/SignOut", ah.SignOut)
	gr.PUT("/api/Account/Update", ah.Update)

	s.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	//admin auth routes
	adminAuthRouts := s.router.Group("/api/Admin/Account", middleware.CheckAuthification(), middleware.CheckAdminStatus())
	adminAuthRouts.GET("/", ah.GetUsers)
	adminAuthRouts.GET("/:id", ah.GetUser)
	adminAuthRouts.POST("/", ah.CreateUser)
	adminAuthRouts.PUT("/:id", ah.UpdateUser)
	adminAuthRouts.DELETE("/:id", ah.DeleteUser)

	server := http.Server{
		Addr:    s.addr,
		Handler: s.router,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Println("failed to listen")
		}
	}()

	<-ctx.Done()
	log.Println("closing server gracefully...")
	ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	if err := server.Shutdown(ctxTimeout); err != nil {
		log.Println("failed to shutdown server gracefully")
	}
	log.Println("server closed gracefully")
}
