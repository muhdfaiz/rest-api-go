package v1

import (
	"mime/multipart"
	"strings"

	"github.com/jinzhu/gorm"

	"os"

	"bitbucket.org/cliqers/shoppermate-api/services/filesystem"
	"bitbucket.org/cliqers/shoppermate-api/systems"
)

type DealCashbackTransactionService struct {
	AmazonS3FileSystem                *filesystem.AmazonS3Upload
	DealCashbackRepository            DealCashbackRepositoryInterface
	DealCashbackTransactionRepository DealCashbackTransactionRepositoryInterface
	DealRepository                    DealRepositoryInterface
	TransactionRepository             TransactionRepositoryInterface
}

// CreateTransaction function used to upload receipt image to Amazon S3 and create new transaction
// through DealCashbackTransactionRepository.
func (dcts *DealCashbackTransactionService) CreateTransaction(dbTransaction *gorm.DB, receipt *multipart.FileHeader, userGUID string,
	dealCashbackGUIDs string, relations string) (*Transaction, *systems.ErrorData) {

	uploadedReceipt, err := dcts.UploadReceipt(receipt)

	if err != nil {
		return nil, err
	}

	// Split on comma.
	splitDealCashbackGUID := strings.Split(dealCashbackGUIDs, ",")

	dealGUIDs := make([]string, len(splitDealCashbackGUID))

	for key, dealCashbackGUID := range splitDealCashbackGUID {
		deal := dcts.DealCashbackRepository.GetByGUID(dealCashbackGUID)

		if deal.GUID == "" {
			return nil, Error.ResourceNotFoundError("Deal Cashback GUID", "guid", dealCashbackGUID)
		}

		dealGUIDs[key] = deal.DealGUID
	}

	totalCashbackAmount := dcts.DealRepository.SumCashbackAmount(dealGUIDs)

	transactionData := &CreateTransaction{
		UserGUID:              userGUID,
		TransactionTypeGUID:   "8def67d0-fe5a-5d5e-806c-5b395a45396c",
		TransactionStatusGUID: "0f9e1582-d618-590c-bd7c-6850555ef8bb",
		Amount:                totalCashbackAmount,
		ReferenceID:           Helper.GenerateUniqueShortID(),
	}

	transaction, err := dcts.TransactionRepository.Create(dbTransaction, transactionData)

	if err != nil {
		return nil, err
	}

	result, err := dcts.DealCashbackTransactionRepository.Create(dbTransaction, userGUID, transaction.GUID, uploadedReceipt["path"])

	if err != nil {
		return nil, err
	}

	err = dcts.DealCashbackRepository.UpdateDealCashbackTransactionGUID(dbTransaction, splitDealCashbackGUID, result.GUID)

	if err != nil {
		return nil, err
	}

	transaction = dcts.TransactionRepository.GetByGUID(transaction.GUID, relations)

	return transaction, nil
}

// UploadReceipt function used to upload receipt image to Amazon S3.
// Allowed file types are jpg, jpeg, png and gif.
// Maximum file size allow is 5MB.
func (dcts *DealCashbackTransactionService) UploadReceipt(receiptImage *multipart.FileHeader) (map[string]string, *systems.ErrorData) {
	image, err := receiptImage.Open()

	if err != nil {
		return nil, Error.InternalServerError(err.Error(), systems.CannotReadFile)
	}

	err1 := FileValidation.ValidateFileType([]string{"jpg", "jpeg", "png", "gif"}, image)
	if err1 != nil {
		return nil, err1
	}

	_, err1 = FileValidation.ValidateFileSize(image, 5000000, "receipt_image")
	if err1 != nil {
		return nil, err1
	}

	localUploadPath := os.Getenv("GOPATH") + os.Getenv("STORAGE_PATH")
	amazonS3UploadPath := "/deal_cashback_receipts/"

	uploadedReceiptImage, err1 := dcts.AmazonS3FileSystem.Upload(image, localUploadPath, amazonS3UploadPath)

	if err1 != nil {
		return nil, err1
	}

	image.Close()

	return uploadedReceiptImage, nil
}
