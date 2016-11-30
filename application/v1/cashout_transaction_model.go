package v1

import "time"

type CashoutTransaction struct {
	ID                    int        `json:"id"`
	GUID                  string     `json:"guid"`
	UserGUID              string     `json:"facebook_id"`
	TransactionGUID       string     `json:"name"`
	BankAccountHolderName string     `json:"bank_account_name"`
	BankAccountNumber     string     `json:"bank_account_number"`
	BankName              string     `json:"bank_name"`
	BankCountry           string     `json:"bank_country"`
	CreatedAt             time.Time  `json:"created_at"`
	UpdatedAt             time.Time  `json:"updated_at"`
	DeletedAt             *time.Time `json:"deleted_at"`

	Transactions *Transaction `json:"transaction,omitempty" gorm:"ForeignKey:TransactionGUID;AssociationForeignKey:GUID"`
}
