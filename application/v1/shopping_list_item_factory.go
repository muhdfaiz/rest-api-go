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
}

type ShoppingListItemFactory struct {
	DB *gorm.DB
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
