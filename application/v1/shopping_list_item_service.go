package v1

import (
	"strconv"
	"time"

	"github.com/fatih/structs"
	"github.com/jinzhu/gorm"

	"bitbucket.org/cliqers/shoppermate-api/systems"
)

type ShoppingListItemService struct {
	ShoppingListItemRepository ShoppingListItemRepositoryInterface
	ItemService                ItemServiceInterface
	ItemCategoryService        ItemCategoryServiceInterface
	ItemSubCategoryService     ItemSubCategoryServiceInterface
	DealService                DealServiceInterface
	DealRepository             DealRepositoryInterface
	GenericService             GenericServiceInterface
}

// ViewUserShoppingListItem function used to view details of one user shopping list item inside user shopping list
func (slis *ShoppingListItemService) ViewUserShoppingListItem(userGUID string, shoppingListGUID string,
	shoppingListItemGUID string, relations string) (*ShoppingListItem, *systems.ErrorData) {

	shoppingListItem := slis.ShoppingListItemRepository.GetByGUIDUserGUIDAndShoppingListGUID(userGUID, shoppingListGUID, shoppingListItemGUID, relations)

	if shoppingListItem.GUID == "" {
		return nil, Error.ResourceNotFoundError("Shopping List Item", "guid", shoppingListItemGUID)
	}

	return shoppingListItem, nil
}

// ViewAllUserShoppingListItem function used to view details of all user shopping list item inside user shopping list
func (slis *ShoppingListItemService) ViewAllUserShoppingListItem(dbTransaction *gorm.DB, userGUID string, shoppingListGUID string, addedToCart string,
	latitude string, longitude string, relations string) (map[string][]*ShoppingListItem, *systems.ErrorData) {

	addedToCartBool, error := strconv.ParseBool(addedToCart)

	if error != nil {
		userShoppingListItems, error1 := slis.GetAllUserShoppingListItem(dbTransaction, userGUID, shoppingListGUID, relations, latitude, longitude)

		if error1 != nil {
			return nil, error1
		}

		return userShoppingListItems, nil
	}

	if addedToCartBool == true {
		userShoppingListItems := slis.GetUserShoppingListItemsAddedToCart(userGUID, shoppingListGUID, relations, latitude, longitude)

		return userShoppingListItems, nil
	}

	userShoppingListItems, error1 := slis.GetUserShoppingListItemsNotAddedToCart(dbTransaction, userGUID, shoppingListGUID, relations, latitude, longitude)

	if error1 != nil {
		return nil, error1
	}

	return userShoppingListItems, nil
}

// CreateUserShoppingListItem function used to create user shopping list item and store in database
func (slis *ShoppingListItemService) CreateUserShoppingListItem(dbTransaction *gorm.DB, userGUID string, shoppingListGUID string,
	shoppingListItemToCreate CreateShoppingListItem) (*ShoppingListItem, *systems.ErrorData) {

	itemCategory, itemSubcategory := slis.SetShoppingListItemCategoryAndSubcategory(shoppingListItemToCreate.Name)

	shoppingListItemToCreate.ShoppingListGUID = shoppingListGUID
	shoppingListItemToCreate.UserGUID = userGUID
	shoppingListItemToCreate.Category = itemCategory
	shoppingListItemToCreate.SubCategory = itemSubcategory

	createdShoppingListItem, error := slis.ShoppingListItemRepository.Create(dbTransaction, shoppingListItemToCreate)

	if error != nil {
		return nil, error
	}

	return createdShoppingListItem, nil
}

// CreateUserShoppingListItemAddedFromDeal function used to create user shopping list item during adding deal to list.
func (slis *ShoppingListItemService) CreateUserShoppingListItemAddedFromDeal(dbTransaction *gorm.DB,
	shoppingListItemToCreate CreateShoppingListItem) (*ShoppingListItem, *systems.ErrorData) {

	itemCategory, itemSubcategory := slis.SetShoppingListItemCategoryAndSubcategory(shoppingListItemToCreate.Name)

	shoppingListItemToCreate.Category = itemCategory
	shoppingListItemToCreate.SubCategory = itemSubcategory

	createdShoppingListItem, error := slis.ShoppingListItemRepository.Create(dbTransaction, shoppingListItemToCreate)

	if error != nil {
		return nil, error
	}

	return createdShoppingListItem, nil
}

