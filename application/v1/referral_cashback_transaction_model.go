package v1

import "time"

type ReferralCashbackTransaction struct {
	ID              uint       `json:"id"`
	GUID            string     `json:"guid"`
	UserGUID        string     `json:"user_guid"`
	ReferrerGUID    string     `json:"referrer_guid"`
	TransactionGUID string     `json:"transaction_guid"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	DeletedAt       *time.Time `json:"deleted_at"`
}
