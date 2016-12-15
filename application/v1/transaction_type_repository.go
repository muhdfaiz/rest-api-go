package v1

import "github.com/jinzhu/gorm"

// TransactionTypeRepositoryInterface ia a contradt that defines the method needed
// for Transaction Type Repository.
type TransactionTypeRepositoryInterface interface {
	GetBySlug(transactionTypeSlug string) *TransactionType
	GetByGUID(transactionTypeGUID string) *TransactionType
}

// TransactionTypeRepository will handle all CRUD function for Transaction Type.
type TransactionTypeRepository struct {
	DB *gorm.DB
}

// GetBySlug function used to retrieve transaction type by slug.
func (ttr *TransactionTypeRepository) GetBySlug(transactionTypeSlug string) *TransactionType {
	transactionType := &TransactionType{}

	ttr.DB.Model(&TransactionType{}).Where(&TransactionType{Slug: transactionTypeSlug}).First(&transactionType)

	return transactionType
}

// GetByGUID function used to retrieve transaction type by transaction type GUID.
func (ttr *TransactionTypeRepository) GetByGUID(transactionTypeGUID string) *TransactionType {
	transactionType := &TransactionType{}

	ttr.DB.Model(&TransactionType{}).Where(&TransactionType{GUID: transactionTypeGUID}).First(&transactionType)

	return transactionType
}
