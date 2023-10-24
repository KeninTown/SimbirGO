package entities

type Rent struct {
	Id          uint    `json:"id"`
	TransportId uint    `json:"transportId"`
	UserId      uint    `json:"userId"`
	TimeStart   string  `json:"timeStart"`
	TimeEnd     string  `json:"timeEnd"`
	PriceOfUnit float64 `json:"priceOfUnit"`
	PriceType   string  `json:"priceType"`
	FinalPrice  float64 `json:"finalPrice"`
}
