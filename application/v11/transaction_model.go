package v11

import "time"

// Transaction model
type Transaction struct {
	ID                    int        `json:"id"`
	GUID                  string     `json:"guid"`
	UserGUID              string     `json:"user_guid"`
	ReferenceID           string     `json:"reference_id"`
	TransactionTypeGUID   string     `json:"transaction_type_guid"`
	TransactionStatusGUID string     `json:"transaction_status_guid"`
	ReadStatus            int        `json:"read_status"`
	TotalAmount           float64    `json:"total_amount"`
	ApprovedAmount        *float64   `json:"approved_amount,omitempty"`
	RejectedAmount        *float64   `json:"rejected_amount,omitempty"`
	CreatedAt             time.Time  `json:"created_at"`
	UpdatedAt             time.Time  `json:"updated_at"`
	DeletedAt             *time.Time `json:"deleted_at"`

	// Transaction has one User.
	Users *User `json:"user" gorm:"ForeignKey:GUID;AssociationForeignKey:UserGUID"`

	// Transaction has one Deal Cashback Transaction.
	Dealcashbacktransactions *DealCashbackTransaction `json:"deal_cashback_transaction,omitempty" gorm:"ForeignKey:TransactionGUID;AssociationForeignKey:GUID"`

	// Transaction has one Cashout Transaction.
	Cashouttransactions *CashoutTransaction `json:"cashout_transaction,omitempty" gorm:"ForeignKey:TransactionGUID;AssociationForeignKey:GUID"`

	ReferralCashbackTransactions *ReferralCashbackTransaction `json:"referral_cashback_transaction,omitempty" gorm:"ForeignKey:TransactionGUID;AssociationForeignKey:GUID"`

	// Transaction has one Transaction Type.
	Transactiontypes *TransactionType `json:"transaction_type" gorm:"ForeignKey:TransactionTypeGUID;AssociationForeignKey:GUID"`

	// Transaction has one Transaction Status.
	Transactionstatuses *TransactionStatus `json:"transaction_status" gorm:"ForeignKey:TransactionStatusGUID;AssociationForeignKey:GUID"`
}
