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

func (db Database) FindTranspots(start, count int, transportId uint) []models.Transport {
	var transports []models.Transport
	db.db.Where("id >= ? AND type_id = ?", start, transportId).Limit(int(count)).Find(&transports)
	return transports
}

func (db Database) DeleteTransport(id uint) {
	db.db.Delete(&models.Transport{}, "id=?", id)
}

// rent repository
func (db Database) FindAvalibleTransports(lat, long, radius float64, typeId uint) []models.Transport {
	var query string
	if typeId == 0 {
		query = fmt.Sprintf("SELECT * FROM transports WHERE SQRT(power((latitude - %f), 2) + power((longitude - %f), 2)) <= %f AND can_be_rented = true", lat, long, radius)
	} else {
		query = fmt.Sprintf("SELECT * FROM transports WHERE SQRT(power((latitude - %f), 2) + power((longitude - %f), 2)) <= %f AND type_id = %d AND can_be_rented = true", lat, long, radius, typeId)
	}

	rows, err := db.db.Raw(query).Rows()
	if err != nil {
		fmt.Println("err: ", err.Error())
		return []models.Transport{}
	}
	defer rows.Close()
	var transports []models.Transport

	for rows.Next() {
		var (
			id          uint
			ownerId     uint
			typeId      uint
			canBeRented bool
			model       string
			color       string
			identifier  string
			description string
			latitude    float64
			longitude   float64
			minutePrice float64
			dayPrice    float64
		)

		err := rows.Scan(&id, &ownerId, &typeId, &model, &color, &identifier, &description,
			&latitude, &longitude, &minutePrice, &dayPrice, &canBeRented)
		if err != nil {
			log.Println("database.FindAvalibleTransports(): ", err.Error())
		}

		transports = append(transports, models.Transport{
			Id:          id,
			OwnerId:     ownerId,
			TypeId:      typeId,
			CanBeRented: canBeRented,
			Model:       model,
			Color:       color,
			Identifier:  identifier,
			Description: description,
			Latitude:    latitude,
			Longitude:   longitude,
			MinutePrice: minutePrice,
			DayPrice:    dayPrice,
		})
	}
	fmt.Println(transports)
	return transports
}

func (db Database) FindRentById(id int) models.Rent {
	var rent models.Rent
	db.db.Find(&rent, "id = ?", id)
	return rent
}

func (db Database) CreateRent(rent models.Rent) models.Rent {
	db.db.Create(&rent)
	return rent
}
