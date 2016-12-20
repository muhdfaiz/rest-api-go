package v1

import (
	"net/http"

	"bitbucket.org/cliqers/shoppermate-api/systems"
)

// TransactionServiceInterface is a contract that defines the methods needed for Transaction Service
type TransactionServiceInterface interface {
	CreateTransaction(userGUID string, transactionTypeGUID string, amount float64) (*Transaction, *systems.ErrorData)
	ViewTransactionDetails(transactionGUID string, relations string) *Transaction
	ViewDealCashbackTransactionAndUpdateReadStatus(userGUID string, transactionGUID string) (*Transaction, *systems.ErrorData)
	ViewCashoutTransactionAndUpdateReadStatus(userGUID string, transactionGUID string) (*Transaction, *systems.ErrorData)
	GetUserTransactions(request *http.Request, userGUID string, transactionStatus string,
		isRead string, pageNumber string, pageLimit string) *TransactionResponse
	SumTotalAmountOfUserPendingTransaction(userGUID string) float64
	SumTotalAmountOfUserCashoutTransaction(userGUID string) float64
}

// TransactionRepositoryInterface is a contract that define the methods needed for Transaction Repository
type TransactionRepositoryInterface interface {
	Create(createTransactionData *CreateTransaction) (*Transaction, *systems.ErrorData)
	UpdateReadStatus(transactionGUID string, readStatus int) *systems.ErrorData
	GetByGUID(GUID string, relations string) *Transaction
	GetByUserGUID(userGUID string, relations string) []*Transaction
	SumTotalAmountOfUserPendingTransactions(userGUID string) float64
	SumTotalAmountOfUserCashoutTransaction(userGUID string) float64
}
