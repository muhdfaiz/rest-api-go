package v1

type LoginViaPhone struct {
	PhoneNo string `form:"phone_no" json:"phone_no" binding:"required,numeric"`
}

type LoginViaFacebook struct {
	FacebookID string `form:"facebook_id" json:"facebook_id" binding:"required,numeric"`
	DeviceUUID string `form:"device_uuid" json:"device_uuid" binding:"required,alphanum"`
}
