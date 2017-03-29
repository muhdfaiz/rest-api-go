package v11

import "time"

// Occasion Model
type Occasion struct {
	ID        int        `json:"id"`
	GUID      string     `json:"guid"`
	Slug      string     `json:"slug"`
	Name      string     `json:"name"`
	Image     string     `json:"image"`
	Active    int        `json:"active"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}
