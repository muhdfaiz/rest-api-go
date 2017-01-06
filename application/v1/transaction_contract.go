package v1

import (
	"net/http"

	"github.com/jinzhu/gorm"

	"bitbucket.org/cliqers/shoppermate-api/systems"
)

// TransactionServiceInterface is a contract that defines the methods needed for Transaction Service
type TransactionServiceInterface interface {
	CreateTransaction(dbTransaction *gorm.DB, userGUID, transactionTypeGUID, transactionStatusGUID string, amount float64) (*Transaction, *systems.ErrorData)
	ViewTransactionDetails(transactionGUID, relations string) *Transaction
	ViewDealCashbackTransactionAndUpdateReadStatus(dbTransaction *gorm.DB, userGUID, transactionGUID string) (*Transaction, *systems.ErrorData)
	CheckIfUserHasPendingCashoutTransaction(userGUID string) *systems.ErrorData
	ViewCashoutTransactionAndUpdateReadStatus(dbTransaction *gorm.DB, userGUID, transactionGUID string) (*Transaction, *systems.ErrorData)
	ViewReferralCashbackTransaction(userGUID string, transactionGUID string) (*Transaction, *systems.ErrorData)
	GetUserTransactions(request *http.Request, userGUID, transactionStatus, isRead, pageNumber, pageLimit string) *TransactionResponse
	SumTotalAmountOfUserPendingTransaction(userGUID string) float64
	SumTotalAmountOfUserCashoutTransaction(userGUID string) float64
}

// TransactionRepositoryInterface is a contract that define the methods needed for Transaction Repository
type TransactionRepositoryInterface interface {
	Create(dbTransaction *gorm.DB, createTransactionData *CreateTransaction) (*Transaction, *systems.ErrorData)
	UpdateReadStatus(dbTransaction *gorm.DB, transactionGUID string, readStatus int) *systems.ErrorData
	GetByGUID(GUID, relations string) *Transaction
	GetByUserGUID(userGUID, relations string) []*Transaction
	GetByUserGUIDAndTransactionTypeGUIDAndTransactionStatusGUID(userGUID, transactionTypeGUID, transactionStatusGUID, relations string) []*Transaction
	SumTotalAmountOfUserPendingTransactions(userGUID string) float64
	SumTotalAmountOfUserCashoutTransaction(userGUID string) float64
}
