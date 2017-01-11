package v1_1

import "time"

// DefaultShoppingListItem Model
type DefaultShoppingListItem struct {
	ID               int        `json:"id"`
	GUID             string     `json:"guid"`
	ShoppingListGUID string     `json:"shopping_list_guid"`
	Name             string     `json:"name"`
	Category         string     `json:"category"`
	Subcategory      string     `json:"sub_category"`
	Quantity         int        `json:"quantity"`
	Remark           string     `json:"remark"`
	AddedToCart      int        `json:"added_to_cart"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
	DeletedAt        *time.Time `json:"deleted_at"`

	// Default Shopping List Item Has many Default Shopping List Item Image
	Images []*DefaultShoppingListItemImage `json:"images,omitempty" gorm:"ForeignKey:ShoppingListItemGUID;AssociationForeignKey:GUID"`

	// Default Shopping List Item Has many Deals
	Deals []*Deal `json:"deals"`
}
