package v1

import "time"

type ReferralCashback struct {
	ID             uint       `json:"id"`
	GUID           string     `json:"guid"`
	ReferrerGUID   string     `json:"referrer_guid"`
	ReferentGUID   string     `json:"referent_guid"`
	CashbackAmount float32    `json:"cashback_amount"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	DeletedAt      *time.Time `json:"deleted_at"`
}
