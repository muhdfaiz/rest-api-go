package v1

import "github.com/jinzhu/gorm"

// DefaultShoppingListRepositoryInterface is a contract that defines the method needed for Default Shopping
// List Repository.
type DefaultShoppingListRepositoryInterface interface {
	GetAll() []*DefaultShoppingList
}

// DefaultShoppingListRepository will handle all task related to CRUD
type DefaultShoppingListRepository struct {
	DB *gorm.DB
}

// GetAll function used to retrieve all default shopping lists from database.
func (dslr *DefaultShoppingListRepository) GetAll() []*DefaultShoppingList {
	defaultShoppingLists := []*DefaultShoppingList{}

	dslr.DB.Model(&DefaultShoppingList{}).Find(&defaultShoppingLists)

	return defaultShoppingLists
}
