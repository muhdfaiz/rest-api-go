package v1

import (
	"bitbucket.org/cliqers/shoppermate-api/systems"
	"github.com/fatih/structs"
	"github.com/jinzhu/gorm"
)

type ShoppingListFactoryInterface interface {
	Create(userGUID string, data CreateShoppingList) (*ShoppingList, *systems.ErrorData)
	Update(userGUID string, shoppingListGUID string, data UpdateShoppingList) *systems.ErrorData
	Delete(attribute string, value string) *systems.ErrorData
}

type ShoppingListFactory struct {
	DB *gorm.DB
}

func (slf *ShoppingListFactory) Create(userGUID string, data CreateShoppingList) (*ShoppingList, *systems.ErrorData) {
	shoppingList := &ShoppingList{
		GUID:         Helper.GenerateUUID(),
		UserGUID:     userGUID,
		Name:         data.Name,
		OccasionGUID: data.OccasionGUID,
	}

	result := slf.DB.Create(shoppingList)

	if result.Error != nil || result.RowsAffected == 0 {
		return nil, Error.InternalServerError(result.Error, systems.DatabaseError)
	}

	return result.Value.(*ShoppingList), nil
}

// Update function used to update shopping list
// Require device uuid. Must provide in url
func (slf *ShoppingListFactory) Update(userGUID string, shoppingListGUID string, data UpdateShoppingList) *systems.ErrorData {
	updateData := map[string]string{}

	for key, value := range structs.Map(data) {
		if value != "" {
			updateData[key] = value.(string)
		}
	}

	result := slf.DB.Model(&ShoppingList{}).Where(&ShoppingList{UserGUID: userGUID, GUID: shoppingListGUID}).Updates(updateData)

	if result.Error != nil {
		return Error.InternalServerError(result.Error, systems.DatabaseError)
	}

	return nil
}

func (slf *ShoppingListFactory) Delete(attribute string, value string) *systems.ErrorData {
	result := slf.DB.Where(attribute+" = ?", value).Delete(&ShoppingList{})

	if result.Error != nil {
		return Error.InternalServerError(result.Error, systems.DatabaseError)
	}

	return nil
}
