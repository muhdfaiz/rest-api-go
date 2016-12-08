package v1

import (
	"bitbucket.org/cliqers/shoppermate-api/systems"
	"github.com/jinzhu/gorm"
	"github.com/serenize/snaker"
)

type ShoppingListItemFactoryInterface interface {
	Create(data CreateShoppingListItem) (*ShoppingListItem, *systems.ErrorData)
	CreateForDeal(data CreateShoppingListItem) (*ShoppingListItem, *systems.ErrorData)
	UpdateByUserGUIDShoppingListGUIDAndShoppingListItemGUID(userGUID string, shoppingListGUID string, shoppingListItemGUID string,
		data map[string]interface{}) *systems.ErrorData
	UpdateByUserGUIDAndShoppingListGUID(userGUID string, shoppingListGUID string, data map[string]interface{}) *systems.ErrorData
	UpdateByUserGUIDAndDealGUID(userGUID string, dealGUID string, data map[string]interface{}) *systems.ErrorData
	UpdateByUserGUIDShoppingListGUIDAndDealGUID(userGUID string, shoppingListGUID string, dealGUID string,
		data map[string]interface{}) *systems.ErrorData
	DeleteByGUID(guid string) *systems.ErrorData
	DeleteByShoppingListGUID(shoppingListGUID string) *systems.ErrorData
	DeleteByUserGUID(userGUID string) *systems.ErrorData
	DeleteItemsHasBeenAddedToCartByUserGUIDAndShoppingListGUID(userGUID string, shoppingListGUID string) *systems.ErrorData
	DeleteItemsHasNotBeenAddedToCartByUserGUIDAndShoppingListGUID(userGUID string, shoppingListGUID string) *systems.ErrorData
}

// ShoppingListItemFactory contain functions to create, update and delete shopping list item
type ShoppingListItemFactory struct {
	DB                              *gorm.DB
	ItemRepository                  ItemRepositoryInterface
	DealRepository                  DealRepositoryInterface
	ShoppingListItemImageFactory    ShoppingListItemImageFactoryInterface
	ShoppingListItemImageRepository ShoppingListItemImageRepositoryInterface
	ItemCategoryRepository          ItemCategoryRepositoryInterface
	ItemSubCategoryRepository       ItemSubCategoryRepositoryInterface
}

// Create function used to create user shopping list item
func (slif *ShoppingListItemFactory) Create(data CreateShoppingListItem) (*ShoppingListItem, *systems.ErrorData) {
	item := slif.ItemRepository.GetByName(data.Name, "")

	shoppingListItemCategory := "Others"
	shoppingListItemSubCategory := "Others"

	if item.GUID != "" {
		itemCategory := slif.ItemCategoryRepository.GetByID(item.CategoryID)

		shoppingListItemCategory = itemCategory.Name

		itemSubCategory := slif.ItemSubCategoryRepository.GetByID(item.SubcategoryID)

		shoppingListItemSubCategory = itemSubCategory.Name
	}

	shoppingListItem := &ShoppingListItem{
		GUID:             Helper.GenerateUUID(),
		UserGUID:         data.UserGUID,
		ShoppingListGUID: data.ShoppingListGUID,
		Name:             data.Name,
		Category:         shoppingListItemCategory,
		SubCategory:      shoppingListItemSubCategory,
		Quantity:         data.Quantity,
	}

	result := slif.DB.Create(shoppingListItem)

	if result.Error != nil || result.RowsAffected == 0 {
		return nil, Error.InternalServerError(result.Error, systems.DatabaseError)
	}

	return result.Value.(*ShoppingListItem), nil
}

