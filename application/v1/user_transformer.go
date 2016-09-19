package v1

import "time"

type UserTransformer struct {
}

type BaseTransformer struct {
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

// UserTransformer used to format response data
type BaseUserTransformer struct {
	ID                int    `json:"id"`
	GUID              string `json:"guid"`
	FacebookID        string `json:"facebook_id"`
	Name              string `json:"name"`
	Email             string `json:"email"`
	PhoneNo           string `json:"phone_no"`
	ProfilePicture    string `json:"profile_picture"`
	ReferralCode      string `json:"referral_code"`
	BankCountry       string `json:"bank_country"`
	BankName          string `json:"bank_name"`
	BankAccountName   string `json:"bank_account_name"`
	BankAccountNumber string `json:"bank_account_number"`
	RegisterBy        string `json:"register_by"`
	BaseTransformer
}

type CreateUserTransformer struct {
	BaseUserTransformer
	BaseTransformer
}

type ReadUserTransformer struct {
	BaseUserTransformer
	Verified               string `json:"verified"`
	TotalCashbackEarned    string `json:"total_cashback_earned"`
	TotalCashbackAvailable string `json:"total_cashback_available"`
	TotalCashbackExpected  string `json:"total_cashback_expected"`
	BaseTransformer
}

func (ut *UserTransformer) TransformCreateData(user *User) interface{} {
	return nil
}
