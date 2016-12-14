package v1

import "time"

type DefaultShoppingList struct {
	ID           int        `json:"id"`
	GUID         string     `json:"guid"`
	OccasionGUID string     `json:"occasion_guid"`
	Name         string     `json:"name"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at"`

	// DefaultShoppingList has one Occasion
	Occasions *Occasion `json:"occasions,omitempty" gorm:"ForeignKey:OccasionGUID;AssociationForeignKey:GUID"`

	// Has many Shopping List Item
	Items []*DefaultShoppingListItem `json:"items,omitempty" gorm:"ForeignKey:ShoppingListGUID;AssociationForeignKey:GUID"`
}
