package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"simbirGo/internal/config"
	"simbirGo/internal/database"
	"simbirGo/internal/server"
	"simbirGo/internal/usecase"
	"syscall"
)

// @title           Simbir.Go REST API
// @version         1.0
// @description     Server for transport booking
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @host      localhost:80
// @BasePath  /api/v1

// @securityDefinitions.basic  BasicAuth

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
	fmt.Println(db)
	uc := usecase.New(db)
	srv := server.New(":80")
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, os.Interrupt)
	defer stop()

	srv.Run(ctx, uc)
}
