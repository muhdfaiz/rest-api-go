package v1

import (
	"bitbucket.org/cliqers/shoppermate-api/systems"
	"github.com/jinzhu/gorm"
)

type CashoutTransactionRepositoryInterface interface {
	Create(userGUID string, transactionGUID string, cashoutTransactionData *CreateCashoutTransaction) (*CashoutTransaction, *systems.ErrorData)
}

type CashoutTransactionRepository struct {
	DB *gorm.DB
}

func (ctr *CashoutTransactionRepository) Create(userGUID string, transactionGUID string, cashoutTransactionData *CreateCashoutTransaction) (*CashoutTransaction, *systems.ErrorData) {
	cashoutTransaction := &CashoutTransaction{
		GUID:                  Helper.GenerateUUID(),
		UserGUID:              userGUID,
		TransactionGUID:       transactionGUID,
		BankName:              cashoutTransactionData.BankName,
		BankAccountHolderName: cashoutTransactionData.BankAccountHolderName,
		BankAccountNumber:     cashoutTransactionData.BankAccountNumber,
		BankCountry:           cashoutTransactionData.BankCountry,
	}

	result := ctr.DB.Create(cashoutTransaction)

	if result.Error != nil || result.RowsAffected == 0 {
		return nil, Error.InternalServerError(result.Error, systems.DatabaseError)
	}

	return result.Value.(*CashoutTransaction), nil
}
