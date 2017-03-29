package v11

// SmsSend is a request data binding that will be used to bind request body to struct.
// When API receives request with header `application/json`, GIN will used `json` tag to find the data.
// When API receives request with header `multipart/form-data` or `application/x-www-form-urlencoded`,
// GIN will used `form` tag to find the data.
// Used in SMS Handler. See `Send` function.
type SmsSend struct {
	RecipientNo string `form:"recipient_no" json:"recipient_no" binding:"required,numeric,min=11,max=13"`
	Event       string `form:"event" json:"event" binding:"required,alpha"`
}

// SmsVerification is a request data binding that will be used to bind request body to struct.
// When API receives request with header `application/json`, GIN will used `json` tag to find the data.
// When API receives request with header `multipart/form-data` or `application/x-www-form-urlencoded`,
// GIN will used `form` tag to find the data.
// Used in SMS Handler. See `Verify` function.
type SmsVerification struct {
	PhoneNo          string `form:"phone_no" json:"phone_no" binding:"required,numeric,min=11,max=13"`
	NewPhoneNo       string `form:"new_phone_no" json:"new_phone_no" binding:"omitempty,numeric,min=11,max=13"`
	DeviceUUID       string `form:"device_uuid" json:"device_uuid" binding:"required,alphanum"`
	VerificationCode string `form:"verification_code" json:"verification_code" binding:"required,numeric,len=4"`
}
