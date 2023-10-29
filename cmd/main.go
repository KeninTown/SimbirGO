package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"simbirGo/internal/config"
	"simbirGo/internal/database"
	"simbirGo/internal/server"
	"simbirGo/internal/tokens"
	"simbirGo/internal/usecase/authUsecase"
	"simbirGo/internal/usecase/paymentUsecase"
	"simbirGo/internal/usecase/rentUsecase"
	transportusecase "simbirGo/internal/usecase/transportUsecase"
	"syscall"
)

// @title           SimbirGO REST API
// @version         1.0
// @description     Server for transport booking
// @termsOfService  http://swagger.io/terms/

// @contact.name   Alexander Soldatov
// @contact.email  soldatovalex207z@gmail.com

// @host      localhost:80
// @BasePath  /

// @securitydefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/

func main() {
	cfgPath := "./config/config.yml"
	cfg, err := config.Init(cfgPath)
	if err != nil {
		log.Fatal(err.Error())
	}

	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Print("succesfully connect to database")

	tokens.InitBlackList()

	authUc := authUsecase.New(db)
	paymentUc := paymentUsecase.New(db)
	transportUc := transportusecase.New(db)
	rentUc := rentUsecase.New(db)
	srv := server.New(":80")
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, os.Interrupt)
	defer stop()

	srv.Run(ctx, authUc, paymentUc, transportUc, rentUc)
}
