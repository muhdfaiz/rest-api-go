package v1_1

// CreateDevice will bind request data based on header content type
type CreateDevice struct {
	UserGUID   string `form:"user_guid" json:"user_guid" binding:"omitempty,uuid5"`
	UUID       string `form:"uuid" json:"uuid" binding:"required,alphanum"`
	Os         string `form:"os" json:"os" binding:"required,alpha"`
	Model      string `form:"model" json:"model" binding:"required"`
	PushToken  string `form:"push_token" json:"push_token" binding:"required"`
	AppVersion string `form:"app_version" json:"app_version" binding:"required"`
}

// UpdateDevice will bind request data to JSON during update device
type UpdateDevice struct {
	UserGUID   string `form:"user_guid" json:"user_guid" binding:"omitempty,uuid5"`
	Os         string `form:"os" json:"os" binding:"omitempty,alpha"`
	Model      string `form:"model" json:"model" binding:"omitempty"`
	PushToken  string `form:"push_token" json:"push_token" binding:"omitempty"`
	AppVersion string `form:"app_version" json:"app_version" binding:"omitempty"`
}
