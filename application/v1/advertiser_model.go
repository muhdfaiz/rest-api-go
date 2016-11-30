package v1

import "time"

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

// TableName function used to set Item table name to be `item``
func (g Advertiser) TableName() string {
	return "advertiser"
}
