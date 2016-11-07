package v1

import "time"

type Grocer struct {
	ID              int               `json:"id"`
	GUID            string            `json:"guid"`
	Name            string            `json:"name"`
	Email           string            `json:"email"`
	Img             string            `json:"img"`
	Status          string            `json:"status"`
	CreatedAt       time.Time         `json:"created_at"`
	UpdatedAt       time.Time         `json:"updated_at"`
	DeletedAt       *time.Time        `json:"deleted_at"`
	GrocerLocations []*GrocerLocation `json:"grocer_locations,omitempty" gorm:"ForeignKey:GrocerID;AssociationForeignKey:ID"`
}

// TableName function used to set Item table name to be `item``
func (g Grocer) TableName() string {
	return "grocer"
}
