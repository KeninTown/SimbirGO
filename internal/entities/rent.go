package entities

import (
	"time"
)

type Rent struct {
	Id          uint       `json:"id"`
	TransportId uint       `json:"transportId"`
	UserId      uint       `json:"userId"`
	TimeStart   time.Time  `json:"timeStart"`
	TimeEnd     *time.Time `json:"timeEnd"`
	PriceOfUnit float64    `json:"priceOfUnit"`
	PriceType   string     `json:"priceType"`
	FinalPrice  float64    `json:"finalPrice"`
}
