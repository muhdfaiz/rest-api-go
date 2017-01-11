package v1

type CreateUser struct {
	FacebookID       string `form:"facebook_id" json:"facebook_id" binding:"omitempty,numeric"`
	Name             string `form:"name" json:"name" binding:"required"`
	Email            string `form:"email" json:"email" binding:"required,email"`
	PhoneNo          string `form:"phone_no" json:"phone_no" binding:"required,numeric,min=11,max=13"`
	ProfilePicture   string `form:"profile_picture" json:"profile_picture" binding:"omitempty"`
	ReferralCode     string `form:"referral_code" json:"referral_code" binding:"omitempty,alphanum,max=8"`
	VerificationCode string `form:"verification_code" json:"verification_code" binding:"required,numeric,len=4"`
	DeviceUUID       string `form:"device_uuid" json:"device_uuid" binding:"required,alphanum"`
	Debug            int    `form:"debug" json:"debug" binding:"omitempty"`
}

type UpdateUser struct {
	Name              string `form:"name" json:"name" binding:"omitempty"`
	Email             string `form:"email" json:"email" binding:"omitempty,email"`
	PhoneNo           string `form:"phone_no" json:"phone_no" binding:"omitempty,numeric,min=11,max=13"`
	ProfilePicture    string `form:"profile_picture" json:"profile_picture" binding:"omitempty"`
	BankCountry       string `form:"bank_country" json:"bank_country" binding:"omitempty"`
	BankName          string `form:"bank_name" json:"bank_name" binding:"omitempty"`
	BankAccountName   string `form:"bank_account_name" json:"bank_account_name" binding:"omitempty"`
	BankAccountNumber string `form:"bank_account_number" json:"bank_account_number" binding:"omitempty,numeric"`
}
