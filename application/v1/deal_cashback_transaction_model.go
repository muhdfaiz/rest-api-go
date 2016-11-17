package v1

import "time"

type DealCashbackTransaction struct {
	ID               int        `json:"id"`
	GUID             string     `json:"guid"`
	UserGUID         string     `json:"user_guid"`
	TransactionGUID  string     `json:"transaction_guid"`
	ReferenceID      string     `json:"reference_id"`
	ReceiptURL       string     `json:"receipt_url"`
	VerificationDate *string    `json:"verification_date"`
	Remark           *string    `json:"remark"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
	DeletedAt        *time.Time `json:"deleted_at"`

	// Has many Deal Cashback
	Dealcashbacks []*DealCashback `json:"deal_cashbacks,omitempty" gorm:"ForeignKey:DealCashbackTransactionGUID;AssociationForeignKey:GUID"`
}
