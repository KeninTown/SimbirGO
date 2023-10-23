package database

import (
	"fmt"
	"log"
	"simbirGo/internal/config"
	"simbirGo/internal/database/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database struct {
	db *gorm.DB
}

func Connect(cfg *config.Config) (Database, error) {
	op := "database.Connect()"
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s ",
		cfg.Host, cfg.User, cfg.Password, cfg.DBName, cfg.Port, cfg.SSLMode)

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: dsn,
	}), &gorm.Config{})

	if err != nil {
		return Database{}, fmt.Errorf("%s: failed to connect to postgres: %w", op, err)
	}

	if err := db.AutoMigrate(&models.Rent{}, &models.User{}, &models.Transport{}, models.TransportType{}); err != nil {
		return Database{}, fmt.Errorf("%s: failed to migrate database: %w", op, err)
	}
	log.Println("succesfully migrate database")
	return Database{db: db}, nil
}
