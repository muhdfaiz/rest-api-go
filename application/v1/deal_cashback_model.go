package v1

import "time"

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

	// Has One Deal Cashback Transaction
	DealCashbackTransaction []*DealCashbackTransaction `json:"deal_cashback_transaction,omitempty" gorm:"ForeignKey:GUID;AssociationForeignKey:DealCashbackTransactionGUID"`
}
