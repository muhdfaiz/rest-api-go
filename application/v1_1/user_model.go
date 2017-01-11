package v1_1

import "time"

// User Model
type User struct {
	ID                int        `json:"id"`
	GUID              string     `json:"guid"`
	FacebookID        *string    `json:"facebook_id"`
	Name              string     `json:"name"`
	Email             string     `json:"email"`
	PhoneNo           string     `json:"phone_no"`
	ProfilePicture    *string    `json:"profile_picture"`
	ReferralCode      string     `json:"referral_code"`
	BankCountry       *string    `json:"bank_country"`
	BankName          *string    `json:"bank_name"`
	BankAccountName   *string    `json:"bank_account_name"`
	BankAccountNumber *string    `json:"bank_account_number"`
	RegisterBy        string     `json:"register_by"`
	Wallet            float64    `json:"available_cashback_amount"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
	DeletedAt         *time.Time `json:"deleted_at"`

	// Virtual Column. Use this column in the response.
	AllTimeAmount *float64 `sql:"-" json:"all_time_amount"`
	PendingAmount *float64 `sql:"-" json:"pending_amount"`

	// User have many Devices
	Devices []*Device `json:"devices,omitempty" gorm:"ForeignKey:UserGUID"`

	// User have many Shopping Lists
	ShoppingLists []*ShoppingList `json:"shopping_lists,omitempty" gorm:"ForeignKey:UserGUID;AssociationForeignKey:GUID"`
}
