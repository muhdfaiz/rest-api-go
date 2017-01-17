package v1_1

import (
	"time"

	"github.com/jinzhu/gorm"
)

// OccasionRepository will handle task related to retrieve and search shopping list occasion in database
type OccasionRepository struct {
	DB *gorm.DB
}

// GetAllWithActiveStatus function used retrieve all active shopping list occasions available from database
func (or *OccasionRepository) GetAllWithActiveStatus() ([]*Occasion, int) {
	occasions := []*Occasion{}

	or.DB.Model(&Occasion{}).Where(&Occasion{Active: 1}).Order("updated_at desc").Find(&occasions)

	var totalOccasion *int

	or.DB.Model(&Occasion{}).Where(&Occasion{Active: 1}).Count(&totalOccasion)

	return occasions, *totalOccasion
}

// GetLatestUpdateWithActiveStatus function used to retrieve latest active shopping list occasions that happen after
// last sync date in the query string
func (or *OccasionRepository) GetLatestUpdateWithActiveStatus(lastSyncDate string) ([]*Occasion, int) {
	lastSync, _ := time.Parse(time.RFC3339, lastSyncDate)

	occasions := []*Occasion{}

	or.DB.Model(&Occasion{}).Where(&Occasion{Active: 1}).Where("updated_at > ?", lastSync).Order("updated_at desc").Find(&occasions)

	var totalOccasion *int

	or.DB.Model(&Occasion{}).Where(&Occasion{Active: 1}).Where("updated_at > ?", lastSyncDate).Count(&totalOccasion)

	return occasions, *totalOccasion
}

// GetByGUID function used to retrieve shopping list occasion by guid.
// Return shopping list  occasion data if found and return empty occasion if not found
func (or *OccasionRepository) GetByGUID(guid string) *Occasion {
	occasions := &Occasion{}
	or.DB.Where(&Occasion{GUID: guid}).First(&occasions)

	return occasions
}
