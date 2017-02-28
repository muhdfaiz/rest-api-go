package v1_1

import "time"

// Deal Model
// Deal contain two models because from API side perspective, API call Ads as Deal
// but from Admin side perspective, Admin call Ads.
type Deal struct {
	ID              int        `json:"id"`
	GUID            string     `json:"guid"`
	AdvertiserID    int        `json:"advertiser_id"`
	CampaignID      int        `json:"campaign_id"`
	ItemID          int        `json:"item_id"`
	Img             string     `json:"img"`
	FrontName       string     `json:"front_name"`
	Name            string     `json:"name"`
	Body            string     `json:"body"`
	CategoryID      int        `json:"category_id"`
	PositiveTag     string     `json:"positive_tag"`
	NegativeTag     string     `json:"negative_tag"`
	Type            string     `json:"type"`
	StartDate       time.Time  `json:"start_date"`
	EndDate         time.Time  `json:"end_date"`
	Time            string     `json:"time"`
	RefreshPeriod   int        `json:"refresh_period"`
	Perlimit        int        `json:"perlimit"`
	CashbackAmount  float64    `json:"cashback_amount"`
	Quota           int        `json:"quota"`
	Status          string     `json:"status"`
	GrocerExclusive int        `json:"grocer_exclusive"`
	Terms           string     `json:"terms"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	DeletedAt       *time.Time `json:"deleted_at"`

	// Virtual Columns. Use to return the column in the response.
	NumberOfDealAddedToList   int     `sql:"-" json:"number_of_deal_added_to_list"`
	CanAddTolist              int     `sql:"-" json:"can_add_to_list"`
	RemainingAddToList        int     `sql:"-" json:"remaining_add_to_list"`
	NearestGrocerDistanceInKm float64 `sql:"-" json:"nearest_grocer_distance_in_km,omitempty"`
	NearestGrocerName         string  `sql:"-" json:"nearest_grocer_name,omitempty"`
	NearestGrocerLatitude     float64 `sql:"-" json:"nearest_grocer_latitude,omitempty"`
	NearestGrocerLongitude    float64 `sql:"-" json:"nearest_grocer_longitude,omitempty"`
	TotalDealCashback         int     `sql:"-" json:"total_deal_cashback,omitempty"`
	TotalUserDealCashback     int     `sql:"-" json:"total_user_deal_cashback,omitempty"`

	// Deal Has One Item
	Items *Item `json:"items,omitempty" gorm:"ForeignKey:ItemID;AssociationForeignKey:ID"`

	// Deal Has One Item Category
	Category *ItemCategory `json:"category,omitempty" gorm:"ForeignKey:CategoryID;AssociationForeignKey:ID"`

	// Deal Has One Grocer Exclusive
	Grocerexclusives *Grocer `json:"grocer_exclusives,omitempty" gorm:"ForeignKey:GrocerExclusive;AssociationForeignKey:ID"`

	// Deal Has Many Grocers
	Grocers []*Grocer `json:"grocers,omitempty" gorm:"many2many:ads_grocer;"`
}

// TableName function used to override default plural table name used by gorm based on struct name.
func (d Deal) TableName() string {
	return "ads"
}

// Ads Model
type Ads struct {
	ID                        int        `json:"id"`
	GUID                      string     `json:"guid"`
	AdvertiserID              int        `json:"advertiser_id"`
	CampaignID                int        `json:"campaign_id"`
	ItemID                    int        `json:"item_id"`
	Img                       string     `json:"img"`
	FrontName                 string     `json:"front_name"`
	Name                      string     `json:"name"`
	Body                      string     `json:"body"`
	CategoryID                int        `json:"category_id"`
	PositiveTag               string     `json:"positive_tag"`
	NegativeTag               string     `json:"negative_tag"`
	Type                      string     `json:"type"`
	StartDate                 time.Time  `json:"start_date"`
	EndDate                   time.Time  `json:"end_date"`
	Time                      string     `json:"time"`
	RefreshPeriod             int        `json:"refresh_period"`
	Perlimit                  int        `json:"perlimit"`
	CashbackAmount            float64    `json:"cashback_amount"`
	Quota                     int        `json:"quota"`
	Status                    string     `json:"status"`
	GrocerExclusive           *int       `json:"grocer_exclusive"`
	Terms                     string     `json:"terms"`
	NumberOfDealAddedToList   int        `sql:"-" json:"number_of_deal_added_to_list"`
	CanAddTolist              int        `sql:"-" json:"can_add_to_list"`
	RemainingAddToList        int        `sql:"-" json:"remaining_add_to_list"`
	NearestGrocerDistanceInKm float64    `sql:"-" json:"nearest_grocer_distance_in_km,omitempty"`
	NearestGrocerName         string     `sql:"-" json:"nearest_grocer_name,omitempty"`
	NearestGrocerLatitude     float64    `sql:"-" json:"nearest_grocer_latitude,omitempty"`
	NearestGrocerLongitude    float64    `sql:"-" json:"nearest_grocer_longitude,omitempty"`
	TotalDealCashback         int        `sql:"-" json:"total_deal_cashback,omitempty"`
	TotalUserDealCashback     int        `sql:"-" json:"total_user_deal_cashback,omitempty"`
	CreatedAt                 time.Time  `json:"created_at"`
	UpdatedAt                 time.Time  `json:"updated_at"`
	DeletedAt                 *time.Time `json:"deleted_at"`

	// Deal Has One Item
	Items *Item `json:"items,omitempty" gorm:"ForeignKey:ItemID;AssociationForeignKey:ID"`

	// Deal Has One Item Category
	Category *ItemCategory `json:"category,omitempty" gorm:"ForeignKey:CategoryID;AssociationForeignKey:ID"`

	// Deal Has One Grocer Exclusive
	Grocerexclusives *Grocer `json:"grocer_exclusives,omitempty" gorm:"ForeignKey:GrocerExclusive;AssociationForeignKey:ID"`

	// Deal Has Many Grocers
	Grocers []*Grocer `json:"grocers,omitempty" gorm:"many2many:ads_grocer;"`
}
