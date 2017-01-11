package v1

import "time"

// DefaultShoppingListItemImage Model
type DefaultShoppingListItemImage struct {
	ID                   int        `json:"id"`
	GUID                 string     `json:"guid"`
	ShoppingListGUID     string     `json:"shopping_list_guid"`
	ShoppingListItemGUID string     `json:"shopping_list_item_guid"`
	URL                  string     `json:"url"`
	CreatedAt            time.Time  `json:"created_at"`
	UpdatedAt            time.Time  `json:"updated_at"`
	DeletedAt            *time.Time `json:"deleted_at"`
}
