package v1

import "time"

type DealCashbackTransaction struct {
	ID               int        `json:"id"`
	GUID             string     `json:"guid"`
	UserGUID         string     `json:"user_guid"`
	TransactionGUID  string     `json:"transaction_guid"`
	ReceiptURL       string     `json:"receipt_url"`
	VerificationDate *string    `json:"verification_date"`
	RemarkTitle      *string    `json:"remark_title"`
	RemarkBody       *string    `json:"remark_body"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
	DeletedAt        *time.Time `json:"deleted_at"`

	// Has many Deal Cashback
	Dealcashbacks []*DealCashback `json:"deal_cashbacks,omitempty" gorm:"ForeignKey:DealCashbackTransactionGUID;AssociationForeignKey:GUID"`

	Transactions *Transaction `json:"transaction,omitempty" gorm:"ForeignKey:TransactionGUID;AssociationForeignKey:GUID"`
}
