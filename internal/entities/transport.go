package entities

type Transport struct {
	Id            uint    `json:"id" gorm:"primaryKey"`
	OwnerId       uint    `json:"owner_id"`
	TransportType string  `json:"transportType"`
	CanBeRanted   bool    `json:"canBeRented"`
	Model         string  `json:"model"`
	Color         string  `json:"color"`
	Identifier    string  `json:"identifier"`
	Description   string  `json:"description"`
	Latitude      float64 `json:"latitude"`
	Longitude     float64 `json:"longitude"`
	MinutePrice   float64 `json:"minutePrice"`
	DayPrice      float64 `json:"dayPrice"`
}
