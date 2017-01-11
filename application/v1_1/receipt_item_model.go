package v1_1

import "time"

// ReceiptItem Model
type ReceiptItem struct {
	ID          int        `json:"id"`
	GUID        string     `json:"guid"`
	ReceiptGUID string     `json:"receipt_guid"`
	Name        string     `json:"name"`
	Quantity    int        `json:"quantity"`
	Price       float64    `json:"price"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at"`
}

// TableName function used to override default plural table name used by gorm based on struct name.
func (ri ReceiptItem) TableName() string {
	return "receipt_item"
}
