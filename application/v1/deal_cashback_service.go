package v1

import (
	"time"

	"bitbucket.org/cliqers/shoppermate-api/systems"
)

type DealCashbackServiceInterface interface {
	CountTotalNumberOfDealUserAddedToList(userGUID string, dealGUID string) int
	SumTotalAmountOfDealAddedTolistByUser(userGUID string) float64
	GetDealCashbacksByTransactionGUIDAndGroupByShoppingList(dealCashbackTrasnactionGUID string) []*DealCashback
	CreateDealCashbackAndShoppingListItem(userGUID string, dealCashbackData CreateDealCashback) *systems.ErrorData
	GetUserDealCashbackForUserShoppingList(userGUID string, shoppingListGUID string, transactionStatus string,
		pageNumber string, pageLimit string, relations string) ([]*DealCashback, int)
}

type DealCashbackService struct {
	ShoppingListItemService ShoppingListItemServiceInterface
	DealService             DealServiceInterface
	DealCashbackRepository  DealCashbackRepositoryInterface
	DealCashbackFactory     DealCashbackFactoryInterface
}

// CreateDealCashbackAndShoppingListItem function used to create deal cashback and store new shopping list item based on deal item
func (dcs *DealCashbackService) CreateDealCashbackAndShoppingListItem(userGUID string, dealCashbackData CreateDealCashback) *systems.ErrorData {
	_, err := dcs.DealCashbackFactory.Create(userGUID, dealCashbackData)

	if err != nil {
		return err
	}

	deal := dcs.DealService.GetDealByGUID(dealCashbackData.DealGUID)

	shoppingListItemData := CreateShoppingListItem{
		UserGUID:         userGUID,
		ShoppingListGUID: dealCashbackData.ShoppingListGUID,
		Name:             deal.Name,
		Quantity:         1,
		AddedFromDeal:    1,
		DealGUID:         deal.GUID,
		CashbackAmount:   deal.CashbackAmount,
	}

	_, err = dcs.ShoppingListItemService.CreateUserShoppingListItemAddedFromDeal(shoppingListItemData)

	if err != nil {
		return err
	}

	return nil
}

// CountTotalNumberOfDealUserAddedToList function used to count total number of deal added to list by user.
func (dcs *DealCashbackService) CountTotalNumberOfDealUserAddedToList(userGUID string, dealGUID string) int {
	total := dcs.DealCashbackRepository.CountByDealGUIDAndUserGUID(dealGUID, userGUID)

	return total
}

// SumTotalAmountOfDealAddedTolistByUser function used to sum total amount of deal added to list.
func (dcs *DealCashbackService) SumTotalAmountOfDealAddedTolistByUser(userGUID string) float64 {
	totalCashbackAmount := dcs.DealCashbackRepository.CalculateTotalCashbackAmountFromDealCashbackAddedTolist(userGUID)

	return totalCashbackAmount
}

func (dcs *DealCashbackService) GetDealCashbacksByTransactionGUIDAndGroupByShoppingList(dealCashbackTransactionGUID string) []*DealCashback {
	dealsCashbacksGroupByShoppingList := dcs.DealCashbackRepository.GetByDealCashbackTransactionGUIDAndGroupByShoppingListGUID(&dealCashbackTransactionGUID)

	return dealsCashbacksGroupByShoppingList
}

func (dcs *DealCashbackService) GetUserDealCashbackForUserShoppingList(userGUID string, shoppingListGUID string, transactionStatus string, pageNumber string,
	pageLimit string, relations string) ([]*DealCashback, int) {

	userDealCashbacks, totalUserDealCashbacks := dcs.DealCashbackRepository.GetByUserGUIDShoppingListGUIDAndTransactionStatus(userGUID, shoppingListGUID,
		transactionStatus, pageNumber, pageLimit, "deals")

	currentDateInGMT8 := time.Now().UTC().Add(time.Hour * 8)

	for _, userDealCashback := range userDealCashbacks {
		diffInDays := currentDateInGMT8.Sub(userDealCashback.Deals.EndDate).Hours() / 24

		// When the deal already expired more than 7 days
		if diffInDays > 7 {
			dcs.DealService.RemoveDealCashbackAndSetItemDealExpired(userGUID, shoppingListGUID, userDealCashback.Deals.GUID)
		}

	}

	relations = relations + ",deals"

	userDealCashbacks, totalUserDealCashbacks = dcs.DealCashbackRepository.GetByUserGUIDShoppingListGUIDAndTransactionStatus(userGUID, shoppingListGUID,
		transactionStatus, pageNumber, pageLimit, relations)

	for key, userDealCashback := range userDealCashbacks {
		diffInDays := currentDateInGMT8.Sub(userDealCashback.Deals.EndDate).Hours() / 24

		// When the deal already expired not more than 7 days
		if diffInDays > 0 && diffInDays < 7 {
			userDealCashbacks[key].Expired = 1
		}
	}

	return userDealCashbacks, totalUserDealCashbacks
}
