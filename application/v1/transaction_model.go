package v1

import "time"

type Transaction struct {
	ID                    int        `json:"id"`
	GUID                  string     `json:"guid"`
	UserGUID              string     `json:"user_guid"`
	TransactionTypeGUID   string     `json:"transaction_type_guid"`
	TransactionStatusGUID string     `json:"transaction_status_guid"`
	Amount                float32    `json:"amount"`
	CreatedAt             time.Time  `json:"created_at"`
	UpdatedAt             time.Time  `json:"updated_at"`
	DeletedAt             *time.Time `json:"deleted_at"`

	Dealcashbacktransactions *DealCashbackTransaction `json:"deal_cashback_transaction,omitempty" gorm:"ForeignKey:TransactionGUID;AssociationForeignKey:GUID"`

	Transactiontypes *TransactionType `json:"transaction_type" gorm:"ForeignKey:TransactionTypeGUID;AssociationForeignKey:GUID"`

	Transactionstatuses *TransactionStatus `json:"transaction_status" gorm:"ForeignKey:TransactionStatusGUID;AssociationForeignKey:GUID"`
}
