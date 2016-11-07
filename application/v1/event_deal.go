package v1

type EventDeal struct {
	EventID int    `json:"event_id"`
	DealID  string `json:"deal_id"`

	Deals []*Deal `json:"deals,omitempty" gorm:"many2many:event_deal;"`
}
