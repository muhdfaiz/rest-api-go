package v1

import "github.com/jinzhu/gorm"

type EventRepositoryInterface interface {
	GetAllIncludingRelations(todayDateInGMT8 string) []*Event
}

type EventRepository struct {
	DB *gorm.DB
}

func (er *EventRepository) GetAllIncludingRelations(todayDateInGMT8 string) []*Event {
	events := []*Event{}

	er.DB.Model(&Event{}).Preload("Deals", func(db *gorm.DB) *gorm.DB {
		return db.Where("ads.start_date <= ? AND ads.end_date > ? AND ads.status = ?", todayDateInGMT8, todayDateInGMT8, "publish")
	}).Preload("Deals.Items").Preload("Deals.Category").Preload("Deals.Items.Categories").Preload("Deals.Items.Subcategories").Find(&events)

	return events
}
