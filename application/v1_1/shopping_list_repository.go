package v1_1

import (
	"bitbucket.org/cliqers/shoppermate-api/systems"
	"github.com/fatih/structs"
	"github.com/jinzhu/gorm"
)

// ShoppingListRepository used to retrieve user shopping list.
type ShoppingListRepository struct {
	DB *gorm.DB
}

// Create function used to create user shopping list.
func (slr *ShoppingListRepository) Create(dbTransaction *gorm.DB, userGUID string, data CreateShoppingList) (*ShoppingList, *systems.ErrorData) {
	shoppingList := &ShoppingList{
		GUID:         Helper.GenerateUUID(),
		UserGUID:     userGUID,
		Name:         data.Name,
		OccasionGUID: data.OccasionGUID,
	}

	result := dbTransaction.Create(shoppingList)

	if result.Error != nil || result.RowsAffected == 0 {
		return nil, Error.InternalServerError(result.Error, systems.DatabaseError)
	}

	return result.Value.(*ShoppingList), nil
}

// Update function used to update shopping list
// Require device uuid. Must provide in url
func (slr *ShoppingListRepository) Update(dbTransaction *gorm.DB, userGUID string, shoppingListGUID string, data UpdateShoppingList) *systems.ErrorData {
	updateData := map[string]string{}

	for key, value := range structs.Map(data) {
		if value != "" {
			updateData[key] = value.(string)
		}
	}

	result := dbTransaction.Model(&ShoppingList{}).Where(&ShoppingList{UserGUID: userGUID, GUID: shoppingListGUID}).Updates(updateData)

	if result.Error != nil {
		return Error.InternalServerError(result.Error, systems.DatabaseError)
	}

	return nil
}

// Delete function used to soft delete shopping list
func (slr *ShoppingListRepository) Delete(dbTransaction *gorm.DB, attribute string, value string) *systems.ErrorData {
	result := dbTransaction.Where(attribute+" = ?", value).Delete(&ShoppingList{})

	if result.Error != nil {
		return Error.InternalServerError(result.Error, systems.DatabaseError)
	}

	return nil
}

// GetByUserGUID function used to retrieve user shopping list by User GUID.
func (slr *ShoppingListRepository) GetByUserGUID(userGUID string, relations string) []*ShoppingList {
	shoppingLists := []*ShoppingList{}

	DB := slr.DB.Model(&ShoppingList{})

	if relations != "" {
		DB = LoadRelations(DB, relations)
	}

	DB.Where(&ShoppingList{UserGUID: userGUID}).Find(&shoppingLists)

	return shoppingLists
}

// GetByGUIDAndUserGUID function used to retrieve user shopping list by and Shopping List GUID.
func (slr *ShoppingListRepository) GetByGUIDAndUserGUID(GUID string, userGUID string, relations string) *ShoppingList {
	shoppingLists := &ShoppingList{}

	DB := slr.DB.Model(&ShoppingList{})

	if relations != "" {
		DB = LoadRelations(DB, relations)
	}

	DB.Where(&ShoppingList{GUID: GUID, UserGUID: userGUID}).First(&shoppingLists)

	return shoppingLists
}

// GetByGUID function used to retrieve user shopping list by Shopping List GUID.
func (slr *ShoppingListRepository) GetByGUID(GUID string, relations string) *ShoppingList {
	shoppingLists := &ShoppingList{}

	DB := slr.DB.Model(&ShoppingList{})

	if relations != "" {
		DB = LoadRelations(DB, relations)
	}

	DB.Where(&ShoppingList{GUID: GUID}).First(&shoppingLists)

	return shoppingLists
}

// GetUnscopedByGUID function used to retrieve user shopping list by Shopping List GUID including
// soft deleted records
func (slr *ShoppingListRepository) GetUnscopedByGUID(GUID string, relations string) *ShoppingList {
	shoppingLists := &ShoppingList{}

	DB := slr.DB.Unscoped()

	if relations != "" {
		DB = LoadRelations(DB, relations)
	}

	DB.Where(&ShoppingList{GUID: GUID}).First(&shoppingLists)

	return shoppingLists
}

// GetByGUIDPreloadWithDealCashbacks function used to retrieve user shopping list by and Shopping List GUID and load relation deal cashback
func (slr *ShoppingListRepository) GetByGUIDPreloadWithDealCashbacks(GUID string, dealCashbackTransactionGUID string, relations string) *ShoppingList {
	shoppingLists := &ShoppingList{}

	slr.DB.Model(&ShoppingList{}).Preload("Dealcashbacks", func(db *gorm.DB) *gorm.DB {
		return db.Where("deal_cashback_transaction_guid = ?", dealCashbackTransactionGUID)
	}).Preload("Dealcashbacks.Deals").Preload("Dealcashbacks.Dealcashbackstatus").Preload("Dealcashbacks.Dealcashbackstatus").Where(&ShoppingList{GUID: GUID}).First(&shoppingLists)

	return shoppingLists
}

// GetByUserGUIDOccasionGUIDAndName function used to retrieve user shopping list by User GUID and Shopping List Name.
func (slr *ShoppingListRepository) GetByUserGUIDOccasionGUIDAndName(userGUID string, name string, occasionGUID string, relations string) *ShoppingList {
	shoppingLists := &ShoppingList{}

	DB := slr.DB.Model(&ShoppingList{})

	if relations != "" {
		DB = LoadRelations(DB, relations)
	}

	DB.Where(&ShoppingList{UserGUID: userGUID, Name: name, OccasionGUID: occasionGUID}).First(&shoppingLists)

	return shoppingLists
}
