package v1

import (
	"net/http"
	"time"

	"github.com/jinzhu/gorm"

	"bitbucket.org/cliqers/shoppermate-api/systems"
)

// DealCashbackTransformerInterface is a contract that defines the method needed for deal cashback transformer.
type DealCashbackTransformerInterface interface {
	transformCollection(currentURI *http.Request, data interface{}, totalData int, limit string) *DealCashbackResponse
}

// DealCashbackServiceInterface is a contract that defines the method needed for Deal Cashback Service.
type DealCashbackServiceInterface interface {
	CheckDealAlreadyAddedToShoppingList(userGUID, shoppingListGUID, dealGUID string) *systems.ErrorData
	CreateDealCashbackAndShoppingListItem(dbTransaction *gorm.DB, userGUID string, dealCashbackData CreateDealCashback) *systems.ErrorData
	CountTotalNumberOfDealUserAddedToList(userGUID string, dealGUID string) int
	SumTotalAmountOfDealAddedTolistByUser(userGUID string) float64
	GetDealCashbacksByTransactionGUIDAndGroupByShoppingList(dealCashbackTrasnactionGUID string) []*DealCashback
	GetUserDealCashbacksByDealGUID(userGUID, dealGUID, pageNumber, pageLimit, relations string) ([]*DealCashback, int)
	GetUserDealCashbacksByShoppingList(dbTransaction *gorm.DB, userGUID string, shoppingListGUID string, transactionStatus string,
		pageNumber string, pageLimit string, relations string) ([]*DealCashback, int, *systems.ErrorData)
	RemoveDealCashbackIfDealExpiredMoreThan7Days(dbTransaction *gorm.DB, userGUID, shoppingListGUID,
		dealGUID string, dealEndDate time.Time) *systems.ErrorData
}

// DealCashbackRepositoryInterface is a contract that defines the method needed for deal cashback repository.
type DealCashbackRepositoryInterface interface {
	Create(dbTransaction *gorm.DB, userGUID string, data CreateDealCashback) (*DealCashback, *systems.ErrorData)
	UpdateDealCashbackTransactionGUID(dbTransaction *gorm.DB, dealCashbackGUIDs []string, dealCashbackTransactionGUID string) *systems.ErrorData
	DeleteByUserGUIDAndDealGUID(dbTransaction *gorm.DB, userGUID string, dealGUID string) *systems.ErrorData
	DeleteByUserGUIDAndShoppingListGUIDAndDealGUID(dbTransaction *gorm.DB, userGUID string, shoppingListGUID string, dealGUID string) *systems.ErrorData
	GetByGUID(GUID string) *DealCashback
	GetByUserGUIDAndDealGUIDGroupByShoppingList(userGUID, dealGUID, pageNumber, pageLimit, relations string) ([]*DealCashback, int)
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
