package v1

import (
	"mime/multipart"
	"os"

	"bitbucket.org/cliqers/shoppermate-api/services/filesystem"
	"bitbucket.org/cliqers/shoppermate-api/systems"
)

type DealCashbackTransactionServiceInterface interface {
	UploadReceipt(images *multipart.FileHeader) (map[string]string, *systems.ErrorData)
}

type DealCashbackTransactionService struct {
	AmazonS3FileSystem *filesystem.AmazonS3Upload
}

func (dcts *DealCashbackTransactionService) UploadReceipt(receiptImage *multipart.FileHeader) (map[string]string, *systems.ErrorData) {
	image, err := receiptImage.Open()

	if err != nil {
		return nil, Error.InternalServerError(err.Error(), systems.CannotReadFile)
	}

	// Validate file type is image
	err1 := FileValidation.ValidateFileType([]string{"jpg", "jpeg", "png", "gif"}, image)
	if err1 != nil {
		return nil, err1
	}

	// Validate file size
	_, err1 = FileValidation.ValidateFileSize(image, 5000000, "receipt_image")
	if err1 != nil {
		return nil, err1
	}

	localUploadPath := os.Getenv("GOPATH") + Config.Get("app.yaml", "storage_path", "src/bitbucket.org/cliqers/shoppermate-api/storages/")
	amazonS3UploadPath := "/deal_cashback_receipts/"

	uploadedReceiptImage, err1 := dcts.AmazonS3FileSystem.Upload(image, localUploadPath, amazonS3UploadPath)

	if err1 != nil {
		return nil, err1
	}

	return uploadedReceiptImage, nil
}
