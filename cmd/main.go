package main

import (
	"fmt"
	"log"
	"simbirGo/internal/config"
	"simbirGo/internal/database"
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
}
