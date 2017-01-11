package v1_1

import "time"

// DefaultShoppingList Model
type DefaultShoppingList struct {
	ID           int        `json:"id"`
	GUID         string     `json:"guid"`
	OccasionGUID string     `json:"occasion_guid"`
	Name         string     `json:"name"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at"`

	// Default Shopping List Has One Occasion
	Occasions *Occasion `json:"occasions,omitempty" gorm:"ForeignKey:OccasionGUID;AssociationForeignKey:GUID"`

	// Default Shopping List Has many Default Shopping List Item
	Items []*DefaultShoppingListItem `json:"items,omitempty" gorm:"ForeignKey:ShoppingListGUID;AssociationForeignKey:GUID"`
}
