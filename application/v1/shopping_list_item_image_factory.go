package v1

import (
	"fmt"

	"bitbucket.org/cliqers/shoppermate-api/systems"
	"github.com/jinzhu/gorm"
)

type ShoppingListItemImageFactoryInterface interface {
	Create(shoppingListItemGUID string, images []map[string]string) ([]*ShoppingListItemImage, *systems.ErrorData)
	Delete(attribute string, values []string) *systems.ErrorData
}

// ShoppingListItemImageFactory used to handle create, update and delete shopping list item images
type ShoppingListItemImageFactory struct {
	DB *gorm.DB
}

// Create function used to store shopping list item images in database
func (sliif *ShoppingListItemImageFactory) Create(shoppingListItemGUID string, images []map[string]string) ([]*ShoppingListItemImage, *systems.ErrorData) {
	createdImages := make([]*ShoppingListItemImage, len(images))
	for key, image := range images {

		shoppingListItemImage := &ShoppingListItemImage{
			GUID:                 Helper.GenerateUUID(),
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
func (sliif *ShoppingListItemImageFactory) Delete(attribute string, values []string) *systems.ErrorData {
	fmt.Println(values)
	for _, value := range values {
		result := sliif.DB.Where(attribute+" = ?", value).Delete(&ShoppingListItemImage{})

		if result.Error != nil {
			return Error.InternalServerError(result.Error, systems.DatabaseError)
		}
	}

	return nil
}
