package v1

import "time"

type ReceiptItem struct {
	ID          int        `json:"id"`
	GUID        string     `json:"guid"`
	ReceiptGUID string     `json:"receipt_guid"`
	Name        string     `json:"name"`
	Quantity    int        `json:"quantity"`
	Price       float32    `json:"price"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at"`
}

func (ri ReceiptItem) TableName() string {
	return "receipt_item"
}
