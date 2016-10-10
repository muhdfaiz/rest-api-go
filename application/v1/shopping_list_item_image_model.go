package v1

import "time"

type ShoppingListItemImage struct {
	ID                   int        `json:"id"`
	GUID                 string     `json:"guid"`
	UserGUID             string     `json:"user_guid"`
	ShoppingListGUID     string     `json:"shopping_list_guid"`
	ShoppingListItemGUID string     `json:"shopping_list_item_guid"`
	URL                  string     `json:"url"`
	CreatedAt            time.Time  `json:"created_at"`
	UpdatedAt            time.Time  `json:"updated_at"`
	DeletedAt            *time.Time `json:"deleted_at"`

	// ShoppingListItemImage belongs to Shopping List Item
	Items *ShoppingListItem `json:"item,omitempty" gorm:"ForeignKey:ShoppingListItemGUID;AssociationForeignKey:GUID"`
}
