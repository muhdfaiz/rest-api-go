package v1

import (
	"bitbucket.org/cliqers/shoppermate-api/systems"
	"github.com/jinzhu/gorm"
	"github.com/serenize/snaker"
)

// ShoppingListItemRepository will handle all CRUD functions related to Shopping List resource.
type ShoppingListItemRepository struct {
	DB *gorm.DB
}

// Create function used to create user shopping list item
func (slir *ShoppingListItemRepository) Create(dbTransaction *gorm.DB, data CreateShoppingListItem) (*ShoppingListItem, *systems.ErrorData) {
	dealGUID := data.DealGUID

	cashbackAmount := data.CashbackAmount

	shoppingListItem := &ShoppingListItem{
		GUID:             Helper.GenerateUUID(),
		UserGUID:         data.UserGUID,
		ShoppingListGUID: data.ShoppingListGUID,
		Name:             data.Name,
		Category:         data.Category,
		SubCategory:      data.SubCategory,
		Quantity:         data.Quantity,
		AddedToCart:      data.AddedToCart,
		AddedFromDeal:    data.AddedFromDeal,
		DealGUID:         &dealGUID,
		CashbackAmount:   &cashbackAmount,
	}

	result := dbTransaction.Create(shoppingListItem)

	if result.Error != nil || result.RowsAffected == 0 {
		return nil, Error.InternalServerError(result.Error, systems.DatabaseError)
	}

	return result.Value.(*ShoppingListItem), nil
}

// UpdateByUserGUIDShoppingListGUIDAndShoppingListItemGUID function used to update device data
// Require device uuid. Must provide in url
func (slir *ShoppingListItemRepository) UpdateByUserGUIDShoppingListGUIDAndShoppingListItemGUID(dbTransaction *gorm.DB, userGUID string,
	shoppingListGUID string, shoppingListItemGUID string, data map[string]interface{}) *systems.ErrorData {

	updateData := map[string]interface{}{}

	for key, value := range data {
		if data, ok := value.(string); ok && value.(string) != "" {
			updateData[snaker.CamelToSnake(key)] = data
		}

		if data, ok := value.(int); ok {
			if key == "Quantity" && data != 0 {
				updateData[snaker.CamelToSnake(key)] = data
			}

			if key != "Quantity" {
				updateData[snaker.CamelToSnake(key)] = data
			}
		}
	}

	result := dbTransaction.Model(&ShoppingListItem{}).Where(&ShoppingListItem{GUID: shoppingListItemGUID, ShoppingListGUID: shoppingListGUID, UserGUID: userGUID}).
		Updates(updateData)

	if result.Error != nil {
		return Error.InternalServerError(result.Error, systems.DatabaseError)
	}

	return nil
}

// UpdateByUserGUIDAndShoppingListGUID function used to update user shopping list item data by user GUID and shopping list GUID
func (slir *ShoppingListItemRepository) UpdateByUserGUIDAndShoppingListGUID(dbTransaction *gorm.DB, userGUID string, shoppingListGUID string,
	data map[string]interface{}) *systems.ErrorData {

	updateData := map[string]interface{}{}

	for key, value := range data {
		if data, ok := value.(string); ok && value.(string) != "" {
			updateData[snaker.CamelToSnake(key)] = data
		}

		if data, ok := value.(int); ok {
			if key == "Quantity" && data != 0 {
				updateData[snaker.CamelToSnake(key)] = data
			}

			if key != "Quantity" {
				updateData[snaker.CamelToSnake(key)] = data
			}
		}
	}

	result := dbTransaction.Model(&ShoppingListItem{}).Where(&ShoppingListItem{ShoppingListGUID: shoppingListGUID, UserGUID: userGUID}).
		Updates(updateData)

	if result.Error != nil {
		return Error.InternalServerError(result.Error, systems.DatabaseError)
	}

	return nil
}