// CreateForDeal function used to create user shopping list item
func (slif *ShoppingListItemFactory) CreateForDeal(data CreateShoppingListItem) (*ShoppingListItem, *systems.ErrorData) {
	item := slif.ItemRepository.GetByName(data.Name, "")

	itemCategory := slif.ItemCategoryRepository.GetByID(item.CategoryID)
	shoppingListItemCategory := itemCategory.Name

	itemSubCategory := slif.ItemSubCategoryRepository.GetByID(item.SubcategoryID)
	shoppingListItemSubCategory := itemSubCategory.Name

	deal := slif.DealRepository.GetDealByGUID(data.DealGUID)

	shoppingListItem := &ShoppingListItem{
		GUID:             Helper.GenerateUUID(),
		UserGUID:         data.UserGUID,
		ShoppingListGUID: data.ShoppingListGUID,
		Name:             data.Name,
		Category:         shoppingListItemCategory,
		SubCategory:      shoppingListItemSubCategory,
		Quantity:         data.Quantity,
		AddedFromDeal:    data.AddedFromDeal,
		DealGUID:         &data.DealGUID,
		CashbackAmount:   &deal.CashbackAmount,
	}

	result := slif.DB.Create(shoppingListItem)

	if result.Error != nil || result.RowsAffected == 0 {
		return nil, Error.InternalServerError(result.Error, systems.DatabaseError)
	}

	return result.Value.(*ShoppingListItem), nil
}

// UpdateByUserGUIDShoppingListGUIDAndShoppingListItemGUID function used to update device data
// Require device uuid. Must provide in url
func (slif *ShoppingListItemFactory) UpdateByUserGUIDShoppingListGUIDAndShoppingListItemGUID(userGUID string, shoppingListGUID string, shoppingListItemGUID string, data map[string]interface{}) *systems.ErrorData {
	updateData := map[string]interface{}{}

	for key, value := range data {
		if data, ok := value.(string); ok && value.(string) != "" {
			updateData[snaker.CamelToSnake(key)] = data
		}

		if data, ok := value.(int); ok {
			updateData[snaker.CamelToSnake(key)] = data
		}
	}

	result := slif.DB.Model(&ShoppingListItem{}).Where(&ShoppingListItem{GUID: shoppingListItemGUID, ShoppingListGUID: shoppingListGUID, UserGUID: userGUID}).
		Updates(updateData)

	if result.Error != nil {
		return Error.InternalServerError(result.Error, systems.DatabaseError)
	}

	return nil
}

// UpdateByUserGUIDAndShoppingListGUID function used to update user shopping list item data by user GUID and shopping list GUID
func (slif *ShoppingListItemFactory) UpdateByUserGUIDAndShoppingListGUID(userGUID string, shoppingListGUID string, data map[string]interface{}) *systems.ErrorData {
	updateData := map[string]interface{}{}

	for key, value := range data {
		if data, ok := value.(string); ok && value.(string) != "" {
			updateData[snaker.CamelToSnake(key)] = data
		}

		if data, ok := value.(int); ok {
			updateData[snaker.CamelToSnake(key)] = data
		}
	}

	result := slif.DB.Model(&ShoppingListItem{}).Where(&ShoppingListItem{ShoppingListGUID: shoppingListGUID, UserGUID: userGUID}).
		Updates(updateData)

	if result.Error != nil {
		return Error.InternalServerError(result.Error, systems.DatabaseError)
	}

	return nil
}

// UpdateByUserGUIDAndDealGUID function used to update user shopping list item by user GUID and deal GUID
func (slif *ShoppingListItemFactory) UpdateByUserGUIDAndDealGUID(userGUID string, dealGUID string, data map[string]interface{}) *systems.ErrorData {
	updateData := map[string]interface{}{}

	for key, value := range data {
		if data, ok := value.(string); ok && value.(string) != "" {
			updateData[snaker.CamelToSnake(key)] = data
		}

		if data, ok := value.(int); ok {
			updateData[snaker.CamelToSnake(key)] = data
		}
	}

	result := slif.DB.Model(&ShoppingListItem{}).Where(&ShoppingListItem{UserGUID: userGUID, DealGUID: &dealGUID}).
		Updates(updateData)

	if result.Error != nil {
		return Error.InternalServerError(result.Error, systems.DatabaseError)
	}

	return nil
}

