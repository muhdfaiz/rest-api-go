package v1

import "bitbucket.org/cliqers/shoppermate-api/systems"

type CashoutTransactionService struct {
	CashoutTransactionRepository CashoutTransactionRepositoryInterface
	TransactionService           TransactionServiceInterface
	UserRepository               UserRepositoryInterface
	TransactionTypeRepository    TransactionTypeRepositoryInterface
}

// CreateCashoutTransaction function used to create cashout transaction through CashoutTransactionRepository.
func (cts *CashoutTransactionService) CreateCashoutTransaction(userGUID string, cashoutTransactionData *CreateCashoutTransaction) (*Transaction, *systems.ErrorData) {
	user := cts.UserRepository.GetByGUID(userGUID, "")
	availableCashoutAmount := user.Wallet

	if cashoutTransactionData.Amount > availableCashoutAmount {
		return nil, Error.GenericError("422", systems.CashoutAmountExceededLimit, "Cashout Amount Exceeded Limit.", "amount", "Cashout amount more than current amount available.")
	}

	transactionTypeGUID := cts.TransactionTypeRepository.GetBySlug("cashout").GUID

	transaction, error := cts.TransactionService.CreateTransaction(userGUID, transactionTypeGUID, cashoutTransactionData.Amount)

	if error != nil {
		return nil, error
	}

	_, error = cts.CashoutTransactionRepository.Create(userGUID, transaction.GUID, cashoutTransactionData)

	if error != nil {
		return nil, error
	}

	error = cts.UserRepository.UpdateUserWallet(userGUID, availableCashoutAmount-cashoutTransactionData.Amount)

	if error != nil {
		return nil, error
	}

	relations := "transactiontypes,transactionstatuses,cashouttransactions,users"

	transaction = cts.TransactionService.ViewTransactionDetails(transaction.GUID, relations)

	return transaction, nil
}
