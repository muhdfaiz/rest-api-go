package v11

import "github.com/jinzhu/gorm"

// TransactionStatusRepository will handle all CRUD task related to Transaction Status resource.
type TransactionStatusRepository struct {
	DB *gorm.DB
}

// GetBySlug function used to retrieve first transaction status found filter by slug.
func (tsr *TransactionStatusRepository) GetBySlug(slug string) *TransactionStatus {
	transactionStatus := &TransactionStatus{}

	tsr.DB.Model(&TransactionStatus{}).Where(&TransactionStatus{Slug: slug}).First(&transactionStatus)

	return transactionStatus
}

// GetByGUID function used to retrieve first transaction status found filter by GUID.
func (tsr *TransactionStatusRepository) GetByGUID(transactionStatusGUID string) *TransactionStatus {
	transactionStatus := &TransactionStatus{}

	tsr.DB.Model(&TransactionStatus{}).Where(&TransactionStatus{GUID: transactionStatusGUID}).First(&transactionStatus)

	return transactionStatus
}
