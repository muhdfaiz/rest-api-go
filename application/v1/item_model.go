package v1

import "time"

// Item Model
type Item struct {
	ID          int        `json:"id"`
	GUID        string     `json:"guid"`
	GenericID   int        `json:"generic_id"`
	Category    string     `json:"category"`
	SubCategory string     `json:"sub_category"`
	Name        string     `json:"name"`
	Remark      string     `json:"remarks"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at"`
}

// TableName function used to set Item table name to be `item``
func (i Item) TableName() string {
	return "item"
}
