package models

type Rent struct {
	Id          uint      `json:"id" gorm:"primaryKey"`
	TransportId uint      `json:"transportId"`
	Transport   Transport `gorm:"foreignKey:TransportId"`
	UserId      uint      `json:"userId"`
	User        User      `gorm:"foreignKey:UserId; not null"`
	TimeStart   string    `json:"timeStart" gorm:"not null"`
	TimeEnd     string    `json:"timeEnd"`
	PriceOfUnit float64   `json:"priceOfUnit" gorm:"not null"`
	PriceType   string    `json:"priceType" gorm:"not null"`
	FinalPrice  float64   `json:"finalPrice"`
}
