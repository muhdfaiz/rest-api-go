package v1

import "time"

// Setting Model
type Setting struct {
	ID        int       `json:"id"`
	GUID      string    `json:"guid"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	Value     string    `json:"value"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
