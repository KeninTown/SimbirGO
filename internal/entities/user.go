package entities

type User struct {
	Id       uint   `json:"id"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	IsAdmin  bool   `json:"isAdmin"`
}