// UpdateByUserGUIDShoppingListGUIDAndDealGUID function used to update user shopping list item by user GUID and deal GUID
func (slif *ShoppingListItemFactory) UpdateByUserGUIDShoppingListGUIDAndDealGUID(userGUID string, shoppingListGUID string, dealGUID string,
	data map[string]interface{}) *systems.ErrorData {

	updateData := map[string]interface{}{}

	for key, value := range data {
		if data, ok := value.(string); ok && value.(string) != "" {
			updateData[snaker.CamelToSnake(key)] = data
		}

		if data, ok := value.(int); ok {
			updateData[snaker.CamelToSnake(key)] = data
		}
	}

	result := slif.DB.Model(&ShoppingListItem{}).Where(&ShoppingListItem{UserGUID: userGUID, ShoppingListGUID: shoppingListGUID, DealGUID: &dealGUID}).
		Updates(updateData)

	if result.Error != nil {
		return Error.InternalServerError(result.Error, systems.DatabaseError)
	}

	return nil
}

// DeleteByGUID function used to soft delete shopping list via GUID including the relationship from database
func (slif *ShoppingListItemFactory) DeleteByGUID(guid string) *systems.ErrorData {
	deleteShoppingListItem := slif.DB.Where("guid = ?", guid).Delete(&ShoppingListItem{})

	if deleteShoppingListItem.Error != nil {
		return Error.InternalServerError(deleteShoppingListItem.Error, systems.DatabaseError)
	}

	itemImages := slif.ShoppingListItemImageRepository.GetByItemGUID(guid, "")

	if len(itemImages) > 0 {
		imageURLs := make([]string, len(itemImages))

		for key, itemImage := range itemImages {
			imageURLs[key] = itemImage.URL
		}

		err := slif.ShoppingListItemImageFactory.Delete("shopping_list_item_guid", []string{guid}, imageURLs)

		if err != nil {
			return Error.InternalServerError(err.Error, systems.DatabaseError)
		}
	}

	return nil
}

// DeleteByShoppingListGUID function used to soft delete shopping list via shopping list GUID including the relationship from database
func (slif *ShoppingListItemFactory) DeleteByShoppingListGUID(shoppingListGUID string) *systems.ErrorData {
	// Delete shopping list item by shopping list GUID
	deleteShoppingListItem := slif.DB.Where("shopping_list_guid = ?", shoppingListGUID).Delete(&ShoppingListItem{})

	if deleteShoppingListItem.Error != nil {
		return Error.InternalServerError(deleteShoppingListItem.Error, systems.DatabaseError)
	}

	// Retrieve Shopping List Item Images
	itemImages := slif.ShoppingListItemImageRepository.GetByShoppingListGUID(shoppingListGUID, "")

	if len(itemImages) > 0 {
		imageURLs := make([]string, len(itemImages))

		for key, itemImage := range itemImages {
			imageURLs[key] = itemImage.URL
		}

		err := slif.ShoppingListItemImageFactory.Delete("shopping_list_guid", []string{shoppingListGUID}, imageURLs)

		if err != nil {
			return Error.InternalServerError(err.Error, systems.DatabaseError)
		}
	}

	return nil
}

// DeleteByUserGUID function used to soft delete all shopping list item via user GUID including the relationship from database
func (slif *ShoppingListItemFactory) DeleteByUserGUID(userGUID string) *systems.ErrorData {
	// Delete shopping list item by shopping list GUID
	deleteShoppingListItem := slif.DB.Where("user_guid = ?", userGUID).Delete(&ShoppingListItem{})

	if deleteShoppingListItem.Error != nil {
		return Error.InternalServerError(deleteShoppingListItem.Error, systems.DatabaseError)
	}

	// Retrieve Shopping List Item Images
	itemImages := slif.ShoppingListItemImageRepository.GetByUserGUID(userGUID, "")

	if len(itemImages) > 0 {
		imageURLs := make([]string, len(itemImages))

		for key, itemImage := range itemImages {
			imageURLs[key] = itemImage.URL
		}

		err := slif.ShoppingListItemImageFactory.Delete("user_guid", []string{userGUID}, imageURLs)

		if err != nil {
			return Error.InternalServerError(err.Error, systems.DatabaseError)
		}
	}

	return nil
}

