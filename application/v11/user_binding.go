package v11

// CreateUser is a request data binding that will be used to bind request body to struct.
// When API receives request with header `application/json`, GIN will used `json` tag to find the data.
// When API receives request with header `multipart/form-data` or `application/x-www-form-urlencoded`,
// GIN will used `form` tag to find the data.
// Used in User Handler. See `Create` function.
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

// UpdateUser is a request data binding that will be used to bind request body to struct.
// When API receives request with header `application/json`, GIN will used `json` tag to find the data.
// When API receives request with header `multipart/form-data` or `application/x-www-form-urlencoded`,
// GIN will used `form` tag to find the data.
// Used in User Handler. See `Update` function.
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
