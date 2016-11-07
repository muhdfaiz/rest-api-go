package v1

import "time"

type DealCashbackTransaction struct {
	ID               int        `json:"id"`
	GUID             string     `json:"guid"`
	UserGUID         string     `json:"user_guid"`
	ReceiptID        string     `json:"receipt_id"`
	ReceiptImage     string     `json:"receipt_image"`
	VerificationDate string     `json:"verification_date"`
	Remark           string     `json:"remark"`
	Status           string     `json:"status"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
	DeletedAt        *time.Time `json:"deleted_at"`

	// Has many Deal Cashback
	DealCashbacks []*DealCashback `json:"deal_cashbacks,omitempty" gorm:"ForeignKey:CashbackTransactionGUID;AssociationForeignKey:GUID"`
}
