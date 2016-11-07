package v1

import "time"

type Deal struct {
	ID                 int        `json:"id"`
	GUID               string     `json:"guid"`
	AdvertiserID       int        `json:"advertiser_id"`
	CampaignID         int        `json:"campaign_id"`
	ItemID             int        `json:"item_id"`
	Img                string     `json:"img"`
	FrontName          string     `json:"front_name"`
	Name               string     `json:"name"`
	Body               string     `json:"body"`
	Category           string     `json:"category"`
	PositiveTag        string     `json:"positive_tag"`
	NegativeTag        string     `json:"negative_tag"`
	Type               string     `json:"type"`
	StartDate          time.Time  `json:"start_date"`
	EndDate            time.Time  `json:"end_date"`
	Time               string     `json:"time"`
	ConversionLocation string     `json:"conversion_location"`
	RefreshPeriod      int        `json:"refresh_period"`
	Perlimit           int        `json:"perlimit"`
	CashbackAmount     float64    `json:"cashback_amount"`
	Quota              int        `json:"quota"`
	Status             string     `json:"status"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
	DeletedAt          *time.Time `json:"deleted_at"`

	// Has One Shopping List Item
	Items   *Item     `json:"items,omitempty" gorm:"ForeignKey:ItemID;AssociationForeignKey:ID"`
	Grocers []*Grocer `json:"grocer_locations,omitempty" gorm:"many2many:ads_grocer;"`
}

// TableName function used to set Item table name to be `item``
func (d Deal) TableName() string {
	return "ads"
}

type Ads struct {
	ID                 int        `json:"id"`
	GUID               string     `json:"guid"`
	AdvertiserID       int        `json:"advertiser_id"`
	CampaignID         int        `json:"campaign_id"`
	ItemID             int        `json:"item_id"`
	Img                string     `json:"img"`
	FrontName          string     `json:"front_name"`
	Name               string     `json:"name"`
	Body               string     `json:"body"`
	Category           string     `json:"category"`
	PositiveTag        string     `json:"positive_tag"`
	NegativeTag        string     `json:"negative_tag"`
	Type               string     `json:"type"`
	StartDate          time.Time  `json:"start_date"`
	EndDate            time.Time  `json:"end_date"`
	Time               string     `json:"time"`
	ConversionLocation string     `json:"conversion_location"`
	RefreshPeriod      int        `json:"refresh_period"`
	Perlimit           int        `json:"perlimit"`
	CashbackAmount     float64    `json:"cashback_amount"`
	Quota              int        `json:"quota"`
	Status             string     `json:"status"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
	DeletedAt          *time.Time `json:"deleted_at"`

	// Has One Shopping List Item
	Items   *Item     `json:"item,omitempty" gorm:"ForeignKey:ItemID;AssociationForeignKey:ID"`
	Grocers []*Grocer `json:"grocers,omitempty" gorm:"many2many:ads_grocer;"`
}
