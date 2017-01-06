package v1

import (
	"time"

	"github.com/jinzhu/gorm"

	"bitbucket.org/cliqers/shoppermate-api/systems"
)

type DealCashbackService struct {
	ShoppingListItemService    ShoppingListItemServiceInterface
	ShoppingListItemRepository ShoppingListItemRepositoryInterface
	DealCashbackRepository     DealCashbackRepositoryInterface
	DealRepository             DealRepositoryInterface
}

// CheckDealAlreadyAddedToShoppingList function used to check if user already added the deal into the same shopping list.
// Business requirement not allowed user to add same deal multiple times into the same shopping list.
func (dcs *DealCashbackService) CheckDealAlreadyAddedToShoppingList(userGUID, shoppingListGUID, dealGUID string) *systems.ErrorData {
	dealCashback := dcs.DealCashbackRepository.GetByUserGUIDAndShoppingListGUIDAndDealGUID(userGUID, shoppingListGUID, dealGUID)

	if dealCashback.GUID != "" {
		return Error.GenericError("409", systems.UserAlreadyAddDealIntoTheShoppingList,
			"Failed to added deal into the shopping list.", "message", "User already add the deal into the shopping list.")
	}

	return nil
}

// CreateDealCashbackAndShoppingListItem function used to create deal cashback and store new shopping list item based on deal item
func (dcs *DealCashbackService) CreateDealCashbackAndShoppingListItem(dbTransaction *gorm.DB, userGUID string, dealCashbackData CreateDealCashback) *systems.ErrorData {
	error := dcs.CheckDealAlreadyAddedToShoppingList(userGUID, dealCashbackData.ShoppingListGUID, dealCashbackData.DealGUID)

	if error != nil {
		return error
	}

	_, error = dcs.DealCashbackRepository.Create(dbTransaction, userGUID, dealCashbackData)

	if error != nil {
		return error
	}

	deal := dcs.DealRepository.GetDealByGUID(dealCashbackData.DealGUID)

	shoppingListItemData := CreateShoppingListItem{
		UserGUID:         userGUID,
		ShoppingListGUID: dealCashbackData.ShoppingListGUID,
		Name:             deal.Name,
		Quantity:         1,
		AddedFromDeal:    1,
		DealGUID:         deal.GUID,
		CashbackAmount:   deal.CashbackAmount,
	}

	_, error = dcs.ShoppingListItemService.CreateUserShoppingListItemAddedFromDeal(dbTransaction, shoppingListItemData)

	if error != nil {
		return error
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

// GetDealCashbacksByTransactionGUIDAndGroupByShoppingList function used to retrieve deal cashback by Transaction GUID
// and group the deal cashback by Shopping List.
func (dcs *DealCashbackService) GetDealCashbacksByTransactionGUIDAndGroupByShoppingList(dealCashbackTransactionGUID string) []*DealCashback {
	dealCashbacksGroupByShoppingList := dcs.DealCashbackRepository.GetByDealCashbackTransactionGUIDAndGroupByShoppingListGUID(&dealCashbackTransactionGUID)

	return dealCashbacksGroupByShoppingList
}

// GetUserDealCashbacksByDealGUID function used to retrieve all deal cashback for user through Deal Cashback Repository.
func (dcs *DealCashbackService) GetUserDealCashbacksByDealGUID(userGUID, dealGUID, pageNumber, pageLimit, relations string) ([]*DealCashback, int) {
	dealCashbacks, totalDealCashbacks := dcs.DealCashbackRepository.GetByUserGUIDAndDealGUIDGroupByShoppingList(userGUID, dealGUID, pageNumber, pageLimit, relations)

	return dealCashbacks, totalDealCashbacks
}

// GetUserDealCashbacksByShoppingList function used to retrieve all deal cashbacks for specific shopping list.
func (dcs *DealCashbackService) GetUserDealCashbacksByShoppingList(dbTransaction *gorm.DB, userGUID string, shoppingListGUID string, transactionStatus string, pageNumber string,
	pageLimit string, relations string) ([]*DealCashback, int, *systems.ErrorData) {

	userDealCashbacks, totalUserDealCashbacks := dcs.DealCashbackRepository.GetByUserGUIDShoppingListGUIDAndTransactionStatus(userGUID, shoppingListGUID,
		transactionStatus, pageNumber, pageLimit, "deals")

	currentDateInGMT8 := time.Now().UTC().Add(time.Hour * 8)

	for _, userDealCashback := range userDealCashbacks {
		diffInDays := currentDateInGMT8.Sub(userDealCashback.Deals.EndDate).Hours() / 24

		// When the deal already expired more than 7 days
		if diffInDays > 7 {
			error := dcs.RemoveDealCashbackAndSetItemDealExpired(dbTransaction, userGUID, shoppingListGUID, userDealCashback.Deals.GUID)

			return nil, 0, error
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

	return userDealCashbacks, totalUserDealCashbacks, nil
}

// RemoveDealCashbackAndSetItemDealExpired function used to soft delete deal cashback that already expired and set the item deal expired.
func (dcs *DealCashbackService) RemoveDealCashbackAndSetItemDealExpired(dbTransaction *gorm.DB, userGUID, shoppingListGUID, dealGUID string) *systems.ErrorData {
	currentDateInGMT8 := time.Now().UTC().Add(time.Hour * 8).Format("2006-01-02")

	deal := dcs.DealRepository.GetDealByGUIDAndValidStartEndDate(dealGUID, currentDateInGMT8)

	if deal.GUID == "" {
		error := dcs.DealCashbackRepository.DeleteByUserGUIDAndShoppingListGUIDAndDealGUID(dbTransaction, userGUID, shoppingListGUID, dealGUID)

		if error != nil {
			return error
		}

		error = dcs.ShoppingListItemRepository.SetDealExpired(dbTransaction, userGUID, shoppingListGUID, dealGUID)

		if error != nil {
			return error
		}
	}

	return nil
}
