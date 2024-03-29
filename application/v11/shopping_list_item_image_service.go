package v11

import (
	"mime/multipart"
	"net/url"
	"strings"

	"bitbucket.org/cliqers/shoppermate-api/services/filesystem"
	"bitbucket.org/cliqers/shoppermate-api/systems"

	"os"

	"github.com/jinzhu/gorm"
)

// ShoppingListItemImageService used to handle application logic related to Shopping List Item Image resource.
type ShoppingListItemImageService struct {
	DB                              *gorm.DB
	ShoppingListItemService         ShoppingListItemServiceInterface
	AmazonS3FileSystem              *filesystem.AmazonS3Upload
	ShoppingListItemImageRepository ShoppingListItemImageRepositoryInterface
}

// ViewUserShoppingListItemImage function used to retrieve user shopping list item image detail including relations
// like shopping list item, shopping list and others.
func (sliis *ShoppingListItemImageService) ViewUserShoppingListItemImage(userGUID string, shoppingListGUID string, shoppingListItemGUID string,
	shoppingListItemImageGUID string, relations string) (*ShoppingListItemImage, *systems.ErrorData) {

	_, error := sliis.CheckUserShoppingListItemImageExistOrNot(userGUID, shoppingListGUID, shoppingListItemGUID, shoppingListItemImageGUID)

	if error != nil {
		return nil, error
	}

	userShoppingListItemImage := sliis.ShoppingListItemImageRepository.GetByUserGUIDAndShoppingListGUIDAndItemGUIDAndImageGUID(userGUID, shoppingListGUID,
		shoppingListItemGUID, shoppingListItemImageGUID, relations)

	if userShoppingListItemImage.GUID == "" {
		return nil, Error.ResourceNotFoundError("Shopping List Item Image", "guid", shoppingListItemImageGUID)
	}

	return userShoppingListItemImage, nil
}

// CreateUserShoppingListItemImage function used to create multiple user shopping list item image, upload the
// images into Amazon S3 and store the URL in database.
func (sliis *ShoppingListItemImageService) CreateUserShoppingListItemImage(dbTransaction *gorm.DB, userGUID string, shoppingListGUID string, shoppingListItemGUID string,
	imagesToUpload []*multipart.FileHeader) ([]*ShoppingListItemImage, *systems.ErrorData) {

	_, error := sliis.ShoppingListItemService.CheckUserShoppingListItemExistOrNot(shoppingListItemGUID, userGUID, shoppingListGUID)

	if error != nil {
		return nil, error
	}

	if len(imagesToUpload) < 1 {
		return nil, Error.FileRequireErrors("images")
	}

	imagesReadyToUpload, error := sliis.ValidateShoppingListItemImages(imagesToUpload)

	if error != nil {
		return nil, error
	}

	uploadedImages, error := sliis.UploadShoppingListItemImages(imagesReadyToUpload)

	if error != nil {
		return nil, error
	}

	createdImages, error := sliis.ShoppingListItemImageRepository.Create(dbTransaction, userGUID, shoppingListGUID, shoppingListItemGUID, uploadedImages)

	if error != nil {
		return nil, error
	}

	return createdImages, nil
}

// UploadShoppingListItemImages function used to upload multiple image for shopping list item to Amazon S3.
func (sliis *ShoppingListItemImageService) UploadShoppingListItemImages(imagesToUpload []multipart.File) ([]map[string]string, *systems.ErrorData) {

	uploadedFiles := make([]map[string]string, len(imagesToUpload))

	for key, image := range imagesToUpload {
		localUploadPath := os.Getenv("GOPATH") + os.Getenv("STORAGE_PATH")

		amazonS3UploadPath := "/item_images/"

		uploadedFile, error := sliis.AmazonS3FileSystem.Upload(image, localUploadPath, amazonS3UploadPath)

		if error != nil {
			return nil, error
		}

		uploadedFiles[key] = uploadedFile

		image.Close()
	}

	return uploadedFiles, nil
}

// ValidateShoppingListItemImages function used to validate file type and file size multiple shopping list item image
// before upload to Amazon S3. Maximum file size allow is 1MB.
func (sliis *ShoppingListItemImageService) ValidateShoppingListItemImages(imagesToUpload []*multipart.FileHeader) ([]multipart.File, *systems.ErrorData) {
	imagesReadyToUpload := make([]multipart.File, len(imagesToUpload))

	for key, image := range imagesToUpload {

		openedImage, error := image.Open()

		if error != nil {
			return nil, Error.InternalServerError(error.Error(), systems.CannotReadFile)
		}

		error1 := FileValidation.ValidateFileType([]string{"jpg", "jpeg", "png", "gif"}, openedImage)

		if error1 != nil {
			return nil, error1
		}

		_, error1 = FileValidation.ValidateFileSize(openedImage, 1000000, "images")

		if error1 != nil {
			return nil, error1
		}

		imagesReadyToUpload[key] = openedImage
	}

	return imagesReadyToUpload, nil
}

