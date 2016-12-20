package v1

import "github.com/jinzhu/gorm"

// DefaultShoppingListRepository will handle all task related to CRUD
type DefaultShoppingListRepository struct {
	DB *gorm.DB
}

// GetAll function used to retrieve all default shopping lists from database.
func (dslr *DefaultShoppingListRepository) GetAll(relations string) []*DefaultShoppingList {
	defaultShoppingLists := []*DefaultShoppingList{}

	DB := dslr.DB.Model(&DefaultShoppingList{})

	if relations != "" {
		DB = LoadRelations(DB, relations)
	}

	DB.Find(&defaultShoppingLists)

	return defaultShoppingLists
}