// UpdateUserShoppingListItem function used to update one of the user shopping list item inside shopping list  in database
func (slis *ShoppingListItemService) UpdateUserShoppingListItem(dbTransaction *gorm.DB, userGUID string, shoppingListGUID string,
	shoppingListItemGUID string, shoppingListItemToUpdate UpdateShoppingListItem, relations string) (*ShoppingListItem, *systems.ErrorData) {

	shoppingListItem := slis.ShoppingListItemRepository.GetByGUIDUserGUIDAndShoppingListGUID(userGUID, shoppingListGUID,
		shoppingListItemGUID, relations)

	if shoppingListItem.GUID == "" {
		return nil, Error.ResourceNotFoundError("Shopping List Item", "guid", shoppingListItemGUID)
	}

	error := slis.ShoppingListItemRepository.UpdateByUserGUIDShoppingListGUIDAndShoppingListItemGUID(dbTransaction, userGUID, shoppingListGUID,
		shoppingListItemGUID, structs.Map(shoppingListItemToUpdate))

	if error != nil {
		return nil, error
	}

	updatedShoppingListItem := slis.ShoppingListItemRepository.GetByGUID(shoppingListItem.GUID, "")

	return updatedShoppingListItem, nil
}

// UpdateAllUserShoppingListItem function used to update all of the user shopping list item inside shopping list in database
func (slis *ShoppingListItemService) UpdateAllUserShoppingListItem(dbTransaction *gorm.DB, userGUID string, shoppingListGUID string,
	shoppingListItemToUpdate UpdateShoppingListItem) ([]*ShoppingListItem, *systems.ErrorData) {

	error := slis.ShoppingListItemRepository.UpdateByUserGUIDAndShoppingListGUID(dbTransaction, userGUID,
		shoppingListGUID, structs.Map(shoppingListItemToUpdate))

	if error != nil {
		return nil, error
	}

	updatedShoppingListItems := slis.ShoppingListItemRepository.GetByUserGUIDAndShoppingListGUID(userGUID, shoppingListGUID, "")

	return updatedShoppingListItems, nil
}

// DeleteUserShoppingListItem function used to soft delete all of the user shopping list items inside shopping list,
// soft delete all of the user shopping list items has been added to cart and has not been added to cart.
func (slis *ShoppingListItemService) DeleteUserShoppingListItem(dbTransaction *gorm.DB, userGUID string,
	shoppingListGUID string, deleteItemInCart string) (map[string]string, *systems.ErrorData) {

	if deleteItemInCart == "1" {
		result, error := slis.DeleteShoppingListItemHasBeenAddtoCart(dbTransaction, userGUID, shoppingListGUID)

		if error != nil {
			return nil, error
		}

		return result, nil
	}

	if deleteItemInCart == "0" {
		result, error := slis.DeleteShoppingListItemHasNotBeenAddtoCart(dbTransaction, userGUID, shoppingListGUID)

		if error != nil {
			return nil, error
		}

		return result, nil
	}

	result, error := slis.DeleteAllShoppingListItemsInShoppingList(dbTransaction, userGUID, shoppingListGUID)

	if error != nil {
		return nil, error
	}

	return result, nil
}