// UpdateByUserGUIDAndDealGUID function used to update user shopping list item by user GUID and deal GUID
func (slir *ShoppingListItemRepository) UpdateByUserGUIDAndDealGUID(dbTransaction *gorm.DB, userGUID string, dealGUID string,
	data map[string]interface{}) *systems.ErrorData {

	updateData := map[string]interface{}{}

	for key, value := range data {
		if data, ok := value.(string); ok && value.(string) != "" {
			updateData[snaker.CamelToSnake(key)] = data
		}

		if data, ok := value.(int); ok {
			if key == "Quantity" && data != 0 {
				updateData[snaker.CamelToSnake(key)] = data
			}

			if key != "Quantity" {
				updateData[snaker.CamelToSnake(key)] = data
			}
		}
	}

	result := dbTransaction.Model(&ShoppingListItem{}).Where(&ShoppingListItem{UserGUID: userGUID, DealGUID: &dealGUID}).
		Updates(updateData)

	if result.Error != nil {
		return Error.InternalServerError(result.Error, systems.DatabaseError)
	}

	return nil
}

// UpdateByUserGUIDShoppingListGUIDAndDealGUID function used to update user shopping list item by user GUID and deal GUID
func (slir *ShoppingListItemRepository) UpdateByUserGUIDShoppingListGUIDAndDealGUID(dbTransaction *gorm.DB, userGUID string, shoppingListGUID string,
	dealGUID string, data map[string]interface{}) *systems.ErrorData {

	updateData := map[string]interface{}{}

	for key, value := range data {
		if data, ok := value.(string); ok && value.(string) != "" {
			updateData[snaker.CamelToSnake(key)] = data
		}

		if data, ok := value.(int); ok {
			if key == "Quantity" && data != 0 {
				updateData[snaker.CamelToSnake(key)] = data
			}

			if key != "Quantity" {
				updateData[snaker.CamelToSnake(key)] = data
			}
		}
	}

	result := dbTransaction.Model(&ShoppingListItem{}).Where(&ShoppingListItem{UserGUID: userGUID, ShoppingListGUID: shoppingListGUID, DealGUID: &dealGUID}).
		Updates(updateData)

	if result.Error != nil {
		return Error.InternalServerError(result.Error, systems.DatabaseError)
	}

	return nil
}

// SetDealExpired function used to set deal expired on shopping list item when the deal already expired.
func (slir *ShoppingListItemRepository) SetDealExpired(dbTransaction *gorm.DB, dealGUID string) *systems.ErrorData {
	result := dbTransaction.Model(&ShoppingListItem{}).Where("deal_guid = ?", dealGUID).Select("deal_expired").Updates(map[string]interface{}{"deal_expired": 1})

	if result.Error != nil {
		return Error.InternalServerError(result.Error, systems.DatabaseError)
	}

	return nil
}

// DeleteByGUID function used to soft delete shopping list via GUID including the relationship from database
func (slir *ShoppingListItemRepository) DeleteByGUID(dbTransaction *gorm.DB, shoppingListItemGUID string) *systems.ErrorData {
	deleteShoppingListItem := dbTransaction.Where("guid = ?", shoppingListItemGUID).Delete(&ShoppingListItem{})

	if deleteShoppingListItem.Error != nil {
		return Error.InternalServerError(deleteShoppingListItem.Error, systems.DatabaseError)
	}

	return nil
}

// DeleteByShoppingListGUID function used to soft delete shopping list via shopping list GUID including the relationship from database
func (slir *ShoppingListItemRepository) DeleteByShoppingListGUID(dbTransaction *gorm.DB, shoppingListGUID string) *systems.ErrorData {
	deleteShoppingListItem := dbTransaction.Where("shopping_list_guid = ?", shoppingListGUID).Delete(&ShoppingListItem{})

	if deleteShoppingListItem.Error != nil {
		return Error.InternalServerError(deleteShoppingListItem.Error, systems.DatabaseError)
	}

	return nil
}

