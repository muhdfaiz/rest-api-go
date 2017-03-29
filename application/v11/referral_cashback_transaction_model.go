package v11

import "time"

// ReferralCashbackTransaction Model
type ReferralCashbackTransaction struct {
	ID              uint       `json:"id"`
	GUID            string     `json:"guid"`
	UserGUID        string     `json:"user_guid"`
	ReferrerGUID    string     `json:"referrer_guid"`
	TransactionGUID string     `json:"transaction_guid"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	DeletedAt       *time.Time `json:"deleted_at"`

	// Referral Cashback Transaction has one User
	Users *User `json:"user,omitempty" gorm:"ForeignKey:UserGUID;AssociationForeignKey:GUID"`

	// Referral Cashback Transaction has one Referrer
	Referrers *User `json:"referrer,omitempty" gorm:"ForeignKey:ReferrerGUID;AssociationForeignKey:GUID"`
}