// DeleteImagesForShoppingList function used to soft delete all images for all of user shopping list items inside
// shopping list.
func (sliis *ShoppingListItemImageService) DeleteImagesForShoppingList(dbTransaction *gorm.DB, shoppingListGUID string) *systems.ErrorData {

	shoppingListItemImages := sliis.ShoppingListItemImageRepository.GetByShoppingListGUID(shoppingListGUID, "")

	deletedImagesURI := make([]string, len(shoppingListItemImages))

	if len(shoppingListItemImages) > 0 {
		for key, shoppingListItemImage := range shoppingListItemImages {

			error := sliis.ShoppingListItemImageRepository.Delete(dbTransaction, "guid", shoppingListItemImage.GUID)

			if error != nil {
				return error
			}

			url, _ := url.Parse(shoppingListItemImage.URL)

			uriSegments := strings.SplitN(url.Path, "/", 3)

			deletedImagesURI[key] = uriSegments[2]
		}

		error := sliis.AmazonS3FileSystem.Delete(deletedImagesURI)

		if error != nil {
			return error
		}
	}

	return nil
}

// DeleteImagesForShoppingListItem function used to soft delete all images for shopping list item
func (sliis *ShoppingListItemImageService) DeleteImagesForShoppingListItem(dbTransaction *gorm.DB, shoppingListItemGUID string) *systems.ErrorData {

	shoppingListItemImages := sliis.ShoppingListItemImageRepository.GetByItemGUID(shoppingListItemGUID, "")

	deletedImagesURI := make([]string, len(shoppingListItemImages))

	if len(shoppingListItemImages) > 0 {

		for key, shoppingListItemImage := range shoppingListItemImages {

			error := sliis.ShoppingListItemImageRepository.Delete(dbTransaction, "guid", shoppingListItemImage.GUID)

			if error != nil {
				return error
			}

			url, _ := url.Parse(shoppingListItemImage.URL)

			uriSegments := strings.SplitN(url.Path, "/", 3)

			deletedImagesURI[key] = uriSegments[2]
		}

		error := sliis.AmazonS3FileSystem.Delete(deletedImagesURI)

		if error != nil {
			return error
		}
	}

	return nil
}

// DeleteShoppingListItemImages function used to soft delete multiple shopping list item images in database
// and delete the images from Amazon S3.
func (sliis *ShoppingListItemImageService) DeleteShoppingListItemImages(dbTransaction *gorm.DB, userGUID string, shoppingListGUID string,
	shoppingListItemGUID string, shoppingListItemImageGUIDs string) *systems.ErrorData {

	splitShoppingListItemImageGUID := strings.Split(shoppingListItemImageGUIDs, ",")

	imagesToDelete, error := sliis.CheckMultipleUserShoppingListItemImageExistOrNot(userGUID, shoppingListGUID, shoppingListItemGUID, splitShoppingListItemImageGUID)

	if error != nil {
		return error
	}

	deletedImagesURI := make([]string, len(imagesToDelete))

	for key, imageToDelete := range imagesToDelete {

		error := sliis.ShoppingListItemImageRepository.Delete(dbTransaction, "guid", imageToDelete.GUID)

		if error != nil {
			return error
		}

		url, _ := url.Parse(imageToDelete.URL)

		uriSegments := strings.SplitN(url.Path, "/", 3)

		deletedImagesURI[key] = uriSegments[2]
	}

	error = sliis.AmazonS3FileSystem.Delete(deletedImagesURI)

	if error != nil {
		return error
	}

	return nil
}

// CheckUserShoppingListItemImageExistOrNot function used to check user shopping list item image exist or not in database
// using user GUID, shopping list GUID, shopping list item GUID and shopping list item image GUID.
func (sliis *ShoppingListItemImageService) CheckUserShoppingListItemImageExistOrNot(userGUID string, shoppingListGUID string,
	shoppingListItemGUID string, shoppingListItemImageGUID string) (*ShoppingListItemImage, *systems.ErrorData) {

	shoppingListItemImage := sliis.ShoppingListItemImageRepository.GetByUserGUIDAndShoppingListGUIDAndItemGUIDAndImageGUID(userGUID, shoppingListGUID,
		shoppingListItemGUID, shoppingListItemImageGUID, "")

	if shoppingListItemImage.GUID == "" {
		return nil, Error.ResourceNotFoundError("Shopping List Item Image", "guid", shoppingListItemImageGUID)
	}

	return shoppingListItemImage, nil
}

// CheckMultipleUserShoppingListItemImageExistOrNot function used to check multiple user shopping list item image exist or not in database
// using user GUID, shopping list GUID, shopping list item GUID and shopping list item image GUID.
func (sliis *ShoppingListItemImageService) CheckMultipleUserShoppingListItemImageExistOrNot(userGUID string, shoppingListGUID string,
	shoppingListItemGUID string, shoppingListItemImageGUIDs []string) ([]*ShoppingListItemImage, *systems.ErrorData) {

	shoppingListItemImages := make([]*ShoppingListItemImage, len(shoppingListItemImageGUIDs))

	for key, shoppingListItemImageGUID := range shoppingListItemImageGUIDs {

		shoppingListItemImage := sliis.ShoppingListItemImageRepository.GetByUserGUIDAndShoppingListGUIDAndItemGUIDAndImageGUID(userGUID, shoppingListGUID,
			shoppingListItemGUID, shoppingListItemImageGUID, "")

		if shoppingListItemImage.GUID == "" {
			return nil, Error.ResourceNotFoundError("Shopping List Item Image", "guid", shoppingListItemImageGUID)
		}

		shoppingListItemImages[key] = shoppingListItemImage
	}

	return shoppingListItemImages, nil
}
