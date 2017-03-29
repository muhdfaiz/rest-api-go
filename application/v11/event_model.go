package v11

import "time"

// Event Model
type Event struct {
	ID        int        `json:"id"`
	GUID      string     `json:"guid"`
	Name      string     `json:"name"`
	Color     string     `json:"color"`
	Img       string     `json:"img"`
	StartDate time.Time  `json:"start_date"`
	EndDate   time.Time  `json:"end_date"`
	Status    string     `json:"status"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`

	//Eventdeals []*EventDeal `json:"deals,omitempty" gorm:"ForeignKey:EventID;AssociationForeignKey:ID"`

	// Event Has Many Deals
	Deals []*Deal `json:"deals,omitempty" gorm:"many2many:event_deal;"`
}

// TableName function used to override default plural table name used by gorm based on struct name.
func (e Event) TableName() string {
	return "event"
}
