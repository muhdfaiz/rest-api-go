package v1

import (
	"mime/multipart"
	"net/url"
	"os"
	"strings"

	"bitbucket.org/cliqers/shoppermate-api/services/filesystem"
	"bitbucket.org/cliqers/shoppermate-api/systems"

	"github.com/jinzhu/gorm"
)

type ShoppingListItemImageServiceInterface interface {
	UploadImages(images []*multipart.FileHeader) ([]map[string]string, *systems.ErrorData)
	DeleteImages(ImageURLs []string) *systems.ErrorData
}

// ShoppingListItemImageService will handle all task related to upload & delete image in Amazon S3
type ShoppingListItemImageService struct {
	DB                 *gorm.DB
	AmazonS3FileSystem *filesystem.AmazonS3Upload
}

// UploadImages function used to upload multiple image for shopping list item to Amazon S3
func (sliis *ShoppingListItemImageService) UploadImages(images []*multipart.FileHeader) ([]map[string]string, *systems.ErrorData) {
	uploadedFiles := make([]map[string]string, len(images))

	for key, image := range images {
		image, err := image.Open()

		if err != nil {
			return nil, Error.InternalServerError(err.Error(), systems.CannotReadFile)
		}

		// Validate file type is image
		err1 := FileValidation.ValidateFileType([]string{"jpg", "jpeg", "png", "gif"}, image)
		if err1 != nil {
			return nil, err1
		}

		// Validate file size
		_, err1 = FileValidation.ValidateFileSize(image, 1000000, "images")
		if err1 != nil {
			return nil, err1
		}

		localUploadPath := os.Getenv("GOPATH") + Config.Get("app.yaml", "storage_path", "src/bitbucket.org/cliqers/shoppermate-api/storages/")
		amazonS3UploadPath := "/item_images/"
		uploadedFile, err1 := sliis.AmazonS3FileSystem.Upload(image, localUploadPath, amazonS3UploadPath)

		if err1 != nil {
			return nil, err1
		}

		uploadedFiles[key] = uploadedFile

	}

	return uploadedFiles, nil
}

// DeleteImages function used to delete multiple images from Amazon S3
func (sliis *ShoppingListItemImageService) DeleteImages(ImageURLs []string) *systems.ErrorData {
	imageURLs := make([]string, len(ImageURLs))

	for key, imageURL := range ImageURLs {
		url, _ := url.Parse(imageURL)

		uriSegments := strings.SplitN(url.Path, "/", 3)

		imageURLs[key] = uriSegments[2]
	}

	err := sliis.AmazonS3FileSystem.Delete(imageURLs)

	if err != nil {
		return err
	}

	return nil
}
