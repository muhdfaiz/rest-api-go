package v1

import (
	"strconv"

	"github.com/fatih/structs"

	"fmt"

	"bitbucket.org/cliqers/shoppermate-api/systems"
)

type ShoppingListItemService struct {
	ShoppingListItemRepository ShoppingListItemRepositoryInterface
	ItemService                ItemServiceInterface
	ItemCategoryService        ItemCategoryServiceInterface
	ItemSubCategoryService     ItemSubCategoryServiceInterface
	DealService                DealServiceInterface
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
func (slis *ShoppingListItemService) ViewAllUserShoppingListItem(userGUID string, shoppingListGUID string, addedToCart string,
	latitude string, longitude string, relations string) (map[string][]*ShoppingListItem, *systems.ErrorData) {

	addedToCartBool, error1 := strconv.ParseBool(addedToCart)

	if error1 != nil {
		userShoppingListItems := slis.GetAllUserShoppingListItem(userGUID, shoppingListGUID, relations, latitude, longitude)

		return userShoppingListItems, nil
	}

	if addedToCartBool == true {
		userShoppingListItems := slis.GetUserShoppingListItemsAddedToCart(userGUID, shoppingListGUID, relations, latitude, longitude)

		return userShoppingListItems, nil
	}

	userShoppingListItems := slis.GetUserShoppingListItemsNotAddedToCart(userGUID, shoppingListGUID, relations, latitude, longitude)

	return userShoppingListItems, nil
}

// CreateUserShoppingListItem function used to create user shopping list item and store in database
func (slis *ShoppingListItemService) CreateUserShoppingListItem(userGUID string, shoppingListGUID string,
	shoppingListItemToCreate CreateShoppingListItem) (*ShoppingListItem, *systems.ErrorData) {

	itemCategory, itemSubcategory := slis.SetShoppingListItemCategoryAndSubcategory(shoppingListItemToCreate.Name)

	shoppingListItemToCreate.ShoppingListGUID = shoppingListGUID
	shoppingListItemToCreate.UserGUID = userGUID
	shoppingListItemToCreate.Category = itemCategory
	shoppingListItemToCreate.SubCategory = itemSubcategory

	createdShoppingListItem, error := slis.ShoppingListItemRepository.Create(shoppingListItemToCreate)

	if error != nil {
		return nil, error
	}

	return createdShoppingListItem, nil
}

// CreateUserShoppingListItemAddedFromDeal function used to create user shopping list item during adding deal to list.
func (slis *ShoppingListItemService) CreateUserShoppingListItemAddedFromDeal(shoppingListItemToCreate CreateShoppingListItem) (*ShoppingListItem, *systems.ErrorData) {

	itemCategory, itemSubcategory := slis.SetShoppingListItemCategoryAndSubcategory(shoppingListItemToCreate.Name)

	shoppingListItemToCreate.Category = itemCategory
	shoppingListItemToCreate.SubCategory = itemSubcategory

	createdShoppingListItem, error := slis.ShoppingListItemRepository.Create(shoppingListItemToCreate)

	if error != nil {
		return nil, error
	}

	return createdShoppingListItem, nil
}

// UpdateUserShoppingListItem function used to update one of the user shopping list item inside shopping list  in database
func (slis *ShoppingListItemService) UpdateUserShoppingListItem(userGUID string, shoppingListGUID string, shoppingListItemGUID string,
	shoppingListItemToUpdate UpdateShoppingListItem, relations string) (*ShoppingListItem, *systems.ErrorData) {

	shoppingListItem := slis.ShoppingListItemRepository.GetByGUIDUserGUIDAndShoppingListGUID(userGUID, shoppingListGUID, shoppingListItemGUID, relations)

	if shoppingListItem.GUID == "" {
		return nil, Error.ResourceNotFoundError("Shopping List Item", "guid", shoppingListItemGUID)
	}
	fmt.Println("QUantity:")
	fmt.Println(shoppingListItemToUpdate.Quantity)

	error := slis.ShoppingListItemRepository.UpdateByUserGUIDShoppingListGUIDAndShoppingListItemGUID(userGUID, shoppingListGUID,
		shoppingListItemGUID, structs.Map(shoppingListItemToUpdate))

	if error != nil {
		return nil, error
	}

	// Retrieve updated shopping list item
	updatedShoppingListItem := slis.ShoppingListItemRepository.GetByGUID(shoppingListItem.GUID, "")

	return updatedShoppingListItem, nil
}

// UpdateAllUserShoppingListItem function used to update all of the user shopping list item inside shopping list in database
func (slis *ShoppingListItemService) UpdateAllUserShoppingListItem(userGUID string, shoppingListGUID string,
	shoppingListItemToUpdate UpdateShoppingListItem) ([]*ShoppingListItem, *systems.ErrorData) {

	error := slis.ShoppingListItemRepository.UpdateByUserGUIDAndShoppingListGUID(userGUID, shoppingListGUID, structs.Map(shoppingListItemToUpdate))

	if error != nil {
		return nil, error
	}

	updatedShoppingListItems := slis.ShoppingListItemRepository.GetByUserGUIDAndShoppingListGUID(userGUID, shoppingListGUID, "")

	return updatedShoppingListItems, nil
}

// DeleteUserShoppingListItem function used to soft delete all of the user shopping list items inside shopping list,
// soft delete all of the user shopping list items has been added to cart and has not been added to cart.
func (slis *ShoppingListItemService) DeleteUserShoppingListItem(userGUID string, shoppingListGUID string,
	deleteItemInCart string) (map[string]string, *systems.ErrorData) {

	if deleteItemInCart == "1" {
		result, error := slis.DeleteShoppingListItemHasBeenAddtoCart(userGUID, shoppingListGUID)

		if error != nil {
			return nil, error
		}

		return result, nil
	}

	if deleteItemInCart == "0" {
		result, error := slis.DeleteShoppingListItemHasNotBeenAddtoCart(userGUID, shoppingListGUID)

		if error != nil {
			return nil, error
		}

		return result, nil
	}

	result, error := slis.DeleteAllShoppingListItemsInShoppingList(userGUID, shoppingListGUID)

	if error != nil {
		return nil, error
	}

	return result, nil
}

// DeleteShoppingListItemInShoppingList function used to delete user shopping list item by shopping list item GUID,
// user GUID and shopping list GUID.
func (slis *ShoppingListItemService) DeleteShoppingListItemInShoppingList(shoppingListItemGUID string, userGUID string,
	shoppingListGUID string) (map[string]string, *systems.ErrorData) {

	_, error := slis.CheckUserShoppingListItemExistOrNot(shoppingListItemGUID, userGUID, shoppingListGUID)

	if error != nil {
		return nil, error
	}

	error = slis.ShoppingListItemRepository.DeleteByGUIDAndUserGUIDAndShoppingListGUID(shoppingListItemGUID, userGUID, shoppingListGUID)

	if error != nil {
		return nil, error
	}

	result := make(map[string]string)

	result["message"] = "Successfully deleted shopping list item for guid " + shoppingListItemGUID

	return result, nil
}

// DeleteAllShoppingListItemsInShoppingList function used to delete all shopping list item those has been added to cart and
// has not been added to cart by user including shopping list item images.
func (slis *ShoppingListItemService) DeleteAllShoppingListItemsInShoppingList(userGUID string, shoppingListGUID string) (map[string]string, *systems.ErrorData) {
	error := slis.ShoppingListItemRepository.DeleteByUserGUIDAndShoppingListGUID(userGUID, shoppingListGUID)

	if error != nil {
		return nil, error
	}

	result := make(map[string]string)

	result["message"] = "Successfully deleted all shopping list item for user guid " + userGUID

	return result, nil
}

// DeleteShoppingListItemHasBeenAddtoCart function used to delete all shopping list item those has been added to cart by user
// including shopping list item images.
func (slis *ShoppingListItemService) DeleteShoppingListItemHasBeenAddtoCart(userGUID string, shoppingListGUID string) (map[string]string, *systems.ErrorData) {
	error := slis.ShoppingListItemRepository.DeleteItemsHasBeenAddedToCartByUserGUIDAndShoppingListGUID(userGUID, shoppingListGUID)

	if error != nil {
		return nil, error
	}

	result := make(map[string]string)

	result["message"] = "Successfully deleted all shopping list items those has been added to cart for user guid " + userGUID

	return result, nil
}

// DeleteShoppingListItemHasNotBeenAddtoCart function used to delete all shopping list item those has not been added to cart by user
// including shopping list item images.
func (slis *ShoppingListItemService) DeleteShoppingListItemHasNotBeenAddtoCart(userGUID string, shoppingListGUID string) (map[string]string, *systems.ErrorData) {
	error := slis.ShoppingListItemRepository.DeleteItemsHasNotBeenAddedToCartByUserGUIDAndShoppingListGUID(userGUID, shoppingListGUID)

	if error != nil {
		return nil, error
	}

	result := make(map[string]string)

	result["message"] = "Successfully deleted all shopping list items those has not been added to cart for user guid " + userGUID

	return result, nil
}

// GetAllUserShoppingListItem function used to retrieve all user shopping list item by user GUID and shopping list GUID and group
// by subcategory
func (slis *ShoppingListItemService) GetAllUserShoppingListItem(userGUID string, shoppingListGUID string, relations string, latitude string,
	longitude string) map[string][]*ShoppingListItem {

	userShoppingListItemsGroupBySubCategory := make(map[string][]*ShoppingListItem)

	uniqueSubCategories := slis.ShoppingListItemRepository.GetUniqueSubCategoryFromAllUserShoppingListItem(userGUID, shoppingListGUID)

	dealsCollection := []*Deal{}

	for _, uniqueSubCategory := range uniqueSubCategories {
		userShoppingListItems := slis.ShoppingListItemRepository.GetByUserGUIDAndShoppingListGUIDAndSubCategory(userGUID, shoppingListGUID,
			uniqueSubCategory.SubCategory, relations)

		if latitude != "" && longitude != "" {
			userShoppingListItems, deals := slis.GetAndSetDealForShoppingListItems(dealsCollection, userGUID, shoppingListGUID, userShoppingListItems, latitude, longitude)

			dealsCollection = append(dealsCollection, deals...)

			userShoppingListItemsGroupBySubCategory[uniqueSubCategory.SubCategory] = userShoppingListItems
		} else {
			userShoppingListItemsGroupBySubCategory[uniqueSubCategory.SubCategory] = userShoppingListItems
		}
	}

	return userShoppingListItemsGroupBySubCategory
}

// GetUserShoppingListItemsNotAddedToCart function used to retrieve shopping list item by user guid and shopping list guid that not added to cart
func (slis *ShoppingListItemService) GetUserShoppingListItemsNotAddedToCart(userGUID string, shoppingListGUID string, relations string,
	latitude string, longitude string) map[string][]*ShoppingListItem {

	userShoppingListItemsGroupBySubCategory := make(map[string][]*ShoppingListItem)

	uniqueSubCategories := slis.ShoppingListItemRepository.GetUniqueSubCategoryFromUserShoppingListItem(userGUID, shoppingListGUID, 0)

	dealsCollection := []*Deal{}

	for _, uniqueSubCategory := range uniqueSubCategories {
		userShoppingListItems := slis.ShoppingListItemRepository.GetByUserGUIDAndShoppingListGUIDAndAddedToCartAndSubCategory(userGUID, shoppingListGUID, 0,
			uniqueSubCategory.SubCategory, relations)

		if latitude != "" && longitude != "" {
			userShoppingListItems, deals := slis.GetAndSetDealForShoppingListItems(dealsCollection, userGUID, shoppingListGUID, userShoppingListItems, latitude, longitude)

			dealsCollection = append(dealsCollection, deals...)

			userShoppingListItemsGroupBySubCategory[uniqueSubCategory.SubCategory] = userShoppingListItems
		} else {
			userShoppingListItemsGroupBySubCategory[uniqueSubCategory.SubCategory] = userShoppingListItems
		}
	}

	return userShoppingListItemsGroupBySubCategory
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
func (slis *ShoppingListItemService) GetAndSetDealForShoppingListItems(dealsCollection []*Deal, userGUID string, shoppingListGUID string,
	userShoppingListItems []*ShoppingListItem, latitude string, longitude string) ([]*ShoppingListItem, []*Deal) {

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

		// If user shopping list item was added from deal and not added to cart, check deal expired or not
		if userShoppingListItem.AddedFromDeal == 1 && userShoppingListItem.AddedToCart == 0 {
			slis.DealService.RemoveDealCashbackAndSetItemDealExpired(userGUID, shoppingListGUID, *userShoppingListItem.DealGUID)
		}
	}

	return userShoppingListItems, dealsCollection
}
