package v1_1

import "github.com/jinzhu/gorm"

type TransactionStatusRepository struct {
	DB *gorm.DB
}

// GetBySlug function used to retrieve transaction status by transaction status slug.
func (tsr *TransactionStatusRepository) GetBySlug(slug string) *TransactionStatus {
	transactionStatus := &TransactionStatus{}

	tsr.DB.Model(&TransactionStatus{}).Where(&TransactionStatus{Slug: slug}).First(&transactionStatus)

	return transactionStatus
}

// GetByGUID function used to retrieve transaction status by transaction GUID.
func (tsr *TransactionStatusRepository) GetByGUID(transactionStatusGUID string) *TransactionStatus {
	transactionStatus := &TransactionStatus{}

	tsr.DB.Model(&TransactionStatus{}).Where(&TransactionStatus{GUID: transactionStatusGUID}).First(&transactionStatus)

	return transactionStatus
}
