package v1

import "github.com/jinzhu/gorm"

type ShoppingListRepositoryInterface interface {
	GetByUserGUID(userGUID string, relations string) []*ShoppingList
	GetByGUID(GUID string, relations string) *ShoppingList
	GetByGUIDAndUserGUID(GUID string, userGUID string, relations string) *ShoppingList
	GetByUserGUIDOccasionGUIDAndName(userGUID string, name string, occasionGUID string, relations string) *ShoppingList
}

// ShoppingListRepository used to retrieve user shopping list.
type ShoppingListRepository struct {
	DB *gorm.DB
}

// GetByUserGUID function used to retrieve user shopping list by User GUID.
func (slr *ShoppingListRepository) GetByUserGUID(userGUID string, relations string) []*ShoppingList {
	shoppingLists := []*ShoppingList{}

	DB := slr.DB.Model(&ShoppingList{})

	if relations != "" {
		DB = LoadRelations(DB, relations)
	}

	DB.Where(&ShoppingList{UserGUID: userGUID}).Find(&shoppingLists)

	return shoppingLists
}

// GetByGUIDAndUserGUID function used to retrieve user shopping list by and Shopping List GUID.
func (slr *ShoppingListRepository) GetByGUIDAndUserGUID(GUID string, userGUID string, relations string) *ShoppingList {
	shoppingLists := &ShoppingList{}

	DB := slr.DB.Model(&ShoppingList{})

	if relations != "" {
		DB = LoadRelations(DB, relations)
	}

	DB.Where(&ShoppingList{GUID: GUID, UserGUID: userGUID}).First(&shoppingLists)

	return shoppingLists
}

// GetByGUID function used to retrieve user shopping list by and Shopping List GUID.
func (slr *ShoppingListRepository) GetByGUID(GUID string, relations string) *ShoppingList {
	shoppingLists := &ShoppingList{}

	DB := slr.DB.Model(&ShoppingList{})

	if relations != "" {
		DB = LoadRelations(DB, relations)
	}

	DB.Where(&ShoppingList{GUID: GUID}).First(&shoppingLists)

	return shoppingLists
}

// GetByUserGUIDOccasionGUIDAndName function used to retrieve user shopping list by User GUID and Shopping List Name.
func (slr *ShoppingListRepository) GetByUserGUIDOccasionGUIDAndName(userGUID string, name string, occasionGUID string, relations string) *ShoppingList {
	shoppingLists := &ShoppingList{}

	DB := slr.DB.Model(&ShoppingList{})

	if relations != "" {
		DB = LoadRelations(DB, relations)
	}

	DB.Where(&ShoppingList{UserGUID: userGUID, Name: name, OccasionGUID: occasionGUID}).First(&shoppingLists)

	return shoppingLists
}
