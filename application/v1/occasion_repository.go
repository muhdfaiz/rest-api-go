package v1

import "github.com/jinzhu/gorm"

type OccasionRepositoryInterface interface {
	GetAll() *Occasion
	GetByGUID(guid string) *Occasion
}

type OccasionRepository struct {
	DB *gorm.DB
}

func (or *OccasionRepository) GetAll() *Occasion {
	result := or.DB.Find(&Occasion{})

	if result.RowsAffected == 0 {
		return &Occasion{}
	}

	return result.Value.(*Occasion)
}

// GetByGUID function used to retrieve occasion by guid.
// Return occasion data if found and return empty occasion if not found
func (or *OccasionRepository) GetByGUID(guid string) *Occasion {
	result := or.DB.Where(&User{GUID: guid}).First(&User{})

	if result.RowsAffected == 0 {
		return &Occasion{}
	}

	return result.Value.(*Occasion)
}
