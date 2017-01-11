package v1_1

import "time"

// Receipt Model
type Receipt struct {
	ID                          int        `json:"id"`
	GUID                        string     `json:"guid"`
	DealCashbackTransactionGUID string     `json:"deal_cashback_transaction_guid"`
	Name                        string     `json:"name"`
	Outlet                      string     `json:"outlet"`
	Cashier                     string     `json:"cashier"`
	ReceiptNo                   string     `json:"receipt_no"`
	Date                        string     `json:"date"`
	Time                        string     `json:"time"`
	CreatedAt                   time.Time  `json:"created_at"`
	UpdatedAt                   time.Time  `json:"updated_at"`
	DeletedAt                   *time.Time `json:"deleted_at"`

	// Receipt Has Many Receipt Items.
	Receiptitems []*ReceiptItem `json:"receipt_items" gorm:"ForeignKey:ReceiptGUID;AssociationForeignKey:GUID"`
}

// TableName function used to override default plural table name used by gorm based on struct name.
func (r Receipt) TableName() string {
	return "receipt"
}
