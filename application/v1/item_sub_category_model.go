package v1

import "time"

type ItemSubCategory struct {
	ID        int        `json:"id"`
	GUID      string     `json:"guid"`
	Name      string     `json:"name"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`

	// Have Many Deals
	TotalDeals int     `json:"total_deals"`
	Deals      []*Deal `json:"deals"`
}

// TableName function used to set Item table name to be `item``
func (isc ItemSubCategory) TableName() string {
	return "subcategory"
}
