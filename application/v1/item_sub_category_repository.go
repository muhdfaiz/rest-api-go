package v1

import "github.com/jinzhu/gorm"

type ItemSubCategoryRepositoryInterface interface {
	GetByID(id int) *ItemSubCategory
}

type ItemSubCategoryRepository struct {
	DB *gorm.DB
}

// GetByID function used to retrieve Item Sub Category by ID
func (iscr *ItemSubCategoryRepository) GetByID(id int) *ItemSubCategory {
	itemSubCategory := &ItemSubCategory{}

	iscr.DB.Model(&ItemSubCategory{}).Where(&ItemSubCategory{ID: id}).First(&itemSubCategory)

	return itemSubCategory
}
