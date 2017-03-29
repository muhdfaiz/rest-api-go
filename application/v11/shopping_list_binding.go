package v11

// CreateShoppingList is a request data binding that will be used to bind request body to struct.
// When API receives request with header `application/json`, GIN will used `json` tag to find the data.
// When API receives request with header `multipart/form-data` or `application/x-www-form-urlencoded`,
// GIN will used `form` tag to find the data.
// Used in Shopping List Handler. See `Create` function.
type CreateShoppingList struct {
	OccasionGUID string `form:"occasion_guid" json:"occasion_guid" binding:"required,uuid5"`
	Name         string `form:"name" json:"name" binding:"required"`
}

// UpdateShoppingList is a request data binding that will be used to bind request body to struct.
// When API receives request with header `application/json`, GIN will used `json` tag to find the data.
// When API receives request with header `multipart/form-data` or `application/x-www-form-urlencoded`,
// GIN will used `form` tag to find the data.
// Used in Shopping List Handler. See `Create` function.
type UpdateShoppingList struct {
	OccasionGUID string `form:"occasion_guid" json:"occasion_guid" binding:"omitempty,uuid5"`
	Name         string `form:"name" json:"name" binding:"omitempty"`
}
