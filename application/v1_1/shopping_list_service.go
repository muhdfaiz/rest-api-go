package v1_1

import (
	"bitbucket.org/cliqers/shoppermate-api/systems"
	"github.com/jinzhu/gorm"
)

type ShoppingListService struct {
	ShoppingListRepository         ShoppingListRepositoryInterface
	OccasionService                OccasionServiceInterface
	DefaultShoppingListService     DefaultShoppingListServiceInterface
	DefaultShoppingListItemService DefaultShoppingListItemServiceInterface
	ShoppingListItemService        ShoppingListItemServiceInterface
	ShoppingListItemImageService   ShoppingListItemImageServiceInterface
}

// CreateUserShoppingList function used to create user shopping lists and store in database.
func (sls *ShoppingListService) CreateUserShoppingList(dbTransaction *gorm.DB, userGUID string, createData CreateShoppingList) (*ShoppingList, *systems.ErrorData) {
	_, error := sls.OccasionService.CheckOccassionExistOrNot(createData.OccasionGUID)

	if error != nil {
		return nil, error
	}

	error = sls.CheckUserShoppingListDuplicate(userGUID, createData.Name, createData.OccasionGUID)

	if error != nil {
		return nil, error
	}

	createdShoppingList, error := sls.ShoppingListRepository.Create(dbTransaction, userGUID, createData)

	if error != nil {
		return nil, error
	}

	return createdShoppingList, nil
}

// UpdateUserShoppingList function used to update user shopping lists in database.
func (sls *ShoppingListService) UpdateUserShoppingList(dbTransaction *gorm.DB, userGUID string, shoppingListGUID string,
	updateData UpdateShoppingList) *systems.ErrorData {

	_, error := sls.CheckUserShoppingListExistOrNot(userGUID, shoppingListGUID)

	if error != nil {
		return error
	}

	_, error = sls.OccasionService.CheckOccassionExistOrNot(updateData.OccasionGUID)

	if error != nil {
		return error
	}

	if updateData.Name != "" {
		error := sls.CheckUserShoppingListDuplicate(userGUID, updateData.Name, updateData.OccasionGUID)

		if error != nil {
			return error
		}
	}

	error = sls.ShoppingListRepository.Update(dbTransaction, userGUID, shoppingListGUID, updateData)

	if error != nil {
		return error
	}

	return nil
}

// DeleteUserShoppingListIncludingItemsAndImages function used to soft delete user shopping list including user shopping list items and
// shopping list item image inside the user shopping list.
func (sls *ShoppingListService) DeleteUserShoppingListIncludingItemsAndImages(dbTransaction *gorm.DB, userGUID string, shoppingListGUID string) *systems.ErrorData {
	_, error := sls.CheckUserShoppingListExistOrNot(userGUID, shoppingListGUID)

	if error != nil {
		return error
	}

	error = sls.ShoppingListRepository.Delete(dbTransaction, "guid", shoppingListGUID)

	if error != nil {
		return error
	}

	_, error = sls.ShoppingListItemService.DeleteAllShoppingListItemsInShoppingList(dbTransaction, userGUID, shoppingListGUID)

	if error != nil {
		return error
	}

	error = sls.ShoppingListItemImageService.DeleteImagesForShoppingList(dbTransaction, shoppingListGUID)

	if error != nil {
		return error
	}

	return nil
}

// GetUserShoppingLists function used to retrieve user shopping lists by user GUID.
func (sls *ShoppingListService) GetUserShoppingLists(userGUID string, relations string) ([]*ShoppingList, *systems.ErrorData) {
	shoppingLists := sls.ShoppingListRepository.GetByUserGUID(userGUID, relations)

	return shoppingLists, nil
}

// ViewShoppingListByGUID function used to retrieve shopping list using shopping list GUID.
func (sls *ShoppingListService) ViewShoppingListByGUID(shoppingListGUID string, relations string) *ShoppingList {
	shoppingList := sls.ShoppingListRepository.GetByGUID(shoppingListGUID, relations)

	return shoppingList
}

// CheckUserShoppingListDuplicate function used to check if user already has shopping list with the same name and occasion type.
func (sls *ShoppingListService) CheckUserShoppingListDuplicate(userGUID string, shoppingListName string,
	occasionGUID string) *systems.ErrorData {

	shoppingList := sls.ShoppingListRepository.GetByUserGUIDOccasionGUIDAndName(userGUID, shoppingListName, occasionGUID, "")

	if shoppingList.Name != "" {
		return Error.DuplicateValueErrors("Shopping List", "name", shoppingListName)
	}

	return nil
}

// CheckUserShoppingListExistOrNot function used to check user shopping list exist or not in database using user GUID and shopping list GUID.
func (sls *ShoppingListService) CheckUserShoppingListExistOrNot(userGUID string, shoppingListGUID string) (*ShoppingList, *systems.ErrorData) {
	shoppingList := sls.ShoppingListRepository.GetByGUIDAndUserGUID(shoppingListGUID, userGUID, "")

	if shoppingList.GUID == "" {
		return nil, Error.ResourceNotFoundError("Shopping List", "guid", shoppingListGUID)
	}

	return shoppingList, nil
}

// GetShoppingListIncludingDealCashbacks function used to retrieve shopping list by shopping list GUID including deal cashback
// related to the shopping list.
func (sls *ShoppingListService) GetShoppingListIncludingDealCashbacks(shoppingListGUID string, dealCashbackTransactionGUID string) *ShoppingList {
	shoppingListWithDealCashbacks := sls.ShoppingListRepository.GetByGUIDPreloadWithDealCashbacks(shoppingListGUID,
		dealCashbackTransactionGUID, "")

	return shoppingListWithDealCashbacks
}

// CreateSampleShoppingListsAndItemsForUser function used to create sample shopping list and shopping list item for user.
func (sls *ShoppingListService) CreateSampleShoppingListsAndItemsForUser(dbTransaction *gorm.DB, userGUID string) *systems.ErrorData {
	defaultShoppingLists := sls.DefaultShoppingListService.GetAllDefaultShoppingLists("")

	for _, defaultShoppingList := range defaultShoppingLists {
		userShoppingList := CreateShoppingList{
			OccasionGUID: defaultShoppingList.OccasionGUID,
			Name:         defaultShoppingList.Name,
		}

		shoppingList, error := sls.CreateUserShoppingList(dbTransaction, userGUID, userShoppingList)

		if error != nil {
			return error
		}

		error = sls.createSampleShoppingListItems(dbTransaction, userGUID, shoppingList.GUID)

		if error != nil {
			return error
		}
	}

	return nil
}

func (sls *ShoppingListService) createSampleShoppingListItems(dbTransaction *gorm.DB, userGUID string, shoppingListGUID string) *systems.ErrorData {
	defaultShoppingListitems := sls.DefaultShoppingListItemService.GetAllDefaultShoppingListItems()

	for _, defaultShoppingListItem := range defaultShoppingListitems {
		userShoppingListItem := CreateShoppingListItem{
			UserGUID:         userGUID,
			ShoppingListGUID: shoppingListGUID,
			Name:             defaultShoppingListItem.Name,
			Quantity:         defaultShoppingListItem.Quantity,
			Remark:           defaultShoppingListItem.Remark,
			AddedToCart:      defaultShoppingListItem.AddedToCart,
		}

		_, error := sls.ShoppingListItemService.CreateUserShoppingListItem(dbTransaction, userGUID, shoppingListGUID, userShoppingListItem)

		if error != nil {
			return error
		}
	}

	return nil
}
