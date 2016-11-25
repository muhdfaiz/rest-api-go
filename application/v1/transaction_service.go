package v1

import "net/http"

// TransactionServiceInterface is a contract that defines the methods needed for Transaction Service
type TransactionServiceInterface interface {
	GetUserTransactionsForSpecificStatus(request *http.Request, userGUID string, transactionStatus string,
		isRead string, pageNumber string, pageLimit string, relations string) *TransactionResponse
}

// TransactionService used to encapsulates semantic gap domain layer (Transaction Handler) and persistence layer (Transaction Repository)
type TransactionService struct {
	TransactionRepository  TransactionRepositoryInterface
	TransactionTransformer TransactionTransformerInterface
}

// GetUserTransactionsForSpecificStatus function used to retrieve list of user transactions that match the transaction status
func (tg *TransactionService) GetUserTransactionsForSpecificStatus(request *http.Request, userGUID string, transactionStatus string,
	isRead string, pageNumber string, pageLimit string, relations string) *TransactionResponse {

	transactions, totalNumberOfTransaction := tg.TransactionRepository.GetByUserGUIDAndStatusAndReadStatus(userGUID, transactionStatus, isRead, pageNumber, pageLimit, relations)

	transactionsResponse := tg.TransactionTransformer.transformCollection(request, transactions, totalNumberOfTransaction, pageLimit)

	return transactionsResponse
}
