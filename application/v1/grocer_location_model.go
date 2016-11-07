package v1

import "time"

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

// TableName function used to set Item table name to be `item``
func (gl GrocerLocation) TableName() string {
	return "grocer_location"
}
