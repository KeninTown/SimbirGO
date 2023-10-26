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
	//fill transport type [Car, Bike, Scooter]
	var transpotType []models.TransportType
	db.Find(&transpotType)
	if len(transpotType) != 3 {
		db.Create(&models.TransportType{Type: "Car"})
		db.Create(&models.TransportType{Type: "Bike"})
		db.Create(&models.TransportType{Type: "Scooter"})
	}

	log.Println("succesfully migrate database")
	return Database{db: db}, nil
}

// auth repository
func (db Database) FindUserByUsername(username string) models.User {
	var user models.User
	db.db.Find(&user, "username=?", username)
	return user
}

func (db Database) FindUserById(id uint) models.User {
	var user models.User
	db.db.Find(&user, "id=?", id)
	return user
}

func (db Database) CreateUser(user models.User) models.User {
	db.db.Create(&user)
	return user
}

func (db Database) SaveUser(user models.User) {
	db.db.Save(&user)
}

func (db Database) GetUsers(start uint, count int) []models.User {
	var users []models.User
	db.db.Limit(int(count)).Find(&users, "id>=?", start)
	return users
}

func (db Database) DeleteUser(id uint) {
	db.db.Delete(&models.User{}, "id=?", id)
}

// transport repository
func (db Database) FindTypeById(id uint) string {
	var trType models.TransportType
	db.db.Find(&trType, "id=?", id)
	return trType.Type
}

func (db Database) FindTypeByName(typeName string) uint {
	var trType models.TransportType
	db.db.Find(&trType, "type=?", typeName)
	return trType.Id
}

func (db Database) FindTranspot(id uint) models.Transport {
	var transport models.Transport
	db.db.Find(&transport, "id=?", id)
	return transport
}

func (db Database) CreateTransport(transport models.Transport) models.Transport {
	db.db.Create(&transport)
	return transport
}

func (db Database) FindUserTransport(userId, transportId uint) models.Transport {
	var transport models.Transport
	db.db.Where("id = ? AND owner_id = ?", transportId, userId).Find(&transport)
	return transport
}

func (db Database) SaveTransport(transport models.Transport) {
	db.db.Save(&transport)
}

func (db Database) DeleteUserTransport(ownerId, transportId uint) {
	db.db.Where("owner_id = ? AND id = ?", ownerId, transportId).Delete(&models.Transport{})
}
