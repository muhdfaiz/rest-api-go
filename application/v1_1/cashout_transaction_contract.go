package v1_1

import (
	"bitbucket.org/cliqers/shoppermate-api/systems"
	"github.com/jinzhu/gorm"
)

// CashoutTransactionServiceInterface is a contract that defines the method needed for Cashout Transaction Service.
type CashoutTransactionServiceInterface interface {
	CreateCashoutTransaction(dbTransaction *gorm.DB, userGUID string, cashoutTransactionData *CreateCashoutTransaction) (*Transaction, *systems.ErrorData)
}

// CashoutTransactionRepositoryInterface is a contract that defines the method
// needed for Cashout Transaction Repository.
type CashoutTransactionRepositoryInterface interface {
	Create(dbTransaction *gorm.DB, userGUID, transactionGUID string, cashoutTransactionData *CreateCashoutTransaction) (*CashoutTransaction, *systems.ErrorData)
}
