package v1

import "time"

// User Mapping
type User struct {
	ID                int        `json:"id"`
	GUID              string     `json:"guid"`
	FacebookID        string     `json:"facebook_id"`
	Name              string     `json:"name"`
	Email             string     `json:"email"`
	PhoneNo           string     `json:"phone_no"`
	ProfilePicture    string     `json:"profile_picture"`
	ReferralCode      string     `json:"referral_code"`
	BankCountry       string     `json:"bank_country"`
	BankName          string     `json:"bank_name"`
	BankAccountName   string     `json:"bank_account_name"`
	BankAccountNumber string     `json:"bank_account_number"`
	RegisterBy        string     `json:"register_by"`
	Verified          int        `json:"verified"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
	DeletedAt         *time.Time `json:"deleted_at"`

	// Model Relationship
	// User have many Device
	// Devices []device.Device `gorm:"ForeignKey:UserGUID"`
}
