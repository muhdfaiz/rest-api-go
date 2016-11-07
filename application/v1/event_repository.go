package v1

import "github.com/jinzhu/gorm"

type EventRepositoryInterface interface {
	GeAllIncludingDeals() []*Event
}

type EventRepository struct {
	DB *gorm.DB
}

func (er *EventRepository) GeAllIncludingDeals() []*Event {
	events := []*Event{}

	er.DB.Preload("Deals").Find(&events)

	return events
}
