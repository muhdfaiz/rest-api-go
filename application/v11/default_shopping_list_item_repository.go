package v11

import "github.com/jinzhu/gorm"

// DefaultShoppingListItemRepository will handle all task related to CRUD.
type DefaultShoppingListItemRepository struct {
	DB *gorm.DB
}

// GetAll function used retrieve all default shopping list items from database.
func (dslir *DefaultShoppingListItemRepository) GetAll() []*DefaultShoppingListItem {
	defaultShoppingListItems := []*DefaultShoppingListItem{}

	dslir.DB.Model(&DefaultShoppingListItem{}).Find(&defaultShoppingListItems)

	return defaultShoppingListItems
}