// DeleteShoppingListItemInShoppingList function used to delete user shopping list item by shopping list item GUID,
// user GUID and shopping list GUID.
func (slis *ShoppingListItemService) DeleteShoppingListItemInShoppingList(dbTransaction *gorm.DB, shoppingListItemGUID string,
	userGUID string, shoppingListGUID string) (map[string]string, *systems.ErrorData) {

	_, error := slis.CheckUserShoppingListItemExistOrNot(shoppingListItemGUID, userGUID, shoppingListGUID)

	if error != nil {
		return nil, error
	}

	error = slis.ShoppingListItemRepository.DeleteByGUIDAndUserGUIDAndShoppingListGUID(dbTransaction, shoppingListItemGUID,
		userGUID, shoppingListGUID)

	if error != nil {
		return nil, error
	}

	result := make(map[string]string)

	result["message"] = "Successfully deleted shopping list item for guid " + shoppingListItemGUID

	return result, nil
}

// DeleteAllShoppingListItemsInShoppingList function used to delete all shopping list item those has been added to cart and
// has not been added to cart by user including shopping list item images.
func (slis *ShoppingListItemService) DeleteAllShoppingListItemsInShoppingList(dbTransaction *gorm.DB, userGUID string,
	shoppingListGUID string) (map[string]string, *systems.ErrorData) {

	error := slis.ShoppingListItemRepository.DeleteByUserGUIDAndShoppingListGUID(dbTransaction, userGUID, shoppingListGUID)

	if error != nil {
		return nil, error
	}

	result := make(map[string]string)

	result["message"] = "Successfully deleted all shopping list item for user guid " + userGUID

	return result, nil
}

// DeleteShoppingListItemHasBeenAddtoCart function used to delete all shopping list item those has been added to cart by user
// including shopping list item images.
func (slis *ShoppingListItemService) DeleteShoppingListItemHasBeenAddtoCart(dbTransaction *gorm.DB, userGUID string,
	shoppingListGUID string) (map[string]string, *systems.ErrorData) {

	error := slis.ShoppingListItemRepository.DeleteItemsHasBeenAddedToCartByUserGUIDAndShoppingListGUID(dbTransaction, userGUID, shoppingListGUID)

	if error != nil {
		return nil, error
	}

	result := make(map[string]string)

	result["message"] = "Successfully deleted all shopping list items those has been added to cart for user guid " + userGUID

	return result, nil
}

// DeleteShoppingListItemHasNotBeenAddtoCart function used to delete all shopping list item those has not been added to cart by user
// including shopping list item images.
func (slis *ShoppingListItemService) DeleteShoppingListItemHasNotBeenAddtoCart(dbTransaction *gorm.DB, userGUID string, shoppingListGUID string) (map[string]string, *systems.ErrorData) {
	error := slis.ShoppingListItemRepository.DeleteItemsHasNotBeenAddedToCartByUserGUIDAndShoppingListGUID(dbTransaction, userGUID, shoppingListGUID)

	if error != nil {
		return nil, error
	}

	result := make(map[string]string)

	result["message"] = "Successfully deleted all shopping list items those has not been added to cart for user guid " + userGUID

	return result, nil
}

// GetAllUserShoppingListItem function used to retrieve all user shopping list item by user GUID and shopping list GUID and group
// by subcategory
func (slis *ShoppingListItemService) GetAllUserShoppingListItem(dbTransaction *gorm.DB, userGUID string, shoppingListGUID string, relations string, latitude string,
	longitude string) (map[string][]*ShoppingListItem, *systems.ErrorData) {

	userShoppingListItemsGroupBySubCategory := make(map[string][]*ShoppingListItem)

	uniqueSubCategories := slis.ShoppingListItemRepository.GetUniqueSubCategoryFromAllUserShoppingListItem(userGUID, shoppingListGUID)

	dealsCollection := []*Deal{}

	for _, uniqueSubCategory := range uniqueSubCategories {
		userShoppingListItems := slis.ShoppingListItemRepository.GetByUserGUIDAndShoppingListGUIDAndSubCategory(userGUID, shoppingListGUID,
			uniqueSubCategory.SubCategory, relations)

		if latitude != "" && longitude != "" {
			userShoppingListItems, deals, error := slis.GetAndSetDealForShoppingListItems(dbTransaction, dealsCollection, userGUID, shoppingListGUID, userShoppingListItems, latitude, longitude)

			if error != nil {
				return nil, error
			}

			dealsCollection = append(dealsCollection, deals...)

			userShoppingListItemsGroupBySubCategory[uniqueSubCategory.SubCategory] = userShoppingListItems
		} else {
			userShoppingListItemsGroupBySubCategory[uniqueSubCategory.SubCategory] = userShoppingListItems
		}
	}

	return userShoppingListItemsGroupBySubCategory, nil
}

