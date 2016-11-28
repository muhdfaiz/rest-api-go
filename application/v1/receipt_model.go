package v1

import "time"

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

	Receiptitems []*ReceiptItem `json:"receipt_items" gorm:"ForeignKey:ReceiptGUID;AssociationForeignKey:GUID"`
}

func (r Receipt) TableName() string {
	return "receipt"
}
