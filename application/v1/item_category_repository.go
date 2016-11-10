package v1

import (
	"github.com/jinzhu/gorm"
)

type ItemCategoryRepositoryInterface interface {
	GetAll() ([]*ItemCategory, int)
	GetAllCategoryNames() ([]string, int)
	GetByID(ID int) *ItemCategory
	GetByGUID(GUID string) *ItemCategory
}

type ItemCategoryRepository struct {
	DB *gorm.DB
}

// GetAll function used to retrieve all shopping list item categories
func (icr *ItemCategoryRepository) GetAll() ([]*ItemCategory, int) {
	itemCategories := []*ItemCategory{}

	icr.DB.Model(&ItemCategory{}).Find(&itemCategories)

	var totalItemCategory *int

	icr.DB.Model(&ItemCategory{}).Count(&totalItemCategory)

	return itemCategories, *totalItemCategory
}

// GetAllCategoryNames function used to retrieve all shopping list item categories name only
func (icr *ItemCategoryRepository) GetAllCategoryNames() ([]string, int) {
	shoppingListCategories := []*ItemCategory{}

	var shoppingListItemCategoryNames []string

	icr.DB.Find(&shoppingListCategories).Pluck("name", &shoppingListItemCategoryNames)

	return shoppingListItemCategoryNames, len(shoppingListItemCategoryNames)
}

// GetByID function used to retrieve shopping list item category by ID
func (icr *ItemCategoryRepository) GetByID(ID int) *ItemCategory {
	itemCategory := &ItemCategory{}

	icr.DB.Model(&ItemCategory{}).Where(&ItemCategory{ID: ID}).First(&itemCategory)

	return itemCategory
}

// GetByGUID function used to retrieve shopping list item category by GUID
func (icr *ItemCategoryRepository) GetByGUID(GUID string) *ItemCategory {
	itemCategory := &ItemCategory{}

	icr.DB.Model(&ItemCategory{}).Where(&ItemCategory{GUID: GUID}).First(&itemCategory)

	return itemCategory
}
