package v11

import (
	"strconv"

	"bitbucket.org/cliqers/shoppermate-api/services/email"
	"bitbucket.org/cliqers/shoppermate-api/systems"
	"github.com/jinzhu/gorm"
)

// CashoutTransactionService used to handle application logic related to Cashout Transaction resource.
type CashoutTransactionService struct {
	CashoutTransactionRepository CashoutTransactionRepositoryInterface
	TransactionService           TransactionServiceInterface
	TransactionStatusService     TransactionStatusServiceInterface
	TransactionTypeService       TransactionTypeServiceInterface
	UserRepository               UserRepositoryInterface
	EmailService                 email.EmailServiceInterface
}

// CreateCashoutTransaction function used to create new cashout transaction for specific user using user GUID.
// First, it will use user repository to retrieve user info using user GUID and get the wallet amount. API used wallet amount to know amount of user can cashout.
// Then, it will check if amount of user want to checkout more than current amount available to cashout.
// It will return an error if amount of user want checkout more than current amount available.
// Then, it will retrieve pending transaction status and cashout transaction type because it needs
// pending transaction status guid cashout transaction type to create transaction.
// Then it will create transaction and cashout transaction. It needs to create transaction first because cashout transaction require transaction GUID.
// After successfully created, it will commit database transaction.
// Lastly, it will send email about cashout transaction if the user first time trying to cashout.
func (cts *CashoutTransactionService) CreateCashoutTransaction(dbTransaction *gorm.DB, userGUID string, cashoutTransactionData *CreateCashoutTransaction) (*Transaction, *systems.ErrorData) {
	user := cts.UserRepository.GetByGUID(userGUID, "")

	availableCashoutAmount := user.Wallet

	if cashoutTransactionData.Amount > availableCashoutAmount {
		dbTransaction.Rollback()
		return nil, Error.GenericError("422", systems.CashoutAmountExceededLimit, "Cashout Amount Exceeded Limit.", "amount", "Cashout amount more than current amount available.")
	}

	pendingTransactionStatus := cts.TransactionStatusService.GetTransactionStatusBySlug("pending")

	cashoutTransactionType := cts.TransactionTypeService.GetTransactionTypeBySlug("cashout")

	transaction, error := cts.TransactionService.CreateTransaction(dbTransaction, userGUID, cashoutTransactionType.GUID, pendingTransactionStatus.GUID, cashoutTransactionData.Amount)

	if error != nil {
		dbTransaction.Rollback()
		return nil, error
	}

	_, error = cts.CashoutTransactionRepository.Create(dbTransaction, userGUID, transaction.GUID, cashoutTransactionData)

	if error != nil {
		dbTransaction.Rollback()
		return nil, error
	}

	dbTransaction.Commit()

	totalNumberOfPreviousCashout := cts.CashoutTransactionRepository.CountByUserGUID(userGUID)

	if totalNumberOfPreviousCashout < 1 {
		error = cts.EmailService.SendTemplate(map[string]string{
			"name":     user.Name,
			"email":    user.Email,
			"template": "11-shoppermate-summary-of-first-submission",
			"variables": `[{"name":"user_fullname","content":"` + user.Name + `"},
			{"name":"bank_name","content":"` + cashoutTransactionData.BankName + `"},
			{"name":"bank_acc_number","content":"` + cashoutTransactionData.BankAccountNumber + `"},
			{"name":"bank_acc_name","content":"` + cashoutTransactionData.BankAccountHolderName + `"},
			{"name":"amount","content":"` + strconv.FormatFloat(cashoutTransactionData.Amount, 'f', 2, 64) + `"},
			{"name":"reference_number","content":"` + transaction.ReferenceID + `"},
			{"name":"transaction_number","content":"` + strconv.Itoa(transaction.ID) + `"}]`,
		})

		if error != nil {
			return nil, error
		}
	}

	return transaction, nil
}
