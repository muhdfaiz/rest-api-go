package v1_1

import "github.com/jinzhu/gorm"

type EventRepository struct {
	DB *gorm.DB
}

// GetAllIncludingRelations function used to retrieve all event including other relations.
func (er *EventRepository) GetAllIncludingRelations(todayDateInGMT8 string) []*Event {
	events := []*Event{}

	er.DB.Model(&Event{}).Preload("Deals", func(db *gorm.DB) *gorm.DB {
		return db.Where("ads.start_date <= ? AND ads.end_date > ? AND ads.status = ?", todayDateInGMT8, todayDateInGMT8, "publish")
	}).Preload("Deals.Items").Preload("Deals.Category").Preload("Deals.Items.Categories").Preload("Deals.Items.Subcategories").Preload("Deals.Grocerexclusives").Where(&Event{Status: "publish"}).Find(&events)

	return events
}
