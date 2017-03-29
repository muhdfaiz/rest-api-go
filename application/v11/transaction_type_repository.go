package v11

import "github.com/jinzhu/gorm"

// TransactionTypeRepository will handle all CRUD function for Transaction Type.
type TransactionTypeRepository struct {
	DB *gorm.DB
}

// GetBySlug function used to retrieve first transaction type filter by slug.
func (ttr *TransactionTypeRepository) GetBySlug(transactionTypeSlug string) *TransactionType {
	transactionType := &TransactionType{}

	ttr.DB.Model(&TransactionType{}).Where(&TransactionType{Slug: transactionTypeSlug}).First(&transactionType)

	return transactionType
}

// GetByGUID function used to retrieve first transaction type filter by GUID.
func (ttr *TransactionTypeRepository) GetByGUID(transactionTypeGUID string) *TransactionType {
	transactionType := &TransactionType{}

	ttr.DB.Model(&TransactionType{}).Where(&TransactionType{GUID: transactionTypeGUID}).First(&transactionType)

	return transactionType
}
