package v1

import "github.com/jinzhu/gorm"

type TransactionStatusRepositoryInterface interface {
	GetBySlug(slug string) *TransactionStatus
}

type TransactionStatusRepository struct {
	DB *gorm.DB
}

func (tsr *TransactionStatusRepository) GetBySlug(slug string) *TransactionStatus {
	transactionStatus := &TransactionStatus{}

	tsr.DB.Model(&TransactionStatus{}).Where(&TransactionStatus{Slug: slug}).First(&transactionStatus)

	return transactionStatus
}
