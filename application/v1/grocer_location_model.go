package v1

import "time"

// GrocerLocation model
type GrocerLocation struct {
	ID        int        `json:"id"`
	GUID      string     `json:"guid"`
	GrocerID  int        `json:"grocer_id"`
	Name      string     `json:"name"`
	Lat       float64    `json:"lat"`
	Lng       float64    `json:"lng"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

// TableName function used to override default plural table name used by gorm based on struct name.
func (gl GrocerLocation) TableName() string {
	return "grocer_location"
}