// DeleteByUserGUIDAndShoppingListGUID function used to soft delete all user shopping list item by user
// GUID and shopping list GUID including shopping list item images.
func (slir *ShoppingListItemRepository) DeleteByUserGUIDAndShoppingListGUID(dbTransaction *gorm.DB, userGUID string, shoppingListGUID string) *systems.ErrorData {
	deleteShoppingListItem := dbTransaction.Where(&ShoppingListItem{UserGUID: userGUID, ShoppingListGUID: shoppingListGUID}).Delete(&ShoppingListItem{})

	if deleteShoppingListItem.Error != nil {
		return Error.InternalServerError(deleteShoppingListItem.Error, systems.DatabaseError)
	}

	return nil
}

// DeleteByGUIDAndUserGUIDAndShoppingListGUID function used to soft delete user shopping list item by
// shopping list item GUID, user GUID and shopping list GUID including shopping list item images.
func (slir *ShoppingListItemRepository) DeleteByGUIDAndUserGUIDAndShoppingListGUID(dbTransaction *gorm.DB, shoppingListItemGUID string,
	userGUID string, shoppingListGUID string) *systems.ErrorData {

	deleteShoppingListItem := dbTransaction.Where(&ShoppingListItem{GUID: shoppingListItemGUID, UserGUID: userGUID, ShoppingListGUID: shoppingListGUID}).Delete(&ShoppingListItem{})

	if deleteShoppingListItem.Error != nil {
		return Error.InternalServerError(deleteShoppingListItem.Error, systems.DatabaseError)
	}

	return nil
}

// DeleteItemsHasBeenAddedToCartByUserGUIDAndShoppingListGUID function used to soft delete all of user shopping list item
// those has been added to cart including shopping list item images.
func (slir *ShoppingListItemRepository) DeleteItemsHasBeenAddedToCartByUserGUIDAndShoppingListGUID(dbTransaction *gorm.DB, userGUID string,
	shoppingListGUID string) *systems.ErrorData {

	userShoppingListItemsHasBeenAddedToCart := []*ShoppingListItem{}

	slir.DB.Where(&ShoppingListItem{UserGUID: userGUID, ShoppingListGUID: shoppingListGUID, AddedToCart: 1}).Find(&userShoppingListItemsHasBeenAddedToCart)

	deleteShoppingListItem := dbTransaction.Where(&ShoppingListItem{UserGUID: userGUID, ShoppingListGUID: shoppingListGUID, AddedToCart: 1}).Delete(&ShoppingListItem{})

	if deleteShoppingListItem.Error != nil {
		return Error.InternalServerError(deleteShoppingListItem.Error, systems.DatabaseError)
	}

	return nil
}

