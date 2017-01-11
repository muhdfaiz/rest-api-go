package v1_1

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
