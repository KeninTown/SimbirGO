package server

import (
	"context"
	"log"
	"net/http"
	"simbirGo/internal/server/handlers"
	"time"

	_ "simbirGo/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Usecase interface {
	handlers.AuthUsecase
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
	h := handlers.New(uc)
	s.router.POST("/api/Account/SignUp", h.SignUp)

	s.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

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
