package v1_1

import "time"

// SmsHistory model
type SmsHistory struct {
	ID               uint       `json:"id"`
	GUID             string     `json:"guid"`
	Provider         string     `json:"provider"`
	SmsID            string     `json:"sms_id"`
	Text             string     `json:"text"`
	RecipientNo      string     `json:"recipient_no"`
	VerificationCode string     `json:"verification_code"`
	Event            string     `json:"event"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
	DeletedAt        *time.Time `json:"deleted_at"`
}
