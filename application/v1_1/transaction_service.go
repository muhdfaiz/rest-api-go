package v1_1

import (
	"net/http"
	"strings"

	"github.com/jinzhu/gorm"

	"strconv"

	"bitbucket.org/cliqers/shoppermate-api/systems"
)

// TransactionService used to encapsulates semantic gap domain layer (Transaction Handler) and persistence layer (Transaction Repository)
type TransactionService struct {
	TransactionRepository    TransactionRepositoryInterface
	TransactionTransformer   TransactionTransformerInterface
	DealCashbackService      DealCashbackServiceInterface
	ShoppingListService      ShoppingListServiceInterface
	DealService              DealServiceInterface
	TransactionTypeService   TransactionTypeServiceInterface
	TransactionStatusService TransactionStatusServiceInterface
}

// CreateTransaction function used to create new user transaction and store in database.
func (ts *TransactionService) CreateTransaction(dbTransaction *gorm.DB, userGUID, transactionTypeGUID, transactionStatusGUID string, amount float64) (*Transaction, *systems.ErrorData) {
	transactionData := &CreateTransaction{
		UserGUID:              userGUID,
		TransactionTypeGUID:   transactionTypeGUID,
		TransactionStatusGUID: transactionStatusGUID,
		Amount:                amount,
		ReferenceID:           Helper.GenerateUniqueShortID(),
	}

	transaction, err := ts.TransactionRepository.Create(dbTransaction, transactionData)

	if err != nil {
		return nil, err
	}

	return transaction, nil
}

// ViewTransactionDetails function used to view transaction details.
func (ts *TransactionService) ViewTransactionDetails(transactionGUID string, relations string) *Transaction {
	transaction := ts.TransactionRepository.GetByGUID(transactionGUID, relations)

	return transaction
}

// ViewDealCashbackTransactionAndUpdateReadStatus function used to view transaction details and update `read_status` if the transaction
// not equal to `pending`.
func (ts *TransactionService) ViewDealCashbackTransactionAndUpdateReadStatus(dbTransaction *gorm.DB, userGUID string, transactionGUID string) (*Transaction, *systems.ErrorData) {

	relations := "transactionstatuses,transactiontypes,dealcashbacktransactions.receipt.receiptitems"

	transaction := ts.TransactionRepository.GetByGUID(transactionGUID, relations)

	if transaction.GUID == "" {
		return nil, Error.ResourceNotFoundError("Transaction", "guid", transactionGUID)
	}

	if transaction.Transactionstatuses.Slug != "pending" {
		error := ts.TransactionRepository.UpdateReadStatus(dbTransaction, transactionGUID, 1)

		if error != nil {
			return nil, error
		}
	}

	// Get deal cashbacks by deal_cashback_transaction_guid group by shopping list id to retrieve unique shopping list
	dealsCashbacksGroupByShoppingListGUID := ts.DealCashbackService.GetDealCashbacksByTransactionGUIDAndGroupByShoppingList(transaction.Dealcashbacktransactions.GUID)

	totalDeal := 0

	for _, dealCashbackGroupByShoppingListGUID := range dealsCashbacksGroupByShoppingListGUID {
		shoppingListWithDealCashbacks := ts.ShoppingListService.GetShoppingListIncludingDealCashbacks(dealCashbackGroupByShoppingListGUID.ShoppingListGUID,
			transaction.Dealcashbacktransactions.GUID)

		totalDeal = totalDeal + len(shoppingListWithDealCashbacks.Dealcashbacks)

		transaction.Dealcashbacktransactions.DealCashbackGroupByShoppingListName = append(transaction.Dealcashbacktransactions.DealCashbackGroupByShoppingListName, shoppingListWithDealCashbacks)

		for key := range shoppingListWithDealCashbacks.Dealcashbacks {
			deal := ts.DealService.SetAddTolistInfoAndItemsAndGrocerExclusiveForDeal(shoppingListWithDealCashbacks.Dealcashbacks[key].Deals, userGUID)

			shoppingListWithDealCashbacks.Dealcashbacks[key].Deals = deal
		}
	}

	transaction.Dealcashbacktransactions.TotalDeal = totalDeal

	return transaction, nil
}

// CheckIfUserHasPendingCashoutTransaction function used to check if user has any cashout transaction with status `pending`.
// API only allow one time cashout transaction until the cashout transaction status become `approve` or `reject`.
func (ts *TransactionService) CheckIfUserHasPendingCashoutTransaction(userGUID string) *systems.ErrorData {
	pendingTransactionStatus := ts.TransactionStatusService.GetTransactionStatusBySlug("pending")

	cashoutTransactionType := ts.TransactionTypeService.GetTransactionTypeBySlug("cashout")

	transactions := ts.TransactionRepository.GetByUserGUIDAndTransactionTypeGUIDAndTransactionStatusGUID(userGUID, cashoutTransactionType.GUID,
		pendingTransactionStatus.GUID, "")

	if len(transactions) > 0 {
		return Error.GenericError("422", systems.StillHasPendingCashoutTransaction, "Pending Cashout Transaction.", "message", "You still has cashout transaction with status pending.")
	}

	return nil
}

