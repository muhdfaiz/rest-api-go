package v1

import "time"

// DealCashbackStatus model
type DealCashbackStatus struct {
	ID               int        `json:"id"`
	GUID             string     `json:"guid"`
	DealCashbackGUID string     `json:"deal_cashback_guid"`
	Status           string     `json:"status"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
	DeletedAt        *time.Time `json:"deleted_at"`
}

// TableName function used to override default plural table name will be used by gorm.
func (dcwe DealCashbackStatus) TableName() string {
	return "deal_cashback_status"
}
