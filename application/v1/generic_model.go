package v1

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

	// Generic Has One Category
	Categories *ItemCategory `json:"category" gorm:"ForeignKey:CategoryID;AssociationForeignKey:ID"`

	// Generic Has One Subcategory
	Subcategories *ItemSubCategory `json:"subcategory,omitempty" gorm:"ForeignKey:SubcategoryID;AssociationForeignKey:ID"`
}

// TableName function used to override default table name use by GORM.
// By default, gorm used plural for database table name.
func (g *Generic) TableName() string {
	return "generic"
}