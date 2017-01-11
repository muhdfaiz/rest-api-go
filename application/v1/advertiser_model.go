package v1

import "time"

// Advertiser Model
type Advertiser struct {
	ID               int        `json:"id"`
	GUID             string     `json:"guid"`
	Fullname         string     `json:"fullname"`
	Email            string     `json:"email"`
	Password         string     `json:"-"`
	Mobile           string     `json:"mobile"`
	Company          string     `json:"company"`
	Address          string     `json:"address"`
	Postcode         string     `json:"postcode"`
	City             string     `json:"city"`
	State            string     `json:"state"`
	TotalCredits     float64    `json:"-"`
	RemainingCredtis float64    `json:"-"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
	DeletedAt        *time.Time `json:"deleted_at"`
}

// TableName function used to set advertiser table name to `advertiser`
// By default, gorm used plural table name.
func (g Advertiser) TableName() string {
	return "advertiser"
}
