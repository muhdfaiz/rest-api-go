package v1

import (
	"net/http"

	"bitbucket.org/cliqers/shoppermate-api/systems"
)

// TransactionServiceInterface is a contract that defines the methods needed for Transaction Service
type TransactionServiceInterface interface {
	ViewDealCashbackTransactionAndUpdateReadStatus(userGUID string, transactionGUID string) (*Transaction, *systems.ErrorData)
	GetUserTransactionsForSpecificStatus(request *http.Request, userGUID string, transactionStatus string,
		isRead string, pageNumber string, pageLimit string, relations string) *TransactionResponse
}

// TransactionService used to encapsulates semantic gap domain layer (Transaction Handler) and persistence layer (Transaction Repository)
type TransactionService struct {
	TransactionRepository  TransactionRepositoryInterface
	TransactionTransformer TransactionTransformerInterface
	DealCashbackRepository DealCashbackRepositoryInterface
	ItemRepository         ItemRepositoryInterface
	GrocerRepository       GrocerRepositoryInterface
	ShoppingListRepository ShoppingListRepositoryInterface
}

// ViewTransactionAndUpdateReadStatus function used to view transaction details and update `read_status` if the transaction
// not equal to `pending`
func (ts *TransactionService) ViewDealCashbackTransactionAndUpdateReadStatus(userGUID string, transactionGUID string) (*Transaction, *systems.ErrorData) {

	relations := "transactionstatuses,transactiontypes,dealcashbacktransactions.receipt.receiptitems"

	transaction := ts.TransactionRepository.GetByGUID(transactionGUID, relations)

	if transaction.GUID == "" {
		return nil, Error.ResourceNotFoundError("Transaction", "guid", transactionGUID)
	}

	if transaction.Transactionstatuses.Slug != "pending" {
		error := ts.TransactionRepository.UpdateReadStatus(transactionGUID, 1)

		if error != nil {
			return nil, error
		}
	}

	// Get deal cashbacks by deal_cashback_transaction_guid group by shopping list id to retrieve unique shopping list
	dealsCashbacksGroupByShoppingListGUID := ts.DealCashbackRepository.GetByDealCashbackTransactionGUIDAndGroupByShoppingListGUID(&transaction.Dealcashbacktransactions.GUID)

	totalDeal := 0

	for _, dealCashbackGroupByShoppingListGUID := range dealsCashbacksGroupByShoppingListGUID {

		shoppingListWithDealCashbacks := ts.ShoppingListRepository.GetByGUIDPreloadWithDealCashbacks(dealCashbackGroupByShoppingListGUID.ShoppingListGUID,
			transaction.Dealcashbacktransactions.GUID, "")

		totalDeal = totalDeal + len(shoppingListWithDealCashbacks.Dealcashbacks)

		transaction.Dealcashbacktransactions.DealCashbackGroupByShoppingListName = append(transaction.Dealcashbacktransactions.DealCashbackGroupByShoppingListName, shoppingListWithDealCashbacks)

		for key, dealCashback := range shoppingListWithDealCashbacks.Dealcashbacks {
			shoppingListWithDealCashbacks.Dealcashbacks[key].Deals.CanAddTolist = 1

			total := ts.DealCashbackRepository.CountByDealGUIDAndUserGUID(dealCashback.Deals.GUID, userGUID)

			if total >= shoppingListWithDealCashbacks.Dealcashbacks[key].Deals.Perlimit {
				shoppingListWithDealCashbacks.Dealcashbacks[key].Deals.CanAddTolist = 0
			}

			shoppingListWithDealCashbacks.Dealcashbacks[key].Deals.NumberOfDealAddedToList = total
			shoppingListWithDealCashbacks.Dealcashbacks[key].Deals.RemainingAddToList = dealCashback.Deals.Perlimit - total
			shoppingListWithDealCashbacks.Dealcashbacks[key].Deals.Items = ts.ItemRepository.GetByID(dealCashback.Deals.ItemID, "Categories,Subcategories")
			shoppingListWithDealCashbacks.Dealcashbacks[key].Deals.Grocerexclusives = ts.GrocerRepository.GetByID(dealCashback.Deals.GrocerExclusive, "")
		}
	}

	transaction.Dealcashbacktransactions.TotalDeal = totalDeal

	return transaction, nil
}

// GetUserTransactionsForSpecificStatus function used to retrieve list of user transactions that match the transaction status
func (ts *TransactionService) GetUserTransactionsForSpecificStatus(request *http.Request, userGUID string, transactionStatus string,
	isRead string, pageNumber string, pageLimit string, relations string) *TransactionResponse {

	transactions, totalNumberOfTransaction := ts.TransactionRepository.GetByUserGUIDAndStatusAndReadStatus(userGUID, transactionStatus, isRead, pageNumber, pageLimit, relations)

	transactionsResponse := ts.TransactionTransformer.transformCollection(request, transactions, totalNumberOfTransaction, pageLimit)

	return transactionsResponse
}