// GetUserShoppingListItemsNotAddedToCart function used to retrieve shopping list item by user guid and shopping list guid that not added to cart
func (slis *ShoppingListItemService) GetUserShoppingListItemsNotAddedToCart(dbTransaction *gorm.DB, userGUID string, shoppingListGUID string, relations string,
	latitude string, longitude string) (map[string][]*ShoppingListItem, *systems.ErrorData) {

	userShoppingListItemsGroupBySubCategory := make(map[string][]*ShoppingListItem)

	uniqueSubCategories := slis.ShoppingListItemRepository.GetUniqueSubCategoryFromUserShoppingListItem(userGUID, shoppingListGUID, 0)

	dealsCollection := []*Deal{}

	for _, uniqueSubCategory := range uniqueSubCategories {
		userShoppingListItems := slis.ShoppingListItemRepository.GetByUserGUIDAndShoppingListGUIDAndAddedToCartAndSubCategory(userGUID, shoppingListGUID, 0,
			uniqueSubCategory.SubCategory, relations)

		if latitude != "" && longitude != "" {
			userShoppingListItems, deals, error := slis.GetAndSetDealForShoppingListItems(dbTransaction, dealsCollection, userGUID, shoppingListGUID, userShoppingListItems, latitude, longitude)

			if error != nil {
				return nil, error
			}

			dealsCollection = append(dealsCollection, deals...)

			userShoppingListItemsGroupBySubCategory[uniqueSubCategory.SubCategory] = userShoppingListItems
		} else {
			userShoppingListItemsGroupBySubCategory[uniqueSubCategory.SubCategory] = userShoppingListItems
		}
	}

	return userShoppingListItemsGroupBySubCategory, nil
}

// GetUserShoppingListItemsAddedToCart function used to retrieve shopping list item by user guid and shopping list guid that not added to cart
func (slis *ShoppingListItemService) GetUserShoppingListItemsAddedToCart(userGUID string, shoppingListGUID string, relations string,
	latitude string, longitude string) map[string][]*ShoppingListItem {

	userShoppingListItemsGroupBySubCategory := make(map[string][]*ShoppingListItem)

	uniqueSubCategories := slis.ShoppingListItemRepository.GetUniqueSubCategoryFromUserShoppingListItem(userGUID, shoppingListGUID, 1)

	for _, uniqueSubCategory := range uniqueSubCategories {
		userShoppingListItems := slis.ShoppingListItemRepository.GetByUserGUIDAndShoppingListGUIDAndAddedToCartAndSubCategory(userGUID, shoppingListGUID, 1,
			uniqueSubCategory.SubCategory, relations)

		userShoppingListItemsGroupBySubCategory[uniqueSubCategory.SubCategory] = userShoppingListItems
	}

	return userShoppingListItemsGroupBySubCategory
}

// GetShoppingListItemByGUID function used to view shopping list item by GUID.
func (slis *ShoppingListItemService) GetShoppingListItemByGUID(shoppingListItemGUID, relations string) *ShoppingListItem {
	shoppingListItem := slis.ShoppingListItemRepository.GetByGUID(shoppingListItemGUID, "")

	return shoppingListItem
}

// GetShoppingListItemsByUserGUIDAndShoppingListGUID function used to retrieve user shopping list items for specific shopping list.
func (slis *ShoppingListItemService) GetShoppingListItemsByUserGUIDAndShoppingListGUID(userGUID, shoppingListGUID string, relations string) []*ShoppingListItem {
	shoppingListItems := slis.ShoppingListItemRepository.GetByUserGUIDAndShoppingListGUID(userGUID, shoppingListGUID, "")

	return shoppingListItems
}

