package v1

import "github.com/jinzhu/gorm"

// GenericRepository will handle all task related CRUD.
type GenericRepository struct {
	DB *gorm.DB
}

// GetAll function used to retrieve all generic category from database.
func (gr *GenericRepository) GetAll(pageNumber, pageLimit, relations string) ([]*Generic, int) {
	DB := gr.DB.Model(&Generic{})

	if relations != "" {
		DB = LoadRelations(DB, relations)
	}

	offset := SetOffsetValue(pageNumber, pageLimit)

	generics := []*Generic{}

	if pageLimit != "" {
		totalGeneric := []*Generic{}

		DB.Find(&totalGeneric)

		DB.Offset(offset).Limit(pageLimit).Find(&generics)

		return generics, len(totalGeneric)
	}

	DB.Find(&generics)

	return generics, len(generics)
}

// GetByUpdatedAtGreaterThanLastSyncDate function used to retrieve generic categories by updated_at more than last sync date.
func (gr *GenericRepository) GetByUpdatedAtGreaterThanLastSyncDate(lastSyncDate, pageNumber, pageLimit, relations string) ([]*Generic, int) {
	DB := gr.DB.Model(&Generic{})

	offset := SetOffsetValue(pageNumber, pageLimit)

	if relations != "" {
		DB = LoadRelations(DB, relations)
	}

	generics := []*Generic{}

	if pageLimit != "" {
		totalGeneric := []*Generic{}

		DB.Where("updated_at > ?", lastSyncDate).Find(&totalGeneric)

		DB.Where("updated_at > ?", lastSyncDate).Offset(offset).Limit(pageLimit).Find(&generics)

		return generics, len(totalGeneric)
	}

	DB.Where("updated_at > ?", lastSyncDate).Order("updated_at desc").Find(&generics)

	return generics, len(generics)
}

// GetByID function used to retrieve generic category from database by Generic ID.
// Return empty result if generic ID not valid.
func (gr *GenericRepository) GetByID(genericID int, relations string) *Generic {
	DB := gr.DB.Model(&Generic{})

	DB = LoadRelations(DB, relations)

	generic := &Generic{}

	DB.Where(&Generic{ID: genericID}).First(&generic)

	return generic
}
