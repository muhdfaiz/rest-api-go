package v1

import "bitbucket.org/cliqers/shoppermate-api/systems"

// CashoutTransactionServiceInterface is a contract that defines the method needed for Cashout Transaction Service.
type CashoutTransactionServiceInterface interface {
	CreateCashoutTransaction(userGUID string, cashoutTransactionData *CreateCashoutTransaction) (*Transaction, *systems.ErrorData)
}

// CashoutTransactionRepositoryInterface is a contract that defines the method
// needed for Cashout Transaction Repository.
type CashoutTransactionRepositoryInterface interface {
	Create(userGUID, transactionGUID string, cashoutTransactionData *CreateCashoutTransaction) (*CashoutTransaction, *systems.ErrorData)
}