// ViewCashoutTransactionAndUpdateReadStatus function used to view cashout transaction details and update `read_status` if the transaction
// not equal to `pending`
func (ts *TransactionService) ViewCashoutTransactionAndUpdateReadStatus(dbTransaction *gorm.DB, userGUID string, transactionGUID string) (*Transaction, *systems.ErrorData) {

	relations := "transactionstatuses,transactiontypes,cashouttransactions"

	cashoutTransaction := ts.TransactionRepository.GetByGUID(transactionGUID, relations)

	if cashoutTransaction.GUID == "" {
		return nil, Error.ResourceNotFoundError("Transaction", "guid", transactionGUID)
	}

	if cashoutTransaction.Transactionstatuses.Slug != "pending" {
		error := ts.TransactionRepository.UpdateReadStatus(dbTransaction, transactionGUID, 1)

		if error != nil {
			return nil, error
		}
	}

	return cashoutTransaction, nil
}

// ViewReferralCashbackTransaction function used to view referral cashback transaction.
func (ts *TransactionService) ViewReferralCashbackTransaction(userGUID string, transactionGUID string) (*Transaction, *systems.ErrorData) {

	relations := "transactionstatuses,transactiontypes,referralCashbackTransactions.referrers"

	referralCashbackTransaction := ts.TransactionRepository.GetByGUID(transactionGUID, relations)

	if referralCashbackTransaction.GUID == "" {
		return nil, Error.ResourceNotFoundError("Transaction", "guid", transactionGUID)
	}

	return referralCashbackTransaction, nil
}

// GetUserTransactions function used to retrieve list of user transactions that match the transaction status and read status
// transaction status and read status is optional.
func (ts *TransactionService) GetUserTransactions(request *http.Request, userGUID string, transactionStatus string,
	isRead string, pageNumber string, pageLimit string) *TransactionResponse {

	transactions := ts.TransactionRepository.GetByUserGUID(userGUID, "transactionstatuses,transactiontypes")

	transactions = ts.FilterByTransactionStatus(transactions, transactionStatus)

	transactions = ts.FilterByReadStatus(transactions, isRead)

	totalTransactions := len(transactions)

	if pageLimit != "" && pageNumber != "" {
		transactions = ts.PaginateTransactions(transactions, pageNumber, pageLimit)
	}

	transactionsResponse := ts.TransactionTransformer.transformCollection(request, transactions, totalTransactions, pageLimit)

	return transactionsResponse
}

// FilterByTransactionStatus function used to filter transactions that match with any of transaction status slug specified.
// Able to specify multiple transaction status (pending, approved, rejected, partial_success)
func (ts *TransactionService) FilterByTransactionStatus(transactions []*Transaction, transactionStatus string) []*Transaction {
	if transactionStatus == "" {
		return transactions
	}

	transactionStatuses := strings.Split(transactionStatus, ",")

	filteredTransactions := []*Transaction{}

	for _, transaction := range transactions {
		stringExistInSlice := Helper.StringInSlice(transaction.Transactionstatuses.Slug, transactionStatuses)

		if stringExistInSlice == true {
			filteredTransactions = append(filteredTransactions, transaction)
		}
	}

	return filteredTransactions
}

// FilterByReadStatus function used to filter transactions that match with read status value specified.
func (ts *TransactionService) FilterByReadStatus(transactions []*Transaction, readStatus string) []*Transaction {
	if readStatus == "" {
		return transactions
	}

	filteredTransactions := []*Transaction{}

	readStatusInInt, _ := strconv.Atoi(readStatus)

	for _, transaction := range transactions {
		if transaction.ReadStatus == readStatusInInt {
			filteredTransactions = append(filteredTransactions, transaction)
		}
	}

	return filteredTransactions
}

// SumTotalAmountOfUserPendingTransaction function used to sum total amount of user transaction with status pending.
func (ts *TransactionService) SumTotalAmountOfUserPendingTransaction(userGUID string) float64 {

	totalAmountOfPendingDealCashbackTransactions := ts.TransactionRepository.SumTotalAmountOfUserPendingTransactions(userGUID)

	return totalAmountOfPendingDealCashbackTransactions
}

// SumTotalAmountOfUserCashoutTransaction function used to sum total amount of user transaction with status pending.
func (ts *TransactionService) SumTotalAmountOfUserCashoutTransaction(userGUID string) float64 {

	totalAmountOfCashoutTransaction := ts.TransactionRepository.SumTotalAmountOfUserCashoutTransaction(userGUID)

	return totalAmountOfCashoutTransaction
}

func (ts *TransactionService) PaginateTransactions(transactions []*Transaction, pageNumber string, pageLimit string) []*Transaction {
	paginatedTransactions := []*Transaction{}

	pageNumberInInt, _ := strconv.Atoi(pageNumber)
	pageLimitInInt, _ := strconv.Atoi(pageLimit)

	offset := 0

	if pageNumberInInt != 0 && pageNumberInInt != 1 {
		offset = (pageNumberInInt * pageLimitInInt) - 1
	}

	limit := offset + pageLimitInInt

	for key, transaction := range transactions {
		if key >= offset && key < limit {
			paginatedTransactions = append(paginatedTransactions, transaction)
		}
	}

	return paginatedTransactions
}
