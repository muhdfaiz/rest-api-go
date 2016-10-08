package v1

import "github.com/jinzhu/gorm"

type ShoppingListItemImageRepositoryInterface interface {
	GetByItemGUIDAndGUID(shoppingListItemGUID string, shoppingListItemImageGUID string, relations string) *ShoppingListItemImage
}

type ShoppingListItemImageRepository struct {
	DB *gorm.DB
}

// GetByItemGUIDAndGUID function used to retrieve shopping list item image by shopping
// list item GUID and shoping list item image GUID
func (sliir *ShoppingListItemImageRepository) GetByItemGUIDAndGUID(shoppingListItemGUID string, shoppingListItemImageGUID string, relations string) *ShoppingListItemImage {
	shoppingListItemImage := &ShoppingListItemImage{}

	DB := sliir.DB.Model(&ShoppingListItemImage{})

	if relations != "" {
		DB = LoadRelations(DB, relations)
	}

	DB.Where(&ShoppingListItemImage{GUID: shoppingListItemImageGUID, ShoppingListItemGUID: shoppingListItemGUID}).First(&shoppingListItemImage)

	return shoppingListItemImage
}
