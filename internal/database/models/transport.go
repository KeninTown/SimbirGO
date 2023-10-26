package models

type Transport struct {
	Id            uint `json:"id" gorm:"primaryKey"`
	OwnerId       uint `json:"owner_id"`
	Owner         User `gorm:"foreignKey:OwnerId"`
	TypeId        uint
	TransportType TransportType `gorm:"foreignKey:TypeId"`
	CanBeRented   bool          `json:"canBeRented" gorm:"not null; type:boolean"`
	Model         string        `json:"model" gorm:"not null"`
	Color         string        `json:"color" gorm:"not null"`
	Identifier    string        `json:"identifier" gorm:"not null"`
	Description   string        `json:"description" gorm:"not null"`
	Latitude      float64       `json:"latitude" gorm:"not null; type: numeric"`
	Longitude     float64       `json:"longitude" gorm:"not null; type: numeric"`
	MinutePrice   float64       `json:"minutePrice"`
	DayPrice      float64       `json:"dayPrice"`
}
