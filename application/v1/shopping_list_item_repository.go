package v1

import "github.com/jinzhu/gorm"

type ShoppingListItemRepositoryInterface interface {
	GetByName(name string, relations string) *ShoppingListItem
	GetByGUID(guid string, relations string) *ShoppingListItem
	GetByShoppingListGUIDAndGUID(guid string, shoppingListGUID string, relations string) *ShoppingListItem
}

// ShoppingListItemRepository used to handle all task related to viewing, retrieving shopping list item
type ShoppingListItemRepository struct {
	DB *gorm.DB
}

// GetByGUID function used to retrieve shopping list item by GUID
func (slir *ShoppingListItemRepository) GetByGUID(guid string, relations string) *ShoppingListItem {

	shoppingListItem := &ShoppingListItem{}

	DB := slir.DB.Model(&ShoppingListItem{})

	if relations != "" {
		DB = LoadRelations(DB, relations)
	}

	DB.Where(&ShoppingListItem{GUID: guid}).First(&shoppingListItem)

	return shoppingListItem
}

// GetByName function used to retrieve shopping list item by Name
func (slir *ShoppingListItemRepository) GetByName(name string, relations string) *ShoppingListItem {
	shoppingListItem := &ShoppingListItem{}

	DB := slir.DB.Model(&ShoppingListItem{})

	if relations != "" {
		DB = LoadRelations(DB, relations)
	}

	DB.Where(&ShoppingListItem{Name: name}).First(&shoppingListItem)

	return shoppingListItem
}

// GetByShoppingListGUIDAndGUID function used to retrieve shopping list item by shopping list guid and guid
func (slir *ShoppingListItemRepository) GetByShoppingListGUIDAndGUID(guid string, shoppingListGUID string, relations string) *ShoppingListItem {
	shoppingListItem := &ShoppingListItem{}

	DB := slir.DB.Model(&ShoppingListItem{})

	if relations != "" {
		DB = LoadRelations(DB, relations)
	}

	DB.Where(&ShoppingListItem{GUID: guid, ShoppingListGUID: shoppingListGUID}).First(&shoppingListItem)

	return shoppingListItem
}
