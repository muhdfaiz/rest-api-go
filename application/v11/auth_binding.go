package v11

// LoginViaPhone is a request data binding that will be used to bind request body to struct.
// When API receives request with header `application/json`, GIN will used `json` tag to find the data.
// When API receives request with header `multipart/form-data` or `application/x-www-form-urlencoded`,
// GIN will used `form` tag to find the data.
// Used in Auth Handler. See `LoginViaPhone` function.
type LoginViaPhone struct {
	PhoneNo string `form:"phone_no" json:"phone_no" binding:"required,numeric"`
}

// LoginViaFacebook is a request data binding that will be used to bind request body to struct.
// When API receives request with header `application/json`, GIN will used `json` tag to find the data.
// When API receives request with header `multipart/form-data` or `application/x-www-form-urlencoded`,
// GIN will used `form` tag to find the data.
// Used in Auth Handler. See `LoginViaFacebook` function.
type LoginViaFacebook struct {
	FacebookID string `form:"facebook_id" json:"facebook_id" binding:"required,numeric"`
	DeviceUUID string `form:"device_uuid" json:"device_uuid" binding:"required,alphanum"`
}
