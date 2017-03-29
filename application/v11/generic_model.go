package v11

import "time"

// Generic model
type Generic struct {
	ID            int        `json:"id"`
	GUID          string     `json:"guid"`
	CategoryID    int        `json:"category_id"`
	SubcategoryID int        `json:"subcategory_id"`
	Name          string     `json:"name"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	DeletedAt     *time.Time `json:"deleted_at"`

	// Generic Has One Item Category
	Categories *ItemCategory `json:"category" gorm:"ForeignKey:CategoryID;AssociationForeignKey:ID"`

	// Generic Has One Item Subcategory
	Subcategories *ItemSubCategory `json:"subcategory,omitempty" gorm:"ForeignKey:SubcategoryID;AssociationForeignKey:ID"`
}

// TableName function used to override default plural table name used by gorm based on struct name.
func (g *Generic) TableName() string {
	return "generic"
}
