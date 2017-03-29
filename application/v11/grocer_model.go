package v11

import "time"

// Grocer Model
type Grocer struct {
	ID        int        `json:"id"`
	GUID      string     `json:"guid"`
	Name      string     `json:"name"`
	Email     string     `json:"email"`
	Img       string     `json:"img"`
	Status    string     `json:"status"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`

	// Grocer Has Many Grocer Location.
	GrocerLocations []*GrocerLocation `json:"grocer_locations,omitempty" gorm:"ForeignKey:GrocerID;AssociationForeignKey:ID"`

	// Grocer Has Many Deals.
	Deals []*Deal `json:"deals,omitempty"`

	// Virtual Column. Not exist in real database table.
	TotalDeals int `json:"total_deals,omitempty"`
}

// TableName function used to override default plural table name used by gorm based on struct name.
func (g Grocer) TableName() string {
	return "grocer"
}
