package v1

import (
	"net/http"

	"bitbucket.org/cliqers/shoppermate-api/systems"
)

// DealCashbackTransformerInterface is a contract that defines the method needed for deal cashback transformer.
type DealCashbackTransformerInterface interface {
	transformCollection(currentURI *http.Request, data interface{}, totalData int, limit string) *DealCashbackResponse
}

// DealCashbackServiceInterface is a contract that defines the method needed for Deal Cashback Service.
type DealCashbackServiceInterface interface {
	CountTotalNumberOfDealUserAddedToList(userGUID string, dealGUID string) int
	SumTotalAmountOfDealAddedTolistByUser(userGUID string) float64
	GetDealCashbacksByTransactionGUIDAndGroupByShoppingList(dealCashbackTrasnactionGUID string) []*DealCashback
	CreateDealCashbackAndShoppingListItem(userGUID string, dealCashbackData CreateDealCashback) *systems.ErrorData
	GetUserDealCashbackForUserShoppingList(userGUID string, shoppingListGUID string, transactionStatus string,
		pageNumber string, pageLimit string, relations string) ([]*DealCashback, int)
}

// DealCashbackRepositoryInterface is a contract that defines the method needed for deal cashback repository.
type DealCashbackRepositoryInterface interface {
	Create(userGUID string, data CreateDealCashback) (*DealCashback, *systems.ErrorData)
	UpdateDealCashbackTransactionGUID(dealCashbackGUIDs []string, dealCashbackTransactionGUID string) *systems.ErrorData
	DeleteByUserGUIDAndDealGUID(userGUID string, dealGUID string) *systems.ErrorData
	DeleteByUserGUIDAndShoppingListGUIDAndDealGUID(userGUID string, shoppingListGUID string, dealGUID string) *systems.ErrorData
	GetByGUID(GUID string) *DealCashback
	GetByDealGUIDAndUserGUID(dealGUID, userGUID string) *DealCashback
	GetByDealCashbackTransactionGUIDAndGroupByShoppingListGUID(dealCashbackTransactionGUID *string) []*DealCashback
	GetByDealCashbackTransactionGUIDAndShoppingListGUID(dealCashbackTransactionGUID *string,
		shoppingListGUID, relations string) []*DealCashback
	GetByUserGUIDAndShoppingListGUIDAndDealGUID(userGUID, shoppingListGUID, dealGUID string) *DealCashback
	CountByDealGUIDAndUserGUID(dealGUID, userGUID string) int
	CountByDealGUID(dealGUID string) int
	CalculateTotalCashbackAmountFromDealCashbackAddedTolist(userGUID string) float64
	GetByUserGUIDShoppingListGUIDAndTransactionStatus(userGUID, shoppingListGUID, transactionStatus, pageNumber,
		pageLimit, relations string) ([]*DealCashback, int)
}
