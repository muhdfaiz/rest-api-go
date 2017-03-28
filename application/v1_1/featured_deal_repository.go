package v1_1

import "github.com/jinzhu/gorm"

// FeaturedDealRepository will handle all CRUD function related to resource featured deal.
type FeaturedDealRepository struct {
	BaseRepository
	DB *gorm.DB
}

// GetActiveFeaturedDeals function used to retrieve all featured deals that still active.
func (fdr *FeaturedDealRepository) GetActiveFeaturedDeals(pageNumber, pageLimit, relations string) []*Deal {
	featuredDeals := []*Deal{}

	events := []*Event{}

	offset := fdr.SetOffsetValue(pageNumber, pageLimit)

	if pageLimit != "" && pageNumber != "" {
		fdr.DB.Model(Event{}).Where(&Event{Status: "publish"}).Find(&events).Offset(offset).Limit(pageLimit).Related(&featuredDeals, "Deals")

		return featuredDeals
	}

	fdr.DB.Model(Event{}).Where(&Event{Status: "publish"}).Find(&events).Related(&featuredDeals, "Deals")

	return featuredDeals
}
