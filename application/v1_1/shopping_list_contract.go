package v1_1

import (
	"bitbucket.org/cliqers/shoppermate-api/systems"
	"github.com/jinzhu/gorm"
)

// ShoppingListServiceInterface is a contract that defines the methods needed for Shopping List Service
type ShoppingListServiceInterface interface {
	CreateUserShoppingList(dbTransaction *gorm.DB, userGUID string, createData CreateShoppingList) (*ShoppingList, *systems.ErrorData)
	UpdateUserShoppingList(dbTransaction *gorm.DB, userGUID string, shoppingListGUID string, updateData UpdateShoppingList) *systems.ErrorData
	DeleteUserShoppingListIncludingItemsAndImages(dbTransaction *gorm.DB, userGUID string, shoppingListGUID string) *systems.ErrorData
	GetUserShoppingLists(userGUID string, relations string) ([]*ShoppingList, *systems.ErrorData)
	ViewShoppingListByGUID(shoppingListGUID string, relations string) *ShoppingList
	CheckUserShoppingListDuplicate(userGUID string, shoppingListName string, occasionGUID string) *systems.ErrorData
	CheckUserShoppingListExistOrNot(userGUID string, shoppingListGUID string) (*ShoppingList, *systems.ErrorData)
	GetShoppingListIncludingDealCashbacks(shoppingListGUID string, dealCashbackTransactionGUID string) *ShoppingList
	CreateSampleShoppingListsAndItemsForUser(dbTransaction *gorm.DB, userGUID string) *systems.ErrorData
	createSampleShoppingListItems(dbTransaction *gorm.DB, userGUID string, shoppingListGUID string) *systems.ErrorData
}

// ShoppingListRepositoryInterface is a contract that defines the methods needed for Shopping List Repository.
type ShoppingListRepositoryInterface interface {
	Create(dbTransaction *gorm.DB, userGUID string, data CreateShoppingList) (*ShoppingList, *systems.ErrorData)
	Update(dbTransaction *gorm.DB, userGUID string, shoppingListGUID string, data UpdateShoppingList) *systems.ErrorData
	Delete(dbTransaction *gorm.DB, attribute string, value string) *systems.ErrorData
	GetByUserGUID(userGUID string, relations string) []*ShoppingList
	GetByGUID(GUID string, relations string) *ShoppingList
	GetByGUIDPreloadWithDealCashbacks(GUID string, dealCashbackTransactionGUID string, relations string) *ShoppingList
	GetByGUIDAndUserGUID(GUID string, userGUID string, relations string) *ShoppingList
	GetByUserGUIDOccasionGUIDAndName(userGUID string, name string, occasionGUID string, relations string) *ShoppingList
}
