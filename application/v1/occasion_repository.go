package v1

import (
	"time"

	"github.com/jinzhu/gorm"
)

type OccasionRepositoryInterface interface {
	GetAll() ([]*Occasion, int)
	GetLatestUpdate(lastSyncDate string) ([]*Occasion, int)
	GetByGUID(guid string) *Occasion
}

// OccasionRepository will handle task related to retrieve and search shopping list occasion in database
type OccasionRepository struct {
	DB *gorm.DB
}

// GetAll function used retrieve all shopping list occasions available from database
func (or *OccasionRepository) GetAll() ([]*Occasion, int) {
	occasions := []*Occasion{}

	or.DB.Model(&Occasion{}).Find(&occasions)

	var totalOccasion *int

	or.DB.Model(&Occasion{}).Count(&totalOccasion)

	return occasions, *totalOccasion
}

// GetLatestUpdate function used to retrieve latest shopping list occasions that happen after last sync date in the query string
func (or *OccasionRepository) GetLatestUpdate(lastSyncDate string) ([]*Occasion, int) {
	lastSync, _ := time.Parse(time.RFC3339, lastSyncDate)

	occasions := []*Occasion{}

	or.DB.Model(&Occasion{}).Where("updated_at > ?", lastSync).Order("updated_at desc").Find(&occasions)

	var totalOccasion *int

	or.DB.Model(&Occasion{}).Where("updated_at > ?", lastSyncDate).Count(&totalOccasion)

	return occasions, *totalOccasion
}

// GetByGUID function used to retrieve shopping list occasion by guid.
// Return shopping list  occasion data if found and return empty occasion if not found
func (or *OccasionRepository) GetByGUID(guid string) *Occasion {
	occasions := &Occasion{}
	or.DB.Where(&Occasion{GUID: guid}).First(&occasions)

	return occasions
}
