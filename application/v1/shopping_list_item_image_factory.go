package v1

import (
	"fmt"

	"bitbucket.org/cliqers/shoppermate-api/systems"
	"github.com/jinzhu/gorm"
)

type ShoppingListItemImageFactoryInterface interface {
	Create(userGUID string, shoppingListGUID string, shoppingListItemGUID string, images []map[string]string) ([]*ShoppingListItemImage, *systems.ErrorData)
	Delete(attribute string, values []string, imageURL []string) *systems.ErrorData
}

// ShoppingListItemImageFactory used to handle create, update and delete shopping list item images
type ShoppingListItemImageFactory struct {
	DB                           *gorm.DB
	ShoppingListItemImageService ShoppingListItemImageServiceInterface
}

// Create function used to store shopping list item images in database
func (sliif *ShoppingListItemImageFactory) Create(userGUID string, shoppingListGUID string, shoppingListItemGUID string,
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

		result := sliif.DB.Create(shoppingListItemImage)

		if result.Error != nil || result.RowsAffected == 0 {
			return nil, Error.InternalServerError(result.Error, systems.DatabaseError)
		}

		createdImages[key] = result.Value.(*ShoppingListItemImage)
	}

	return createdImages, nil
}

// Delete function used to delete shopping list item image from database
func (sliif *ShoppingListItemImageFactory) Delete(attribute string, values []string, imageURLs []string) *systems.ErrorData {
	for _, value := range values {
		deleteShoppingListItem := sliif.DB.Where(attribute+" = ?", value).Delete(&ShoppingListItemImage{})

		if deleteShoppingListItem.Error != nil {
			return Error.InternalServerError(deleteShoppingListItem.Error, systems.DatabaseError)
		}
	}
	fmt.Println(imageURLs)
	// Delete shopping list item image from Amazon S3
	err := sliif.ShoppingListItemImageService.DeleteImages(imageURLs)

	if err != nil {
		return err
	}

	return nil
}