// CheckUserShoppingListItemExistOrNot function used to check if user shopping list item inside the shopping
// list exist in database or not.
func (slis *ShoppingListItemService) CheckUserShoppingListItemExistOrNot(shoppingListItemGUID string, userGUID string,
	shoppingListGUID string) (*ShoppingListItem, *systems.ErrorData) {

	shoppingList := slis.ShoppingListItemRepository.GetByGUIDUserGUIDAndShoppingListGUID(userGUID, shoppingListGUID, shoppingListItemGUID, "")

	if shoppingList.GUID == "" {
		return nil, Error.ResourceNotFoundError("Shopping List Item", "guid", shoppingListGUID)
	}

	return shoppingList, nil
}

// SetShoppingListItemCategoryAndSubcategory function used to set shopping list item category and subcategory.
// by checking shopping list item name with item name in database. If not exist set category and subcategory
// to `Others`.
func (slis *ShoppingListItemService) SetShoppingListItemCategoryAndSubcategory(shoppingListItemName string) (string, string) {
	item := slis.ItemService.GetItemByName(shoppingListItemName, "")

	shoppingListItemCategory := "Others"
	shoppingListItemSubCategory := "Others"

	if item.GUID != "" {
		shoppingListItemCategory = slis.ItemCategoryService.GetItemCategoryByID(item.CategoryID).Name

		shoppingListItemSubCategory = slis.ItemSubCategoryService.GetItemSubCategoryByID(item.SubcategoryID).Name
	}

	if item.GUID == "" {
		generic := slis.GenericService.GetGenericByName(shoppingListItemName)

		if generic.GUID != "" {
			shoppingListItemCategory = slis.ItemCategoryService.GetItemCategoryByID(generic.CategoryID).Name

			shoppingListItemSubCategory = slis.ItemSubCategoryService.GetItemSubCategoryByID(generic.SubcategoryID).Name
		}
	}

	return shoppingListItemCategory, shoppingListItemSubCategory
}

// GetAndSetDealForShoppingListItems function used to find deals for shopping list items and set the deals to the shopping
// shopping list items.
func (slis *ShoppingListItemService) GetAndSetDealForShoppingListItems(dbTransaction *gorm.DB, dealsCollection []*Deal, userGUID string, shoppingListGUID string,
	userShoppingListItems []*ShoppingListItem, latitude string, longitude string) ([]*ShoppingListItem, []*Deal, *systems.ErrorData) {

	for key, userShoppingListItem := range userShoppingListItems {

		if userShoppingListItem.AddedFromDeal == 0 && userShoppingListItem.AddedToCart == 0 && latitude != "" && longitude != "" {
			deals := slis.DealService.GetDealsBasedOnUserShoppingListItem(userGUID, shoppingListGUID, userShoppingListItem, latitude, longitude)

			deals = slis.DealService.FilteredDealMustBeUniquePerShoppingList(deals, dealsCollection, userGUID)

			dealsCollection = append(dealsCollection, deals...)

			if len(deals) == 0 {
				deals = nil
			}

			userShoppingListItems[key].Deals = deals
		}

		// If user shopping list item was added from dealcheck deal expired or not
		if userShoppingListItem.AddedFromDeal == 1 {
			currentDateInGMT8 := time.Now().UTC().Add(time.Hour * 8).Format("2006-01-02")

			deal := slis.DealRepository.GetDealByGUIDAndValidStartEndDate(*userShoppingListItem.DealGUID, currentDateInGMT8)

			if deal.GUID == "" {
				error := slis.ShoppingListItemRepository.SetDealExpired(dbTransaction, userGUID, shoppingListGUID, *userShoppingListItem.DealGUID)

				if error != nil {
					return nil, nil, error
				}
			}
		}
	}

	return userShoppingListItems, dealsCollection, nil
}
