package v1

type SmsSend struct {
	RecipientNo string `form:"recipient_no" json:"recipient_no" binding:"required,numeric,min=11,max=13"`
	Event       string `form:"event" json:"event" binding:"required,alpha"`
}

type SmsVerification struct {
	PhoneNo          string `form:"phone_no" json:"phone_no" binding:"required,numeric,min=11,max=13"`
	DeviceUUID       string `form:"device_uuid" json:"device_uuid" binding:"required,alphanum"`
	VerificationCode string `form:"verification_code" json:"verification_code" binding:"required,numeric,len=4"`
}
