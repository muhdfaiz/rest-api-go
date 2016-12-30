package v1

import "time"

// SmsHistory model
type SmsHistory struct {
	ID               uint       `json:"id"`
	GUID             string     `json:"guid"`
	UserGUID         string     `json:"user_guid"`
	Provider         string     `json:"provider"`
	SmsID            string     `json:"sms_id"`
	Text             string     `json:"text"`
	RecipientNo      string     `json:"recipient_no"`
	VerificationCode string     `json:"verification_code"`
	Status           string     `json:"status"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
	DeletedAt        *time.Time `json:"deleted_at"`
}
