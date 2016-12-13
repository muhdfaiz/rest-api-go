package v1

import "time"

type DefaultShoppingListItem struct {
	ID               int        `json:"id"`
	GUID             string     `json:"guid"`
	ShoppingListGUID string     `json:"shopping_list_guid"`
	Name             string     `json:"name"`
	Category         string     `json:"category"`
	SubCategory      string     `json:"sub_category"`
	Quantity         int        `json:"quantity"`
	Remark           string     `json:"remark"`
	AddedToCart      int        `json:"added_to_cart"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
	DeletedAt        *time.Time `json:"deleted_at"`
}
