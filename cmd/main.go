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
