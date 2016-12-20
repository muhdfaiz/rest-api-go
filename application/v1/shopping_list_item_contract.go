package v1

import "bitbucket.org/cliqers/shoppermate-api/systems"

// ShoppingListItemServiceInterface is a contract that defines the methods needed for Shopping List Item Service
type ShoppingListItemServiceInterface interface {
	ViewUserShoppingListItem(userGUID string, shoppingListGUID string,
		shoppingListItemGUID string, relations string) (*ShoppingListItem, *systems.ErrorData)
	ViewAllUserShoppingListItem(userGUID string, shoppingListGUID string, addedToCart string,
		latitude string, longitude string, relations string) (map[string][]*ShoppingListItem, *systems.ErrorData)
	CreateUserShoppingListItem(userGUID string, shoppingListGUID string,
		shoppingListItemToCreate CreateShoppingListItem) (*ShoppingListItem, *systems.ErrorData)
	CreateUserShoppingListItemAddedFromDeal(shoppingListItemToCreate CreateShoppingListItem) (*ShoppingListItem, *systems.ErrorData)
	UpdateUserShoppingListItem(userGUID string, shoppingListGUID string, shoppingListItemGUID string,
		shoppingListItemToUpdate UpdateShoppingListItem, relations string) (*ShoppingListItem, *systems.ErrorData)
	UpdateAllUserShoppingListItem(userGUID string, shoppingListGUID string,
		shoppingListItemToUpdate UpdateShoppingListItem) ([]*ShoppingListItem, *systems.ErrorData)
	DeleteUserShoppingListItem(userGUID string, shoppingListGUID string,
		deleteItemInCart string) (map[string]string, *systems.ErrorData)
	DeleteShoppingListItemInShoppingList(shoppingListItemGUID string, userGUID string,
		shoppingListGUID string) (map[string]string, *systems.ErrorData)
	DeleteAllShoppingListItemsInShoppingList(userGUID string, shoppingListGUID string) (map[string]string, *systems.ErrorData)
	DeleteShoppingListItemHasBeenAddtoCart(userGUID string, shoppingListGUID string) (map[string]string, *systems.ErrorData)
	DeleteShoppingListItemHasNotBeenAddtoCart(userGUID string, shoppingListGUID string) (map[string]string, *systems.ErrorData)
	GetAllUserShoppingListItem(userGUID string, shoppingListGUID string, relations string, latitude string,
		longitude string) map[string][]*ShoppingListItem
	GetUserShoppingListItemsNotAddedToCart(userGUID string, shoppingListGUID string, relations string,
		latitude string, longitude string) map[string][]*ShoppingListItem
	GetUserShoppingListItemsAddedToCart(userGUID string, shoppingListGUID string, relations string,
		latitude string, longitude string) map[string][]*ShoppingListItem
	CheckUserShoppingListItemExistOrNot(shoppingListItemGUID string, userGUID string,
		shoppingListGUID string) (*ShoppingListItem, *systems.ErrorData)
	SetShoppingListItemCategoryAndSubcategory(shoppingListItemName string) (string, string)
	GetAndSetDealForShoppingListItems(userGUID string, shoppingListGUID string, userShoppingListItems []*ShoppingListItem,
		latitude string, longitude string) []*ShoppingListItem
}

// ShoppingListItemRepositoryInterface is a contract that defines the method needed for Shopping List Item Repository.
type ShoppingListItemRepositoryInterface interface {
	Create(data CreateShoppingListItem) (*ShoppingListItem, *systems.ErrorData)
	UpdateByUserGUIDShoppingListGUIDAndShoppingListItemGUID(userGUID string, shoppingListGUID string, shoppingListItemGUID string,
		data map[string]interface{}) *systems.ErrorData
	UpdateByUserGUIDAndShoppingListGUID(userGUID string, shoppingListGUID string, data map[string]interface{}) *systems.ErrorData
	UpdateByUserGUIDAndDealGUID(userGUID string, dealGUID string, data map[string]interface{}) *systems.ErrorData
	UpdateByUserGUIDShoppingListGUIDAndDealGUID(userGUID string, shoppingListGUID string, dealGUID string,
		data map[string]interface{}) *systems.ErrorData
	SetDealExpired(dealGUID string) *systems.ErrorData
	DeleteByGUID(shoppingListItemGUID string) *systems.ErrorData
	DeleteByShoppingListGUID(shoppingListGUID string) *systems.ErrorData
	DeleteByUserGUIDAndShoppingListGUID(userGUID string, shoppingListGUID string) *systems.ErrorData
	DeleteByGUIDAndUserGUIDAndShoppingListGUID(shoppingListItemGUID string,
		userGUID string, shoppingListGUID string) *systems.ErrorData
	DeleteItemsHasBeenAddedToCartByUserGUIDAndShoppingListGUID(userGUID string, shoppingListGUID string) *systems.ErrorData
	DeleteItemsHasNotBeenAddedToCartByUserGUIDAndShoppingListGUID(userGUID string, shoppingListGUID string) *systems.ErrorData
	GetByName(name string, relations string) *ShoppingListItem
	GetByGUID(guid string, relations string) *ShoppingListItem
	GetByGUIDUserGUIDAndShoppingListGUID(userGUID string, shoppingListGUID string,
		shoppingListItemGUID string, relations string) *ShoppingListItem
	GetByUserGUIDAndShoppingListGUID(userGUID string, shoppingListGUID string, relations string) []*ShoppingListItem
	GetByUserGUIDAndShoppingListGUIDAndSubCategory(userGUID string, shoppingListGUID string,
		subcategory string, relations string) []*ShoppingListItem
	GetByUserGUIDAndShoppingListGUIDAndAddedToCartAndSubCategory(userGUID string, shoppingListGUID string,
		addedToCart int, subcategory string, relations string) []*ShoppingListItem
	GetUniqueSubCategoryFromAllUserShoppingListItem(userGUID string, shoppingListGUID string) []*ShoppingListItem
	GetUniqueSubCategoryFromUserShoppingListItem(userGUID string, shoppingListGUID string,
		addedToCart int) []*ShoppingListItem
}
