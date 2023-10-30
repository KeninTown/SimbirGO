package models

type RentType struct {
	Id   uint   `gorm:"primaryKey"`
	Type string `gorm:"not null"`
}
