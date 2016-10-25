package v1

import "github.com/jinzhu/gorm"

type ShoppingListItemImageRepositoryInterface interface {
	GetByItemGUIDAndGUID(shoppingListItemGUID string, shoppingListItemImageGUID string, relations string) *ShoppingListItemImage
	GetByItemGUID(shoppingListItemGUID string, relations string) []*ShoppingListItemImage
	GetByShoppingListGUID(shoppingListGUID string, relations string) []*ShoppingListItemImage
	GetByUserGUID(useerGUID string, relations string) []*ShoppingListItemImage
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

// GetByItemGUID function used to retrieve shopping list item images using shopping list item GUID
func (sliir *ShoppingListItemImageRepository) GetByItemGUID(shoppingListItemGUID string, relations string) []*ShoppingListItemImage {
	shoppingListItemImage := []*ShoppingListItemImage{}

	DB := sliir.DB.Model(&ShoppingListItemImage{})

	if relations != "" {
		DB = LoadRelations(DB, relations)
	}

	DB.Where(&ShoppingListItemImage{ShoppingListItemGUID: shoppingListItemGUID}).Find(&shoppingListItemImage)

	return shoppingListItemImage
}

// GetByShoppingListGUID function used to retrieve shopping list item images using shopping list GUID
func (sliir *ShoppingListItemImageRepository) GetByShoppingListGUID(shoppingListGUID string, relations string) []*ShoppingListItemImage {
	shoppingListItemImage := []*ShoppingListItemImage{}

	DB := sliir.DB.Model(&ShoppingListItemImage{})

	if relations != "" {
		DB = LoadRelations(DB, relations)
	}

	DB.Where(&ShoppingListItemImage{ShoppingListGUID: shoppingListGUID}).Find(&shoppingListItemImage)

	return shoppingListItemImage
}

// GetByUserGUID function used to retrieve shopping list item images using User GUID
func (sliir *ShoppingListItemImageRepository) GetByUserGUID(userGUID string, relations string) []*ShoppingListItemImage {
	shoppingListItemImage := []*ShoppingListItemImage{}

	DB := sliir.DB.Model(&ShoppingListItemImage{})

	if relations != "" {
		DB = LoadRelations(DB, relations)
	}

	DB.Where(&ShoppingListItemImage{UserGUID: userGUID}).Find(&shoppingListItemImage)

	return shoppingListItemImage
}
