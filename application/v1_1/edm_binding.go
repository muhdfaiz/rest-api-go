package v1_1

type SendEdmInsufficientFunds struct {
	Name      string `form:"name" json:"name" binding:"required"`
	Email     string `form:"email" json:"email" binding:"required,email"`
	Latitude  string `form:"latitude" json:"latitude" binding:"required,latitude"`
	Longitude string `form:"longitude" json:"longitude" binding:"required,longitude"`
}
