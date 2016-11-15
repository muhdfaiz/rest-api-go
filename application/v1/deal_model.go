package v1

import "time"

type Deal struct {
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
	ConversionLocation        string     `json:"conversion_location"`
	RefreshPeriod             int        `json:"refresh_period"`
	Perlimit                  int        `json:"perlimit"`
	CashbackAmount            float64    `json:"cashback_amount"`
	Quota                     int        `json:"quota"`
	Status                    string     `json:"status"`
	GrocerExclusive           int        `json:"grocer_exclusive"`
	Terms                     string     `json:"terms"`
	AddedToList               int        `json:"added_to_list"`
	NearestGrocerDistanceInKm float64    `json:"nearest_grocer_distance_in_km,omitempty"`
	NearestGrocerName         string     `json:"nearest_grocer_name,omitempty"`
	NearestGrocerLatitude     float64    `json:"nearest_grocer_latitude,omitempty"`
	NearestGrocerLongitude    float64    `json:"nearest_grocer_longitude,omitempty"`
	TotalDealCashback         int        `json:"total_deal_cashback,omitempty"`
	TotalUserDealCashback     int        `json:"total_user_deal_cashback,omitempty"`
	CreatedAt                 time.Time  `json:"created_at"`
	UpdatedAt                 time.Time  `json:"updated_at"`
	DeletedAt                 *time.Time `json:"deleted_at"`

	// Has One Shopping List Item
	Items *Item `json:"items,omitempty" gorm:"ForeignKey:ItemID;AssociationForeignKey:ID"`

	// Has One Targeted Item Category
	Category *ItemCategory `json:"category,omitempty" gorm:"ForeignKey:CategoryID;AssociationForeignKey:ID"`

	// Grocer Exclusive
	Grocerexclusives *Grocer `json:"grocer_exclusives,omitempty" gorm:"ForeignKey:GrocerExclusive;AssociationForeignKey:ID"`

	// Have Many Grocers
	Grocers []*Grocer `json:"grocer_locations,omitempty" gorm:"many2many:ads_grocer;"`
}

// TableName function used to set Item table name to be `item`
func (d Deal) TableName() string {
	return "ads"
}

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
	ConversionLocation        string     `json:"conversion_location"`
	RefreshPeriod             int        `json:"refresh_period"`
	Perlimit                  int        `json:"perlimit"`
	CashbackAmount            float64    `json:"cashback_amount"`
	Quota                     int        `json:"quota"`
	Status                    string     `json:"status"`
	GrocerExclusive           int        `json:"grocer_exclusive"`
	Terms                     string     `json:"terms"`
	CanAddTolist              int        `json:"can_add_to_list"`
	NearestGrocerDistanceInKm float64    `json:"nearest_grocer_distance_in_km,omitempty"`
	NearestGrocerName         string     `json:"nearest_grocer_name,omitempty"`
	NearestGrocerLatitude     float64    `json:"nearest_grocer_latitude,omitempty"`
	NearestGrocerLongitude    float64    `json:"nearest_grocer_longitude,omitempty"`
	TotalDealCashback         int        `json:"total_deal_cashback,omitempty"`
	TotalUserDealCashback     int        `json:"total_user_deal_cashback,omitempty"`
	CreatedAt                 time.Time  `json:"created_at"`
	UpdatedAt                 time.Time  `json:"updated_at"`
	DeletedAt                 *time.Time `json:"deleted_at"`

	// Grocer Exclusive
	Grocerexclusives *Grocer `json:"grocer_exclusives,omitempty" gorm:"ForeignKey:GrocerExclusive;AssociationForeignKey:ID"`

	// Has One Shopping List Item
	Items *Item `json:"items,omitempty" gorm:"ForeignKey:ItemID;AssociationForeignKey:ID"`

	// Has One Targeted Item Category
	Category *ItemCategory `json:"category,omitempty" gorm:"ForeignKey:CategoryID;AssociationForeignKey:ID"`

	// Have Many Grocers
	Grocers []*Grocer `json:"grocers,omitempty" gorm:"many2many:ads_grocer;"`
}
