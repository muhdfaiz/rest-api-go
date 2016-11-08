package v1

import "time"

// Item Model
type Item struct {
	ID            int        `json:"id"`
	GUID          string     `json:"guid"`
	GenericID     *int       `json:"generic_id"`
	Name          string     `json:"name"`
	CategoryID    int        `json:"category_id"`
	Category      string     `json:"category"`
	SubcategoryID int        `json:"subcategory_id"`
	Remarks       string     `json:"remarks"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	DeletedAt     *time.Time `json:"deleted_at"`

	Categories    *ItemCategory    `json:"item_category,omitempty" gorm:"ForeignKey:CategoryID;AssociationForeignKey:ID"`
	Subcategories *ItemSubCategory `json:"item_subcategory,omitempty" gorm:"ForeignKey:SubcategoryID;AssociationForeignKey:ID"`
}

// TableName function used to set Item table name to be `item``
func (i Item) TableName() string {
	return "item"
}
