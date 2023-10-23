package models

type TransportType struct {
	Id   uint   `gorm:"primaryKey"`
	Type string `gorm:"not null"`
}
