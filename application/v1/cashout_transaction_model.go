package v1

import "time"

// CashoutTransaction model
type CashoutTransaction struct {
	ID                    int        `json:"id"`
	GUID                  string     `json:"guid"`
	UserGUID              string     `json:"user_guid"`
	TransactionGUID       string     `json:"transaction_guid"`
	BankAccountHolderName string     `json:"bank_account_name"`
	BankAccountNumber     string     `json:"bank_account_number"`
	BankName              string     `json:"bank_name"`
	BankCountry           string     `json:"bank_country"`
	RemarkTitle           *string    `json:"remark_title"`
	RemarkBody            *string    `json:"remark_body"`
	TransferDate          *time.Time `json:"transfer_date"`
	ReceiptImage          *string    `json:"receipt_image"`
	CreatedAt             time.Time  `json:"created_at"`
	UpdatedAt             time.Time  `json:"updated_at"`
	DeletedAt             *time.Time `json:"deleted_at"`

	// Cashout Transaction has one Transaction.
	Transactions *Transaction `json:"transaction,omitempty" gorm:"ForeignKey:TransactionGUID;AssociationForeignKey:GUID"`
}
