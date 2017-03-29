package v11

// SendEdmInsufficientFunds is a request data binding that will be used to bind request body to struct.
// When API receives request with header `application/json`, GIN will used `json` tag to find the data.
// When API receives request with header `multipart/form-data` or `application/x-www-form-urlencoded`,
// GIN will used `form` tag to find the data.
// Used in EDM Handler. See `InsufficientFunds` function.
type SendEdmInsufficientFunds struct {
	Name      string `form:"name" json:"name" binding:"required"`
	Email     string `form:"email" json:"email" binding:"required,email"`
	Latitude  string `form:"latitude" json:"latitude" binding:"omitempty,latitude"`
	Longitude string `form:"longitude" json:"longitude" binding:"omitempty,longitude"`
}