// DeleteItemsHasBeenAddedToCartByUserGUIDAndShoppingListGUID function used to soft delete all shopping list item via user GUID
// and items those has been added to cart including the relationship from database
func (slif *ShoppingListItemFactory) DeleteItemsHasBeenAddedToCartByUserGUIDAndShoppingListGUID(userGUID string, shoppingListGUID string) *systems.ErrorData {
	userShoppingListItemsHasBeenAddedToCart := []*ShoppingListItem{}

	// Retrieve shopping list item those has been added to cart by user
	slif.DB.Where("user_guid = ? AND added_to_cart = ?", userGUID, 1).Find(&userShoppingListItemsHasBeenAddedToCart)

	// Delete shopping list item by user_guid and itemsadded to cart
	deleteShoppingListItem := slif.DB.Where("user_guid = ? AND added_to_cart = ? AND shopping_list_guid = ?", userGUID, 1, shoppingListGUID).Delete(&ShoppingListItem{})

	if deleteShoppingListItem.Error != nil {
		return Error.InternalServerError(deleteShoppingListItem.Error, systems.DatabaseError)
	}

	for _, userShoppingListItemHasBeenAddedToCart := range userShoppingListItemsHasBeenAddedToCart {
		// Retrieve Shopping List Item Images
		itemImages := slif.ShoppingListItemImageRepository.GetByItemGUID(userShoppingListItemHasBeenAddedToCart.GUID, "")

		if len(itemImages) > 0 {
			imageURLs := make([]string, len(itemImages))

			for key, itemImage := range itemImages {
				imageURLs[key] = itemImage.URL
			}

			err := slif.ShoppingListItemImageFactory.Delete("shopping_list_item_guid", []string{userShoppingListItemHasBeenAddedToCart.GUID}, imageURLs)

			if err != nil {
				return Error.InternalServerError(err.Error, systems.DatabaseError)
			}
		}
	}

	return nil
}

// DeleteItemsHasNotBeenAddedToCartByUserGUIDAndShoppingListGUID function used to soft delete all shopping list item via user GUID
// and items those has been added to cart including the relationship from database
func (slif *ShoppingListItemFactory) DeleteItemsHasNotBeenAddedToCartByUserGUIDAndShoppingListGUID(userGUID string, shoppingListGUID string) *systems.ErrorData {
	userShoppingListItemsHasBeenAddedToCart := []*ShoppingListItem{}

	// Retrieve shopping list item those has been added to cart by user
	slif.DB.Where("user_guid = ? AND added_to_cart = ?", userGUID, 0).Find(&userShoppingListItemsHasBeenAddedToCart)

	// Delete shopping list item by user_guid and itemsadded to cart
	deleteShoppingListItem := slif.DB.Where("user_guid = ? AND added_to_cart = ? AND shopping_list_guid", userGUID, 0, shoppingListGUID).Delete(&ShoppingListItem{})

	if deleteShoppingListItem.Error != nil {
		return Error.InternalServerError(deleteShoppingListItem.Error, systems.DatabaseError)
	}

	for _, userShoppingListItemHasBeenAddedToCart := range userShoppingListItemsHasBeenAddedToCart {
		// Retrieve Shopping List Item Images
		itemImages := slif.ShoppingListItemImageRepository.GetByItemGUID(userShoppingListItemHasBeenAddedToCart.GUID, "")

		if len(itemImages) > 0 {
			imageURLs := make([]string, len(itemImages))

			for key, itemImage := range itemImages {
				imageURLs[key] = itemImage.URL
			}

			err := slif.ShoppingListItemImageFactory.Delete("shopping_list_item_guid", []string{userShoppingListItemHasBeenAddedToCart.GUID}, imageURLs)

			if err != nil {
				return Error.InternalServerError(err.Error, systems.DatabaseError)
			}
		}
	}

	return nil
}
