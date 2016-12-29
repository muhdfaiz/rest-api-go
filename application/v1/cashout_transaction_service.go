package v1

import "bitbucket.org/cliqers/shoppermate-api/systems"

type CashoutTransactionService struct {
	CashoutTransactionRepository CashoutTransactionRepositoryInterface
	TransactionService           TransactionServiceInterface
	UserRepository               UserRepositoryInterface
}

// CreateCashoutTransaction function used to create cashout transaction through CashoutTransactionRepository.
func (cts *CashoutTransactionService) CreateCashoutTransaction(userGUID string, cashoutTransactionData *CreateCashoutTransaction) (*Transaction, *systems.ErrorData) {
	user := cts.UserRepository.GetByGUID(userGUID, "")
	availableCashoutAmount := user.Wallet

	if cashoutTransactionData.Amount > availableCashoutAmount {
		return nil, Error.GenericError("422", systems.CashoutAmountExceededLimit, "Cashout Amount Exceeded Limit.", "amount", "Cashout amount more than current amount available.")
	}

	transaction, error := cts.TransactionService.CreateTransaction(userGUID, "c96358c0-13ae-59ad-863f-f113ddb33c68", "0f9e1582-d618-590c-bd7c-6850555ef8bb", cashoutTransactionData.Amount)

	if error != nil {
		return nil, error
	}

	_, error = cts.CashoutTransactionRepository.Create(userGUID, transaction.GUID, cashoutTransactionData)

	if error != nil {
		return nil, error
	}

	// error = cts.UserRepository.UpdateUserWallet(userGUID, availableCashoutAmount-cashoutTransactionData.Amount)

	// if error != nil {
	// 	return nil, error
	// }

	relations := "transactiontypes,transactionstatuses,cashouttransactions,users"

	transaction = cts.TransactionService.ViewTransactionDetails(transaction.GUID, relations)

	return transaction, nil
}
