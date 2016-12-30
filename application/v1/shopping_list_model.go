package v1

import "time"

// ShoppingList Model
type ShoppingList struct {
	ID           uint       `json:"id"`
	GUID         string     `json:"guid"`
	UserGUID     string     `json:"user_guid"`
	OccasionGUID string     `json:"occasion_guid"`
	Name         string     `json:"name"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at"`

	// Shopping List Has One Occasion
	Occasions *Occasion `json:"occasions,omitempty" gorm:"ForeignKey:OccasionGUID;AssociationForeignKey:GUID"`

	// Shopping List Has many Shopping List Item
	Items []*ShoppingListItem `json:"items,omitempty" gorm:"ForeignKey:ShoppingListGUID;AssociationForeignKey:GUID"`

	// Shopping List Has Many Deal Cashback
	Dealcashbacks []*DealCashback `json:"deal_cashbacks,omitempty" gorm:"ForeignKey:ShoppingListGUID;AssociationForeignKey:GUID"`
}
