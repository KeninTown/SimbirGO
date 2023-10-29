package models

import (
	"time"
)

type Rent struct {
	Id          uint      `gorm:"primaryKey"`
	TransportId uint      `gorm:"not null"`
	Transport   Transport `gorm:"foreignKey:TransportId"`
	UserId      uint      `gorm:"not null"`
	User        User      `gorm:"foreignKey:UserId; not null"`
	TimeStart   time.Time `gorm:"not null; type: timestamptz"`
	TimeEnd     *time.Time `gorm:"default:null"`
	PriceOfUnit float64   `gorm:"not null"`
	PriceType   string    `gorm:"not null"`
	FinalPrice  float64   `gorm:"default:null"`
}
