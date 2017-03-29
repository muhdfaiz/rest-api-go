package v11

// TransactionStatusServiceInterface is a contract that defines the method needed
// for Transaction Status Service.
type TransactionStatusServiceInterface interface {
	GetTransactionStatusByGUID(transactionStatusGUID string) *TransactionStatus
	GetTransactionStatusBySlug(transactionStatusSlug string) *TransactionStatus
}

// TransactionStatusRepositoryInterface is a contract that defines the method needed
// for Transaction Status Repository.
type TransactionStatusRepositoryInterface interface {
	GetBySlug(slug string) *TransactionStatus
	GetByGUID(transactionStatusGUID string) *TransactionStatus
}