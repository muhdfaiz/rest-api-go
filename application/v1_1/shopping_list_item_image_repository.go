package v1_1

import (
	"bitbucket.org/cliqers/shoppermate-api/systems"
	"github.com/jinzhu/gorm"
)

// ShoppingListItemImageRepository will handle all CRUD task for shopping list item image resource.
type ShoppingListItemImageRepository struct {
	BaseRepository
	DB *gorm.DB
}

// Create function used to store shopping list item images in database.
func (sliir *ShoppingListItemImageRepository) Create(dbTransaction *gorm.DB, userGUID string, shoppingListGUID string, shoppingListItemGUID string,
	images []map[string]string) ([]*ShoppingListItemImage, *systems.ErrorData) {

	createdImages := make([]*ShoppingListItemImage, len(images))

	for key, image := range images {

		shoppingListItemImage := &ShoppingListItemImage{
			GUID:                 Helper.GenerateUUID(),
			UserGUID:             userGUID,
			ShoppingListGUID:     shoppingListGUID,
			ShoppingListItemGUID: shoppingListItemGUID,
			URL:                  image["path"],
		}

		result := dbTransaction.Create(shoppingListItemImage)

		if result.Error != nil || result.RowsAffected == 0 {
			return nil, Error.InternalServerError(result.Error, systems.DatabaseError)
		}

		createdImages[key] = result.Value.(*ShoppingListItemImage)
	}

	return createdImages, nil
}

// Delete function used to soft delete shopping list item image from database.
func (sliir *ShoppingListItemImageRepository) Delete(dbTransaction *gorm.DB, attribute string, value string) *systems.ErrorData {

	deleteShoppingListItemImage := dbTransaction.Where(attribute+" = ?", value).Delete(&ShoppingListItemImage{})

	if deleteShoppingListItemImage.Error != nil {
		return Error.InternalServerError(deleteShoppingListItemImage.Error, systems.DatabaseError)
	}

	return nil
}

// GetByUserGUIDAndShoppingListGUIDAndItemGUIDAndImageGUID function used to retrieve shopping list item image by user GUID,
// shopping list GUID, shopping list item GUID and shopping list item image GUID.
func (sliir *ShoppingListItemImageRepository) GetByUserGUIDAndShoppingListGUIDAndItemGUIDAndImageGUID(userGUID string, shoppingListGUID string,
	shoppingListItemGUID string, shoppingListItemImageGUID string, relations string) *ShoppingListItemImage {

	shoppingListItemImage := &ShoppingListItemImage{}

	DB := sliir.DB.Model(&ShoppingListItemImage{})

	if relations != "" {
		DB = sliir.LoadRelations(DB, relations)
	}

	DB.Where(&ShoppingListItemImage{GUID: shoppingListItemImageGUID, ShoppingListItemGUID: shoppingListItemGUID}).First(&shoppingListItemImage)

	return shoppingListItemImage
}

// GetByItemGUID function used to retrieve shopping list item images using shopping list item GUID.
func (sliir *ShoppingListItemImageRepository) GetByItemGUID(shoppingListItemGUID string, relations string) []*ShoppingListItemImage {
	shoppingListItemImage := []*ShoppingListItemImage{}

	DB := sliir.DB.Model(&ShoppingListItemImage{})

	if relations != "" {
		DB = sliir.LoadRelations(DB, relations)
	}

	DB.Where(&ShoppingListItemImage{ShoppingListItemGUID: shoppingListItemGUID}).Find(&shoppingListItemImage)

	return shoppingListItemImage
}

// GetByShoppingListGUID function used to retrieve shopping list item images using shopping list GUID.
func (sliir *ShoppingListItemImageRepository) GetByShoppingListGUID(shoppingListGUID string, relations string) []*ShoppingListItemImage {
	shoppingListItemImage := []*ShoppingListItemImage{}

	DB := sliir.DB.Model(&ShoppingListItemImage{})

	if relations != "" {
		DB = sliir.LoadRelations(DB, relations)
	}

	DB.Where(&ShoppingListItemImage{ShoppingListGUID: shoppingListGUID}).Find(&shoppingListItemImage)

	return shoppingListItemImage
}

// GetByUserGUID function used to retrieve shopping list item images using User GUID.
func (sliir *ShoppingListItemImageRepository) GetByUserGUID(userGUID string, relations string) []*ShoppingListItemImage {
	shoppingListItemImage := []*ShoppingListItemImage{}

	DB := sliir.DB.Model(&ShoppingListItemImage{})

	if relations != "" {
		DB = sliir.LoadRelations(DB, relations)
	}

	DB.Where(&ShoppingListItemImage{UserGUID: userGUID}).Find(&shoppingListItemImage)

	return shoppingListItemImage
}
