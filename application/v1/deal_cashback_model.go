package v1

import "time"

type DealCashback struct {
	ID                          int        `json:"id"`
	GUID                        string     `json:"guid"`
	UserGUID                    string     `json:"user_guid"`
	ShoppingListGUID            string     `json:"shopping_list_guid"`
	DealGUID                    string     `json:"deal_guid"`
	DealCashbackTransactionGUID *string    `json:"deal_cashback_transaction_guid"`
	Expired                     int        `sql:"-" json:"expired"`
	CreatedAt                   time.Time  `json:"created_at"`
	UpdatedAt                   time.Time  `json:"updated_at"`
	DeletedAt                   *time.Time `json:"deleted_at"`

	Dealcashbackstatus *DealCashbackStatus `json:"deal_cashback_statuses" gorm:"ForeignKey:GUID;AssociationForeignKey:DealCashbackGUID"`

	// Has One Deal Cashback Transaction
	Dealcashbacktransactions *DealCashbackTransaction `json:"deal_cashback_transaction,omitempty" gorm:"ForeignKey:GUID;AssociationForeignKey:DealCashbackTransactionGUID"`

	Deals         *Deal         `json:"deal,omitempty" gorm:"ForeignKey:GUID;AssociationForeignKey:DealGUID"`
	Users         *User         `json:"user,omitempty" gorm:"ForeignKey:GUID;AssociationForeignKey:UserGUID"`
	Shoppinglists *ShoppingList `json:"shopping_list,omitempty" gorm:"ForeignKey:GUID;AssociationForeignKey:ShoppingListGUID"`
}

// type DealCashbackWithoutExpired struct {
// 	ID                          int        `json:"id"`
// 	GUID                        string     `json:"guid"`
// 	UserGUID                    string     `json:"user_guid"`
// 	ShoppingListGUID            string     `json:"shopping_list_guid"`
// 	DealGUID                    string     `json:"deal_guid"`
// 	DealCashbackTransactionGUID *string    `json:"deal_cashback_transaction_guid"`
// 	CreatedAt                   time.Time  `json:"created_at"`
// 	UpdatedAt                   time.Time  `json:"updated_at"`
// 	DeletedAt                   *time.Time `json:"deleted_at"`

// 	Dealcashbackstatus *DealCashbackStatus `json:"deal_cashback_statuses" gorm:"ForeignKey:GUID;AssociationForeignKey:DealCashbackGUID"`

// 	// Has One Deal Cashback Transaction
// 	Dealcashbacktransactions *DealCashbackTransaction `json:"deal_cashback_transaction,omitempty" gorm:"ForeignKey:GUID;AssociationForeignKey:DealCashbackTransactionGUID"`

// 	Deals         *Deal         `json:"deal,omitempty" gorm:"ForeignKey:GUID;AssociationForeignKey:DealGUID"`
// 	Users         *User         `json:"user,omitempty" gorm:"ForeignKey:GUID;AssociationForeignKey:UserGUID"`
// 	Shoppinglists *ShoppingList `json:"shopping_list,omitempty" gorm:"ForeignKey:GUID;AssociationForeignKey:ShoppingListGUID"`
// }

// func (d Deal) TableName() string {
// 	return "deal_cashbacks"
// }