// DeleteItemsHasNotBeenAddedToCartByUserGUIDAndShoppingListGUID function used to soft delete all of user shopping list item
// those has not been added to cart including shopping list item images.
func (slir *ShoppingListItemRepository) DeleteItemsHasNotBeenAddedToCartByUserGUIDAndShoppingListGUID(dbTransaction *gorm.DB, userGUID string,
	shoppingListGUID string) *systems.ErrorData {

	userShoppingListItemsHasBeenAddedToCart := []*ShoppingListItem{}

	// Retrieve shopping list item those has been added to cart by user
	slir.DB.Where(&ShoppingListItem{UserGUID: userGUID, ShoppingListGUID: shoppingListGUID, AddedToCart: 0}).Find(&userShoppingListItemsHasBeenAddedToCart)

	// Delete shopping list item by user_guid and itemsadded to cart
	deleteShoppingListItem := dbTransaction.Where(&ShoppingListItem{UserGUID: userGUID, ShoppingListGUID: shoppingListGUID, AddedToCart: 0}).Delete(&ShoppingListItem{})

	if deleteShoppingListItem.Error != nil {
		return Error.InternalServerError(deleteShoppingListItem.Error, systems.DatabaseError)
	}

	return nil
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

// GetByGUIDUserGUIDAndShoppingListGUID function used to retrieve shopping list item by shopping list GUID and GUID
func (slir *ShoppingListItemRepository) GetByGUIDUserGUIDAndShoppingListGUID(userGUID string, shoppingListGUID string,
	shoppingListItemGUID string, relations string) *ShoppingListItem {

	shoppingListItem := &ShoppingListItem{}

	DB := slir.DB.Model(&ShoppingListItem{})

	if relations != "" {
		DB = LoadRelations(DB, relations)
	}

	DB.Where(&ShoppingListItem{GUID: shoppingListItemGUID, UserGUID: userGUID, ShoppingListGUID: shoppingListGUID}).First(&shoppingListItem)

	return shoppingListItem
}

// GetByUserGUIDAndShoppingListGUID function used to retrieve shopping list item by user GUID and shopping list GUID
func (slir *ShoppingListItemRepository) GetByUserGUIDAndShoppingListGUID(userGUID string, shoppingListGUID string, relations string) []*ShoppingListItem {
	shoppingListItem := []*ShoppingListItem{}

	DB := slir.DB.Model(&ShoppingListItem{})

	if relations != "" {
		DB = LoadRelations(DB, relations)
	}

	DB.Where(&ShoppingListItem{UserGUID: userGUID, ShoppingListGUID: shoppingListGUID}).Find(&shoppingListItem)

	return shoppingListItem
}

// GetByUserGUIDAndShoppingListGUIDAndSubCategory function used to retrieve user shopping list items by user GUID, shopping list GUID,
// and Subcategory Name
func (slir *ShoppingListItemRepository) GetByUserGUIDAndShoppingListGUIDAndSubCategory(userGUID string, shoppingListGUID string,
	subcategory string, relations string) []*ShoppingListItem {

	shoppingListItems := []*ShoppingListItem{}

	DB := slir.DB.Model(&ShoppingListItem{})

	if relations != "" {
		DB = LoadRelations(DB, relations)
	}

	DB.Where("user_guid = ? AND shopping_list_guid = ? AND sub_category = ?", userGUID, shoppingListGUID,
		subcategory).Find(&shoppingListItems)

	return shoppingListItems
}

// GetByUserGUIDAndShoppingListGUIDAndAddedToCartAndSubCategory function used to retrieve user shopping list items by user GUID, shopping list GUID,
// added to cart and Subcategory Name
func (slir *ShoppingListItemRepository) GetByUserGUIDAndShoppingListGUIDAndAddedToCartAndSubCategory(userGUID string, shoppingListGUID string,
	addedToCart int, subcategory string, relations string) []*ShoppingListItem {

	shoppingListItems := []*ShoppingListItem{}

	DB := slir.DB.Model(&ShoppingListItem{})

	if relations != "" {
		DB = LoadRelations(DB, relations)
	}

	DB.Where("user_guid = ? AND shopping_list_guid = ? AND added_to_cart = ? AND sub_category = ?", userGUID, shoppingListGUID, addedToCart,
		subcategory).Find(&shoppingListItems)

	return shoppingListItems
}

// GetUniqueSubCategoryFromAllUserShoppingListItem function used to retrieve unique sub category from all user shopping list items.
func (slir *ShoppingListItemRepository) GetUniqueSubCategoryFromAllUserShoppingListItem(userGUID string, shoppingListGUID string) []*ShoppingListItem {

	shoppingListItemsGroupBySubCategory := []*ShoppingListItem{}

	slir.DB.Where("user_guid = ? AND shopping_list_guid = ?", userGUID, shoppingListGUID).Group("sub_category").
		Find(&shoppingListItemsGroupBySubCategory)

	return shoppingListItemsGroupBySubCategory
}

// GetUniqueSubCategoryFromUserShoppingListItem function used to retrieve unique sub category from user shopping list items those
// has been added to cart or those has not been added to cart.
func (slir *ShoppingListItemRepository) GetUniqueSubCategoryFromUserShoppingListItem(userGUID string, shoppingListGUID string,
	addedToCart int) []*ShoppingListItem {

	shoppingListItemsGroupBySubCategory := []*ShoppingListItem{}

	slir.DB.Where("user_guid = ? AND shopping_list_guid = ? AND added_to_cart = ?", userGUID, shoppingListGUID, addedToCart).Group("sub_category").
		Find(&shoppingListItemsGroupBySubCategory)

	return shoppingListItemsGroupBySubCategory
}
