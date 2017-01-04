package v1

import (
	"mime/multipart"

	"github.com/jinzhu/gorm"

	"bitbucket.org/cliqers/shoppermate-api/systems"
)

// ShoppingListItemImageServiceInterface is a contract that defines the methods needed for Shopping List Item Image Service
type ShoppingListItemImageServiceInterface interface {
	ViewUserShoppingListItemImage(userGUID string, shoppingListGUID string, shoppingListItemGUID string,
		shoppingListItemImageGUID string, relations string) (*ShoppingListItemImage, *systems.ErrorData)
	CreateUserShoppingListItemImage(dbTransaction *gorm.DB, userGUID string, shoppingListGUID string, shoppingListItemGUID string,
		imagesToUpload []*multipart.FileHeader) ([]*ShoppingListItemImage, *systems.ErrorData)
	UploadShoppingListItemImages(imagesToUpload []multipart.File) ([]map[string]string, *systems.ErrorData)
	ValidateShoppingListItemImages(imagesToUpload []*multipart.FileHeader) ([]multipart.File, *systems.ErrorData)
	DeleteImagesForShoppingList(dbTransaction *gorm.DB, shoppingListGUID string) *systems.ErrorData
	DeleteImagesForShoppingListItem(dbTransaction *gorm.DB, shoppingListItemGUID string) *systems.ErrorData
	DeleteShoppingListItemImages(dbTransaction *gorm.DB, userGUID string, shoppingListGUID string, shoppingListItemGUID string,
		shoppingListItemImageGUIDs string) *systems.ErrorData
	CheckUserShoppingListItemImageExistOrNot(userGUID string, shoppingListGUID string,
		shoppingListItemGUID string, shoppingListItemImageGUID string) (*ShoppingListItemImage, *systems.ErrorData)
	CheckMultipleUserShoppingListItemImageExistOrNot(userGUID string, shoppingListGUID string,
		shoppingListItemGUID string, shoppingListItemImageGUIDs []string) ([]*ShoppingListItemImage, *systems.ErrorData)
}

// ShoppingListItemImageRepositoryInterface is a contract that defines the methods needed for Shopping List Item Image Repository.
type ShoppingListItemImageRepositoryInterface interface {
	Create(dbTransaction *gorm.DB, userGUID string, shoppingListGUID string, shoppingListItemGUID string,
		images []map[string]string) ([]*ShoppingListItemImage, *systems.ErrorData)
	Delete(dbTransaction *gorm.DB, attribute string, value string) *systems.ErrorData
	GetByUserGUIDAndShoppingListGUIDAndItemGUIDAndImageGUID(userGUID string, shoppingListGUID string,
		shoppingListItemGUID string, shoppingListItemImageGUID string, relations string) *ShoppingListItemImage
	GetByItemGUID(shoppingListItemGUID string, relations string) []*ShoppingListItemImage
	GetByShoppingListGUID(shoppingListGUID string, relations string) []*ShoppingListItemImage
	GetByUserGUID(useerGUID string, relations string) []*ShoppingListItemImage
}
