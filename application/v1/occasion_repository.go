package v1

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
)

type OccasionRepositoryInterface interface {
	GetAll() []Occasion
	GetLatestUpdate(lastSyncDate string) []Occasion
	GetByGUID(guid string) *Occasion
}

type OccasionRepository struct {
	DB *gorm.DB
}

func (or *OccasionRepository) GetAll() []Occasion {
	rows, _ := or.DB.Model(&Occasion{}).Rows()

	var occasions []Occasion

	for rows.Next() {
		var occasion Occasion

		if err := rows.Scan(&occasion.ID, &occasion.GUID, &occasion.Slug, &occasion.Name, &occasion.Image, &occasion.Active, &occasion.CreatedAt, &occasion.UpdatedAt, &occasion.DeletedAt); err != nil {
			fmt.Println(err)
		}

		occasions = append(occasions, occasion)
	}
	return occasions
}

func (or *OccasionRepository) GetLatestUpdate(lastSyncDate string) []Occasion {
	lastSync, _ := time.Parse(time.RFC3339, lastSyncDate)

	rows, _ := or.DB.Table("occasions").Where("updated_at > ?", lastSync).Order("updated_at desc").Rows()

	var occasions []Occasion

	for rows.Next() {
		var occasion Occasion

		if err := rows.Scan(&occasion.ID, &occasion.GUID, &occasion.Slug, &occasion.Name, &occasion.Image, &occasion.Active, &occasion.CreatedAt, &occasion.UpdatedAt, &occasion.DeletedAt); err != nil {
			fmt.Println(err)
		}

		occasions = append(occasions, occasion)
	}
	return occasions
}

// GetByGUID function used to retrieve occasion by guid.
// Return occasion data if found and return empty occasion if not found
func (or *OccasionRepository) GetByGUID(guid string) *Occasion {
	result := or.DB.Where(&Occasion{GUID: guid}).First(&Occasion{})

	if result.RowsAffected == 0 {
		return &Occasion{}
	}

	return result.Value.(*Occasion)
}
