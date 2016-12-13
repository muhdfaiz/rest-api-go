package v1

import "time"

type ItemCategory struct {
	ID        int        `json:"id"`
	GUID      string     `json:"guid"`
	Img       string     `json:"img"`
	Name      string     `json:"name"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`

	// Have Many Deals
	TotalDeals int     `json:"total_deals,omitempty"`
	Deals      []*Deal `json:"deals"`
}

// TableName function used to set Item table name to be `item``
func (i ItemCategory) TableName() string {
	return "category"
}
