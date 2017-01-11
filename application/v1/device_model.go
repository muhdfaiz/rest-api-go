package v1

import "time"

// Device model
type Device struct {
	ID           int        `json:"id"`
	GUID         string     `json:"guid"`
	UserGUID     *string    `json:"user_guid"`
	UUID         string     `json:"uuid"`
	Os           string     `json:"os"`
	Model        string     `json:"model"`
	PushToken    string     `json:"push_token"`
	AppVersion   string     `json:"app_version"`
	TokenExpired int        `json:"token_expired"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at"`
}
