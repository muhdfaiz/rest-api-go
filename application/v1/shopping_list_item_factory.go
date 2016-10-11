package v1

import (
	"fmt"
	"strconv"

	"bitbucket.org/cliqers/shoppermate-api/systems"
	"github.com/jinzhu/gorm"
)

type ShoppingListItemFactoryInterface interface {
	Create(data CreateShoppingListItem) (*ShoppingListItem, *systems.ErrorData)
	Update(userGUID string, shoppingListGUID string, shoppingListItemGUID string, data map[string]interface{}) *systems.ErrorData
	DeleteByGUID(guid string) *systems.ErrorData
	DeleteByShoppingListGUID(shoppingListGUID string) *systems.ErrorData
}

// ShoppingListItemFactory contain functions to create, update and delete shopping list item
type ShoppingListItemFactory struct {
	DB                              *gorm.DB
	ShoppingListItemImageFactory    ShoppingListItemImageFactoryInterface
	ShoppingListItemImageRepository ShoppingListItemImageRepositoryInterface
}

// Create function used to create user shopping list item
func (slif *ShoppingListItemFactory) Create(data CreateShoppingListItem) (*ShoppingListItem, *systems.ErrorData) {
	quantity, _ := strconv.Atoi(data.Quantity)
	shoppingListItem := &ShoppingListItem{
		GUID:             Helper.GenerateUUID(),
		UserGUID:         data.UserGUID,
		ShoppingListGUID: data.ShoppingListGUID,
		Name:             data.Name,
		Quantity:         quantity,
	}

	result := slif.DB.Create(shoppingListItem)

	if result.Error != nil || result.RowsAffected == 0 {
		return nil, Error.InternalServerError(result.Error, systems.DatabaseError)
	}

	return result.Value.(*ShoppingListItem), nil
}

// Update function used to update device data
// Require device uuid. Must provide in url
func (slif *ShoppingListItemFactory) Update(userGUID string, shoppingListGUID string, shoppingListItemGUID string, data map[string]interface{}) *systems.ErrorData {
	updateData := map[string]interface{}{}

	for key, value := range data {
		if data, ok := value.(string); ok && value.(string) != "" {
			updateData[key] = data
		}
		if data, ok := value.(int); ok && value.(int) != 0 {
			updateData[key] = data
		}
	}
	fmt.Println(updateData)
	result := slif.DB.Model(&ShoppingListItem{}).Where(&ShoppingListItem{GUID: shoppingListItemGUID, ShoppingListGUID: shoppingListGUID, UserGUID: userGUID}).
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

	imageURLs := make([]string, len(itemImages))

	for key, itemImage := range itemImages {
		imageURLs[key] = itemImage.URL
	}

	err := slif.ShoppingListItemImageFactory.Delete("shopping_list_item_guid", []string{guid}, imageURLs)

	if err != nil {
		return Error.InternalServerError(err.Error, systems.DatabaseError)
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
