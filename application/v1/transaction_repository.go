package v1

import "github.com/jinzhu/gorm"

type TransactionRepositoryInterface interface {
	GetByGUID(GUID string, relations string) *Transaction
}

type TransactionRepository struct {
	DB *gorm.DB
}

func (tr *TransactionRepository) GetByGUID(GUID string, relations string) *Transaction {
	transaction := &Transaction{}

	DB := tr.DB.Model(&Transaction{})

	if relations != "" {
		DB = LoadRelations(DB, relations)
	}

	DB.Where(&Transaction{GUID: GUID}).First(&transaction)

	return transaction
}
