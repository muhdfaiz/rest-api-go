package v1

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

type ShoppingListRepositoryInterface interface {
	GetByUserGUID(userGUID string) []ShoppingList
	GetByGUID(GUID string) *ShoppingList
	GetByGUIDAndUserGUID(GUID string, userGUID string) *ShoppingList
	GetByUserGUIDOccasionGUIDAndName(userGUID string, name string, occasionGUID string) *ShoppingList
}

// ShoppingListRepository used to retrieve user shopping list.
type ShoppingListRepository struct {
	DB *gorm.DB
}

// GetByUserGUID function used to retrieve user shopping list by User GUID.
func (slr *ShoppingListRepository) GetByUserGUID(userGUID string) []ShoppingList {
	rows, _ := slr.DB.Model(&ShoppingList{}).Where(&ShoppingList{UserGUID: userGUID}).Rows()

	var shoppingLists []ShoppingList

	for rows.Next() {
		var shoppingList ShoppingList

		if err := rows.Scan(&shoppingList.ID, &shoppingList.GUID, &shoppingList.UserGUID, &shoppingList.OccasionGUID, &shoppingList.Name,
			&shoppingList.CreatedAt, &shoppingList.UpdatedAt, &shoppingList.DeletedAt); err != nil {
			fmt.Println(err)
		}

		occasion := slr.DB.Where(&Occasion{GUID: shoppingList.OccasionGUID}).First(&Occasion{})
		occasionData := occasion.Value.(*Occasion)

		shoppingList.Occasion.ID = occasionData.ID
		shoppingList.Occasion.GUID = occasionData.GUID
		shoppingList.Occasion.Name = occasionData.Name
		shoppingList.Occasion.Image = occasionData.Image
		shoppingList.Occasion.Slug = occasionData.Slug
		shoppingList.Occasion.CreatedAt = occasionData.CreatedAt
		shoppingList.Occasion.UpdatedAt = occasionData.UpdatedAt
		shoppingList.Occasion.UpdatedAt = occasionData.UpdatedAt

		shoppingLists = append(shoppingLists, shoppingList)
	}
	return shoppingLists
}

// GetByGUID function used to retrieve user shopping list by and Shopping List GUID.
func (slr *ShoppingListRepository) GetByGUIDAndUserGUID(GUID string, userGUID string) *ShoppingList {
	result := slr.DB.Preload("Occasion").Where(&ShoppingList{GUID: GUID, UserGUID: userGUID}).First(&ShoppingList{})

	if result.RowsAffected == 0 {
		return &ShoppingList{}
	}

	return result.Value.(*ShoppingList)
}

// GetByGUID function used to retrieve user shopping list by and Shopping List GUID.
func (slr *ShoppingListRepository) GetByGUID(GUID string) *ShoppingList {
	result := slr.DB.Preload("Occasion").Where(&ShoppingList{GUID: GUID}).First(&ShoppingList{})

	if result.RowsAffected == 0 {
		return &ShoppingList{}
	}

	return result.Value.(*ShoppingList)
}

// GetByUserGUIDOccasionGUIDAndName function used to retrieve user shopping list by User GUID and Shopping List Name.
func (slr *ShoppingListRepository) GetByUserGUIDOccasionGUIDAndName(userGUID string, name string, occasionGUID string) *ShoppingList {
	result := slr.DB.Where(&ShoppingList{UserGUID: userGUID, Name: name, OccasionGUID: occasionGUID}).First(&ShoppingList{})

	if result.RowsAffected == 0 {
		return &ShoppingList{}
	}

	return result.Value.(*ShoppingList)
}
