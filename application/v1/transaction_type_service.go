package v1

// TransactionTypeServiceInterface is a contract that defines the method needed
// for Transaction Type Service.
type TransactionTypeServiceInterface interface {
	GetTransactionTypeByGUID(transactionTypeGUID string) *TransactionType
	GetTransactionTypeBySlug(transactionTypeSlug string) *TransactionType
}

type TransactionTypeService struct {
	TransactionTypeRepository TransactionTypeRepositoryInterface
}

// GetTransactionTypeByGUID function used to retrieve transaction type by transaction type GUID.
func (tts *TransactionTypeService) GetTransactionTypeByGUID(transactionTypeGUID string) *TransactionType {
	transactionType := tts.TransactionTypeRepository.GetByGUID(transactionTypeGUID)

	return transactionType
}

// GetTransactionTypeBySlug function used to retrieve transaction type by slug.
func (tts *TransactionTypeService) GetTransactionTypeBySlug(transactionTypeSlug string) *TransactionType {
	transactionType := tts.TransactionTypeRepository.GetBySlug(transactionTypeSlug)

	return transactionType
}
