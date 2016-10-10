package v1

import (
	"time"

	"github.com/shopspring/decimal"
)

// Item Model
type Item struct {
	ID             int             `json:"id"`
	GUID           string          `json:"guid"`
	MasterCategory string          `json:"master_category"`
	Category       string          `json:"category"`
	SubCategory    string          `json:"sub_category"`
	Price          decimal.Decimal `json:"price"`
	Remark         string          `json:"remarks"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
	DeletedAt      *time.Time      `json:"deleted_at"`
}

// TableName function used to set Item table name to be `item``
func (i Item) TableName() string {
	return "item"
}
