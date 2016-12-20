package v1

import "bitbucket.org/cliqers/shoppermate-api/systems"

// ShoppingListServiceInterface is a contract that defines the methods needed for Shopping List Service
type ShoppingListServiceInterface interface {
	CreateUserShoppingList(userGUID string, createData CreateShoppingList) (*ShoppingList, *systems.ErrorData)
	UpdateUserShoppingList(userGUID string, shoppingListGUID string, updateData UpdateShoppingList) (*ShoppingList, *systems.ErrorData)
	DeleteUserShoppingListIncludingItemsAndImages(userGUID string, shoppingListGUID string) *systems.ErrorData
	GetUserShoppingLists(userGUID string, relations string) ([]*ShoppingList, *systems.ErrorData)
	ViewShoppingListByGUID(shoppingListGUID string, relations string) *ShoppingList
	CheckUserShoppingListDuplicate(userGUID string, shoppingListName string, occasionGUID string) *systems.ErrorData
	CheckUserShoppingListExistOrNot(userGUID string, shoppingListGUID string) (*ShoppingList, *systems.ErrorData)
	GetShoppingListIncludingDealCashbacks(shoppingListGUID string, dealCashbackTransactionGUID string) *ShoppingList
	CreateSampleShoppingListsAndItemsForUser(userGUID string) *systems.ErrorData
	createSampleShoppingListItems(userGUID string, shoppingListGUID string) *systems.ErrorData
}

// ShoppingListRepositoryInterface is a contract that defines the methods needed for Shopping List Repository.
type ShoppingListRepositoryInterface interface {
	Create(userGUID string, data CreateShoppingList) (*ShoppingList, *systems.ErrorData)
	Update(userGUID string, shoppingListGUID string, data UpdateShoppingList) *systems.ErrorData
	Delete(attribute string, value string) *systems.ErrorData
	GetByUserGUID(userGUID string, relations string) []*ShoppingList
	GetByGUID(GUID string, relations string) *ShoppingList
	GetByGUIDPreloadWithDealCashbacks(GUID string, dealCashbackTransactionGUID string, relations string) *ShoppingList
	GetByGUIDAndUserGUID(GUID string, userGUID string, relations string) *ShoppingList
	GetByUserGUIDOccasionGUIDAndName(userGUID string, name string, occasionGUID string, relations string) *ShoppingList
}
