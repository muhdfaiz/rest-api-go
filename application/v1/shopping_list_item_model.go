package v1

import "time"

type ShoppingListItem struct {
	ID               uint       `json:"id"`
	GUID             string     `json:"guid"`
	UserGUID         string     `json:"user_guid"`
	ShoppingListGUID string     `json:"shopping_list_guid"`
	Name             string     `json:"name"`
	Category         string     `json:"category"`
	SubCategory      string     `json:"sub_category"`
	Quantity         int        `json:"quantity"`
	Remark           string     `json:"remark"`
	AddedFromDeal    int        `json:"added_from_deal"`
	DealGUID         *string    `json:"deal_guid"`
	CashbackAmount   *float64   `json:"cashback_amount"`
	DealExpired      *int       `json:"deal_expired"`
	AddedToCart      int        `json:"added_to_cart"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
	DeletedAt        *time.Time `json:"deleted_at"`

	// Has many Shopping List Item Image
	Images []*ShoppingListItemImage `json:"images,omitempty" gorm:"ForeignKey:ShoppingListItemGUID;AssociationForeignKey:GUID"`

	// Has many Deals
	Deals []*Deal `json:"deals,omitempty"`
}
