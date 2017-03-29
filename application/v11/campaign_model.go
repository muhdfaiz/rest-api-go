package v11

import "time"

// Campaign Model
type Campaign struct {
	ID           int        `json:"id"`
	GUID         string     `json:"guid"`
	AdvertiserID int        `json:"advertiser_id"`
	Description  string     `json:"description"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at"`
}
