package models

type Transport struct {
	Id            uint `gorm:"primaryKey"`
	OwnerId       uint `gorm:"not null"`
	Owner         User `gorm:"foreignKey:OwnerId"`
	TypeId        uint
	TransportType TransportType `gorm:"foreignKey:TypeId"`
	CanBeRented   bool          `gorm:"not null; type:boolean"`
	Model         string        `gorm:"not null"`
	Color         string        `gorm:"not null"`
	Identifier    string        `gorm:"not null"`
	Description   string        `gorm:"not null"`
	Latitude      float64       `gorm:"not null; type: numeric"`
	Longitude     float64       `gorm:"not null; type: numeric"`
	MinutePrice   float64
	DayPrice      float64
}
