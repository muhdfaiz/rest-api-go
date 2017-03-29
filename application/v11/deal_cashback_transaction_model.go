package v11

import "time"

// DealCashbackTransaction model
type DealCashbackTransaction struct {
	ID               int        `json:"id"`
	GUID             string     `json:"guid"`
	UserGUID         string     `json:"user_guid"`
	TransactionGUID  string     `json:"transaction_guid"`
	ReceiptURL       string     `json:"receipt_url"`
	VerificationDate *time.Time `json:"verification_date"`
	RemarkTitle      *string    `json:"remark_title"`
	RemarkBody       *string    `json:"remark_body"`
	Status           string     `json:"status"`
	TotalDeal        int        `sql:"-" json:"total_deal"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
	DeletedAt        *time.Time `json:"deleted_at"`

	// Deal Cashback Transaction Has Many Shopping Lists.
	DealCashbackGroupByShoppingListName []*ShoppingList `json:"deal_cashbacks_group_by_shopping_list"`

	// Deal Cashback Trasnaction Has many Deal Cashbacks.
	Dealcashbacks []*DealCashback `json:"deal_cashbacks,omitempty" gorm:"ForeignKey:DealCashbackTransactionGUID;AssociationForeignKey:GUID"`

	// DealCashbackTransaction Has One Receipt.
	Receipt *Receipt `json:"receipt" gorm:"ForeignKey:DealCashbackTransactionGUID;AssociationForeignKey:GUID"`

	// Deal Cashback Transaction Has One Transaction.
	Transactions *Transaction `json:"transaction,omitempty" gorm:"ForeignKey:TransactionGUID;AssociationForeignKey:GUID"`
}
