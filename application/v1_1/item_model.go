package v1_1

import "time"

// Item Model
type Item struct {
	ID            int        `json:"id"`
	GUID          string     `json:"guid"`
	GenericID     *int       `json:"generic_id"`
	Name          string     `json:"name"`
	CategoryID    int        `json:"category_id"`
	Category      string     `sql:"-" json:"category"`
	SubcategoryID int        `json:"subcategory_id"`
	Remarks       string     `json:"remarks"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	DeletedAt     *time.Time `json:"deleted_at"`

	// Item Has One Item Category.
	Categories *ItemCategory `json:"item_category,omitempty" gorm:"ForeignKey:CategoryID;AssociationForeignKey:ID"`

	// Item Has One Item Subcategory.
	Subcategories *ItemSubCategory `json:"item_subcategory,omitempty" gorm:"ForeignKey:SubcategoryID;AssociationForeignKey:ID"`
}

// TableName function used to override default plural table name used by gorm based on struct name.
func (i Item) TableName() string {
	return "item"
}
