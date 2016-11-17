package v1

import "github.com/jinzhu/gorm"

type TransactionTypeRepositoryInterface interface {
	GetBySlug(slug string) *TransactionType
}

type TransactionTypeRepository struct {
	DB *gorm.DB
}

func (ttr *TransactionTypeRepository) GetBySlug(slug string) *TransactionType {
	transactionType := &TransactionType{}

	ttr.DB.Model(&TransactionType{}).Where(&TransactionType{Slug: slug}).First(&transactionType)

	return transactionType
}
