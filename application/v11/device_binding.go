package v11

// CreateDevice is a request data binding that will be used to bind request body to struct.
// When API receives request with header `application/json`, GIN will used `json` tag to find the data.
// When API receives request with header `multipart/form-data` or `application/x-www-form-urlencoded`,
// GIN will used `form` tag to find the data.
// Used in Device Handler. See `Create` function.
type CreateDevice struct {
	UserGUID   string `form:"user_guid" json:"user_guid" binding:"omitempty,uuid5"`
	UUID       string `form:"uuid" json:"uuid" binding:"required,alphanum"`
	Os         string `form:"os" json:"os" binding:"required,alpha"`
	Model      string `form:"model" json:"model" binding:"required"`
	PushToken  string `form:"push_token" json:"push_token" binding:"required"`
	AppVersion string `form:"app_version" json:"app_version" binding:"required"`
}

// UpdateDevice is a request data binding that will be used to bind request body to struct.
// When API receives request with header `application/json`, GIN will used `json` tag to find the data.
// When API receives request with header `multipart/form-data` or `application/x-www-form-urlencoded`,
// GIN will used `form` tag to find the data.
// Used in Deal Cashback Handler. See `Update` function.
type UpdateDevice struct {
	UserGUID   string `form:"user_guid" json:"user_guid" binding:"omitempty,uuid5"`
	Os         string `form:"os" json:"os" binding:"omitempty,alpha"`
	Model      string `form:"model" json:"model" binding:"omitempty"`
	PushToken  string `form:"push_token" json:"push_token" binding:"omitempty"`
	AppVersion string `form:"app_version" json:"app_version" binding:"omitempty"`
}
