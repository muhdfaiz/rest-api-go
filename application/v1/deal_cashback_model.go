package v1

import "time"

// DealCashback Model
type DealCashback struct {
	ID                          int        `json:"id"`
	GUID                        string     `json:"guid"`
	UserGUID                    string     `json:"user_guid"`
	ShoppingListGUID            string     `json:"shopping_list_guid"`
	DealGUID                    string     `json:"deal_guid"`
	DealCashbackTransactionGUID *string    `json:"deal_cashback_transaction_guid"`
	CreatedAt                   time.Time  `json:"created_at"`
	UpdatedAt                   time.Time  `json:"updated_at"`
	DeletedAt                   *time.Time `json:"deleted_at"`

	// Virtual Columns
	Expired               int `sql:"-" json:"expired"`
	RemainingDaysToRemove int `sql:"-" json:"remaining_days_to_remove"`

	// Deal Cashback Has One Deal Cashback Status.
	Dealcashbackstatus *DealCashbackStatus `json:"deal_cashback_statuses" gorm:"ForeignKey:GUID;AssociationForeignKey:DealCashbackGUID"`

	// Deal Cashback Has One Deal Cashback Transaction.
	Dealcashbacktransactions *DealCashbackTransaction `json:"deal_cashback_transaction,omitempty" gorm:"ForeignKey:GUID;AssociationForeignKey:DealCashbackTransactionGUID"`

	// Deal Cashback Has One Deal.
	Deals *Deal `json:"deal,omitempty" gorm:"ForeignKey:GUID;AssociationForeignKey:DealGUID"`

	// Deal Cashback Has One User.
	Users *User `json:"user,omitempty" gorm:"ForeignKey:GUID;AssociationForeignKey:UserGUID"`

	// Deal Cashback Has One Shopping List.
	Shoppinglists *ShoppingList `json:"shopping_list,omitempty" gorm:"ForeignKey:GUID;AssociationForeignKey:ShoppingListGUID"`
}
