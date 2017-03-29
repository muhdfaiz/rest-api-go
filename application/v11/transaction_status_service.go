package v11

// TransactionStatusService used to handle application logic related to Transaction Status resource.
type TransactionStatusService struct {
	TransactionStatusRepository TransactionStatusRepositoryInterface
}

// GetTransactionStatusByGUID function used to retrieve transaction status by transaction Status GUID.
func (tss *TransactionStatusService) GetTransactionStatusByGUID(transactionStatusGUID string) *TransactionStatus {
	transactionStatus := tss.TransactionStatusRepository.GetByGUID(transactionStatusGUID)

	return transactionStatus
}

// GetTransactionStatusBySlug function used to retrieve transaction status by transaction status GUID.
func (tss *TransactionStatusService) GetTransactionStatusBySlug(transactionStatusSlug string) *TransactionStatus {
	transactionStatus := tss.TransactionStatusRepository.GetBySlug(transactionStatusSlug)

	return transactionStatus
}
