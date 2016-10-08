package v1

import (
	"time"

	"github.com/jinzhu/gorm"
)

type OccasionRepositoryInterface interface {
	GetAll() []*Occasion
	GetLatestUpdate(lastSyncDate string) []*Occasion
	GetByGUID(guid string) *Occasion
}

type OccasionRepository struct {
	DB *gorm.DB
}

func (or *OccasionRepository) GetAll() []*Occasion {
	occasions := []*Occasion{}

	or.DB.Model(&Occasion{}).Find(&occasions)

	return occasions
}

func (or *OccasionRepository) GetLatestUpdate(lastSyncDate string) []*Occasion {
	lastSync, _ := time.Parse(time.RFC3339, lastSyncDate)

	occasions := []*Occasion{}

	or.DB.Table("occasions").Where("updated_at > ?", lastSync).Order("updated_at desc").Find(&occasions)

	return occasions
}

// GetByGUID function used to retrieve occasion by guid.
// Return occasion data if found and return empty occasion if not found
func (or *OccasionRepository) GetByGUID(guid string) *Occasion {
	occasions := &Occasion{}
	or.DB.Where(&Occasion{GUID: guid}).First(&occasions)

	return occasions
}
