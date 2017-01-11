package v1_1

import "bitbucket.org/cliqers/shoppermate-api/systems"
import "github.com/jinzhu/gorm"

// ShoppingListItemServiceInterface is a contract that defines the methods needed for Shopping List Item Service
type ShoppingListItemServiceInterface interface {
	ViewUserShoppingListItem(userGUID string, shoppingListGUID string,
		shoppingListItemGUID string, relations string) (*ShoppingListItem, *systems.ErrorData)
	ViewAllUserShoppingListItem(dbTransaction *gorm.DB, userGUID string, shoppingListGUID string, addedToCart string,
		latitude string, longitude string, relations string) (map[string][]*ShoppingListItem, *systems.ErrorData)
	CreateUserShoppingListItem(dbTransaction *gorm.DB, userGUID string, shoppingListGUID string,
		shoppingListItemToCreate CreateShoppingListItem) (*ShoppingListItem, *systems.ErrorData)
	CreateUserShoppingListItemAddedFromDeal(dbTransaction *gorm.DB, shoppingListItemToCreate CreateShoppingListItem) (*ShoppingListItem, *systems.ErrorData)
	UpdateUserShoppingListItem(dbTransaction *gorm.DB, userGUID string, shoppingListGUID string, shoppingListItemGUID string,
		shoppingListItemToUpdate UpdateShoppingListItem, relations string) (*ShoppingListItem, *systems.ErrorData)
	UpdateAllUserShoppingListItem(dbTransaction *gorm.DB, userGUID string, shoppingListGUID string,
		shoppingListItemToUpdate UpdateShoppingListItem) ([]*ShoppingListItem, *systems.ErrorData)
	DeleteUserShoppingListItem(dbTransaction *gorm.DB, userGUID string, shoppingListGUID string,
		deleteItemInCart string) (map[string]string, *systems.ErrorData)
	DeleteShoppingListItemInShoppingList(dbTransaction *gorm.DB, shoppingListItemGUID string, userGUID string,
		shoppingListGUID string) (map[string]string, *systems.ErrorData)
	DeleteAllShoppingListItemsInShoppingList(dbTransaction *gorm.DB, userGUID string, shoppingListGUID string) (map[string]string, *systems.ErrorData)
	DeleteShoppingListItemHasBeenAddtoCart(dbTransaction *gorm.DB, userGUID string, shoppingListGUID string) (map[string]string, *systems.ErrorData)
	DeleteShoppingListItemHasNotBeenAddtoCart(dbTransaction *gorm.DB, userGUID string, shoppingListGUID string) (map[string]string, *systems.ErrorData)
	GetAllUserShoppingListItem(dbTransaction *gorm.DB, userGUID string, shoppingListGUID string, relations string, latitude string,
		longitude string) (map[string][]*ShoppingListItem, *systems.ErrorData)
	GetUserShoppingListItemsNotAddedToCart(dbTransaction *gorm.DB, userGUID string, shoppingListGUID string, relations string,
		latitude string, longitude string) (map[string][]*ShoppingListItem, *systems.ErrorData)
	GetUserShoppingListItemsAddedToCart(userGUID string, shoppingListGUID string, relations string,
		latitude string, longitude string) map[string][]*ShoppingListItem
	GetShoppingListItemByGUID(shoppingListItemGUID, relations string) *ShoppingListItem
	GetShoppingListItemsByUserGUIDAndShoppingListGUID(userGUID, shoppingListGUID string, relations string) []*ShoppingListItem
	CheckUserShoppingListItemExistOrNot(shoppingListItemGUID string, userGUID string,
		shoppingListGUID string) (*ShoppingListItem, *systems.ErrorData)
	SetShoppingListItemCategoryAndSubcategory(shoppingListItemName string) (string, string)
	GetAndSetDealForShoppingListItems(dbTransaction *gorm.DB, dealsCollection []*Deal, userGUID string, shoppingListGUID string,
		userShoppingListItems []*ShoppingListItem, latitude string, longitude string) ([]*ShoppingListItem, []*Deal, *systems.ErrorData)
}

// ShoppingListItemRepositoryInterface is a contract that defines the method needed for Shopping List Item Repository.
type ShoppingListItemRepositoryInterface interface {
	Create(dbTransaction *gorm.DB, data CreateShoppingListItem) (*ShoppingListItem, *systems.ErrorData)
	UpdateByUserGUIDShoppingListGUIDAndShoppingListItemGUID(dbTransaction *gorm.DB, userGUID string, shoppingListGUID string,
		shoppingListItemGUID string, data map[string]interface{}) *systems.ErrorData
	UpdateByUserGUIDAndShoppingListGUID(dbTransaction *gorm.DB, userGUID string, shoppingListGUID string,
		data map[string]interface{}) *systems.ErrorData
	UpdateByUserGUIDAndDealGUID(dbTransaction *gorm.DB, userGUID string, dealGUID string, data map[string]interface{}) *systems.ErrorData
	UpdateByUserGUIDShoppingListGUIDAndDealGUID(dbTransaction *gorm.DB, userGUID string, shoppingListGUID string, dealGUID string,
		data map[string]interface{}) *systems.ErrorData
	SetDealExpired(dbTransaction *gorm.DB, userGUID, shoppingListGUID, dealGUID string) *systems.ErrorData
	DeleteByGUID(dbTransaction *gorm.DB, shoppingListItemGUID string) *systems.ErrorData
	DeleteByShoppingListGUID(dbTransaction *gorm.DB, shoppingListGUID string) *systems.ErrorData
	DeleteByUserGUIDAndShoppingListGUID(dbTransaction *gorm.DB, userGUID string, shoppingListGUID string) *systems.ErrorData
	DeleteByGUIDAndUserGUIDAndShoppingListGUID(dbTransaction *gorm.DB, shoppingListItemGUID string,
		userGUID string, shoppingListGUID string) *systems.ErrorData
	DeleteItemsHasBeenAddedToCartByUserGUIDAndShoppingListGUID(dbTransaction *gorm.DB, userGUID string, shoppingListGUID string) *systems.ErrorData
	DeleteItemsHasNotBeenAddedToCartByUserGUIDAndShoppingListGUID(dbTransaction *gorm.DB, userGUID string, shoppingListGUID string) *systems.ErrorData
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
