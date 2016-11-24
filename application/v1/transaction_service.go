package v1

// TransactionServiceInterface is a contract that defines the methods needed for Transaction Service
type TransactionServiceInterface interface {
	GetUserTransactionsForSpecificStatus(userGUID string, transactionStatus string, relations string) []*Transaction
}

// TransactionService used to encapsulates semantic gap domain layer (Transaction Handler) and persistence layer (Transaction Repository)
type TransactionService struct {
	TransactionRepository TransactionRepositoryInterface
}

// GetUserTransactionsForSpecificStatus function used to retrieve list of user transactions that match the transaction status
func (tg *TransactionService) GetUserTransactionsForSpecificStatus(userGUID string, transactionStatus string, relations string) []*Transaction {
	transactions := tg.TransactionRepository.GetByUserGUIDAndStatus(userGUID, transactionStatus, relations)

	return transactions
}
