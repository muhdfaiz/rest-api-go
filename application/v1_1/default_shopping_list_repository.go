package v1_1

import "github.com/jinzhu/gorm"

// DefaultShoppingListRepository will handle all task related to CRUD
type DefaultShoppingListRepository struct {
	BaseRepository
	DB *gorm.DB
}

// GetAll function used to retrieve all default shopping lists from database.
func (dslr *DefaultShoppingListRepository) GetAll(relations string) []*DefaultShoppingList {
	defaultShoppingLists := []*DefaultShoppingList{}

	DB := dslr.DB.Model(&DefaultShoppingList{})

	if relations != "" {
		DB = dslr.LoadRelations(DB, relations)
	}

	DB.Find(&defaultShoppingLists)

	return defaultShoppingLists
}
