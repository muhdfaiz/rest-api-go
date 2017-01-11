package v1_1

import "time"

// ItemCategory Model
type ItemCategory struct {
	ID        int        `json:"id"`
	GUID      string     `json:"guid"`
	Img       string     `json:"img"`
	Name      string     `json:"name"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`

	// Virtual Column. Use to include the column in the response.
	TotalDeals int `json:"total_deals,omitempty"`

	// Item Category Has Many Deals
	Deals []*Deal `json:"deals"`
}

// TableName function used to override default plural table name used by gorm based on struct name.
func (i ItemCategory) TableName() string {
	return "category"
}
