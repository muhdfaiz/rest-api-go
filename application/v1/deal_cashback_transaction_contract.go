package v1

import (
	"mime/multipart"

	"bitbucket.org/cliqers/shoppermate-api/systems"
)

// DealCashbackTransactionServiceInterface is a contract that defines the method needed for Deal Cashback Transaction Service.
type DealCashbackTransactionServiceInterface interface {
	CreateTransaction(receipt *multipart.FileHeader, userGUID string, dealCashbackGUIDs string,
		relations string) (*Transaction, *systems.ErrorData)
	UploadReceipt(images *multipart.FileHeader) (map[string]string, *systems.ErrorData)
}

// DealCashbackTransactionRepositoryInterface is a contract that defines the method needed for Deal Cashback Transaction Repository.
type DealCashbackTransactionRepositoryInterface interface {
	Create(userGUID string, transactionGUID, receiptURL string) (*DealCashbackTransaction, *systems.ErrorData)
	CountByUserGUID(userGUID string) *int
}
