package v11

import "time"

// ItemSubCategory Model
type ItemSubCategory struct {
	ID        int        `json:"id"`
	GUID      string     `json:"guid"`
	Name      string     `json:"name"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`

	// Virtual Column. Use to include the column in the response.
	TotalDeals int `sql:"-" json:"total_deals,omitempty"`

	// Item Subcategory Has Many Deals
	Deals []*Deal `json:"deals"`
}

// TableName function used to override default plural table name used by gorm based on struct name.
func (isc ItemSubCategory) TableName() string {
	return "subcategory"
}
