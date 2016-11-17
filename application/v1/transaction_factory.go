package v1

import (
	"bitbucket.org/cliqers/shoppermate-api/systems"
	"github.com/jinzhu/gorm"
)

type TransactionFactoryInterface interface {
	Create(createTransactionData *CreateTransaction) (*Transaction, *systems.ErrorData)
}

type TransactionFactory struct {
	DB                          *gorm.DB
	TransactionStatusRepository TransactionStatusRepositoryInterface
}

func (tf *TransactionFactory) Create(createTransactionData *CreateTransaction) (*Transaction, *systems.ErrorData) {
	pendingTransactionStatusGUID := tf.TransactionStatusRepository.GetBySlug("pending").GUID

	createdTransaction := &Transaction{}

	transaction := &Transaction{
		GUID:                  Helper.GenerateUUID(),
		UserGUID:              createTransactionData.UserGUID,
		TransactionTypeGUID:   createTransactionData.TransactionTypeGUID,
		TransactionStatusGUID: pendingTransactionStatusGUID,
		Amount:                createTransactionData.Amount,
	}

	result := tf.DB.Create(transaction)

	if result.Error != nil || result.RowsAffected == 0 {
		return nil, Error.InternalServerError(result.Error, systems.DatabaseError)
	}

	createdTransaction = result.Value.(*Transaction)

	return createdTransaction, nil
}
