package v1

// TransactionTypeRepositoryInterface ia a contradt that defines the method needed
// for Transaction Type Repository.
type TransactionTypeRepositoryInterface interface {
	GetBySlug(transactionTypeSlug string) *TransactionType
	GetByGUID(transactionTypeGUID string) *TransactionType
}

// TransactionTypeServiceInterface is a contract that defines the method needed
// for Transaction Type Service.
type TransactionTypeServiceInterface interface {
	GetTransactionTypeByGUID(transactionTypeGUID string) *TransactionType
	GetTransactionTypeBySlug(transactionTypeSlug string) *TransactionType
}
