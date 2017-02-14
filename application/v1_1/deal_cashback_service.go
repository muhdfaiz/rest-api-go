package v1_1

import (
	"time"

	"github.com/jinzhu/gorm"

	"bitbucket.org/cliqers/shoppermate-api/systems"
)

// DealCashbackService will handle all application logic related to Deal Cashback resources.
type DealCashbackService struct {
	ShoppingListService     ShoppingListServiceInterface
	ShoppingListItemService ShoppingListItemServiceInterface
	DealCashbackRepository  DealCashbackRepositoryInterface
	DealRepository          DealRepositoryInterface
}

// CheckDealAlreadyAddedToShoppingList function used to check if user already added the deal into the same shopping list.
// Business requirement not allowed user to add same deal multiple times into the same shopping list.
func (dcs *DealCashbackService) CheckDealAlreadyAddedToShoppingList(userGUID, shoppingListGUID, dealGUID string) *systems.ErrorData {
	dealCashback := dcs.DealCashbackRepository.GetByUserGUIDAndShoppingListGUIDAndDealGUID(userGUID, shoppingListGUID, dealGUID)

	if dealCashback.GUID != "" {
		return Error.GenericError("409", systems.UserAlreadyAddDealIntoShoppingList,
			systems.TitleUserAlreadyAddDealIntoShoppingList, "message", systems.ErrorUserAlreadyAddDealIntoShoppingList)
	}

	return nil
}

// CreateDealCashbackAndShoppingListItem function used to create deal cashback and store new shopping list item based on deal item
func (dcs *DealCashbackService) CreateDealCashbackAndShoppingListItem(dbTransaction *gorm.DB, userGUID string, dealCashbackData CreateDealCashback) *systems.ErrorData {
	currentDateInGMT8 := time.Now().UTC().Add(time.Hour * 8).Format("2006-01-02")

	deal := dcs.DealRepository.GetDealByGUIDAndUserGUIDWithinDateRangeAndValidQuotaAndLimitPerUserAndPublished(userGUID, dealCashbackData.DealGUID, currentDateInGMT8)

	if deal.GUID == "" {
		return Error.GenericError("422", systems.DealAlreadyExpiredOrNotValid,
			systems.TitleDealAlreadyExpiredOrNotValid, "message", systems.ErrorDealAlreadyExpiredOrNotValid)
	}

	error := dcs.CheckDealAlreadyAddedToShoppingList(userGUID, dealCashbackData.ShoppingListGUID, dealCashbackData.DealGUID)

	if error != nil {
		return error
	}

	_, error = dcs.DealCashbackRepository.Create(dbTransaction, userGUID, dealCashbackData)

	if error != nil {
		return error
	}

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

// GetUserDealCashbacksFilterByTransactionStatusGroupByShoppingList function used to retrieve all deal cashback for user for
// specific transaction status like `empty`, `notempty`, `pending`, `approved` and group by shopping list.
func (dcs *DealCashbackService) GetUserDealCashbacksFilterByTransactionStatusGroupByShoppingList(dbTransaction *gorm.DB, userGUID,
	transactionStatus, relations string) ([]*ShoppingList, *systems.ErrorData) {

	userShoppingListsWithDealCashbacks := []*ShoppingList{}

	userDealCashbacksGroupbyShoppingList := dcs.DealCashbackRepository.GetByUserGUIDAndTransactionStatusGroupByShoppingListGUID(userGUID, transactionStatus)

	if transactionStatus == "empty" {
		for _, userDealCashbackGroupbyShoppingList := range userDealCashbacksGroupbyShoppingList {
			dealCashbacks, _ := dcs.DealCashbackRepository.GetByUserGUIDShoppingListGUIDAndTransactionStatus(userGUID,
				userDealCashbackGroupbyShoppingList.ShoppingListGUID, transactionStatus, "1", "", "deals")

			for _, dealCashback := range dealCashbacks {

				error := dcs.RemoveDealCashbackIfDealExpiredMoreThan7Days(dbTransaction, userGUID, dealCashback.ShoppingListGUID,
					dealCashback.Deals.GUID, dealCashback.Deals.EndDate)

				if error != nil {
					dbTransaction.Rollback()
					return nil, error
				}
			}
		}
	}

	dbTransaction.Commit()

	dealCashbacksForOtherShoppingLists := []*DealCashback{}

	for _, userDealCashbackGroupbyShoppingList := range userDealCashbacksGroupbyShoppingList {
		relations = "deals"

		if relations != "" {
			relations = relations + ",deals"
		}

		dealCashbacks, _ := dcs.DealCashbackRepository.GetByUserGUIDShoppingListGUIDAndTransactionStatus(userGUID,
			userDealCashbackGroupbyShoppingList.ShoppingListGUID, transactionStatus, "1", "", relations)

		dealCashbacks = dcs.SetDealsExpiredInfo(dealCashbacks)

		shoppingList := dcs.ShoppingListService.ViewShoppingListByGUIDIncludingSoftDelete(userDealCashbackGroupbyShoppingList.ShoppingListGUID, "")

		if len(dealCashbacks) > 0 {
			if shoppingList.DeletedAt != nil {
				dealCashbacksForOtherShoppingLists = append(dealCashbacksForOtherShoppingLists, dealCashbacks...)
			} else {
				shoppingList.Dealcashbacks = dealCashbacks

				userShoppingListsWithDealCashbacks = append(userShoppingListsWithDealCashbacks, shoppingList)
			}
		}
	}

	if len(dealCashbacksForOtherShoppingLists) > 0 {
		otherShoppingListWithDealCashbacks := &ShoppingList{}

		otherShoppingListWithDealCashbacks.Name = "Deleted Shopping List"

		otherShoppingListWithDealCashbacks.Dealcashbacks = dealCashbacksForOtherShoppingLists

		userShoppingListsWithDealCashbacks = append(userShoppingListsWithDealCashbacks, otherShoppingListWithDealCashbacks)
	}

	return userShoppingListsWithDealCashbacks, nil
}

// GetUserDealCashbacksByDealGUID function used to retrieve all deal cashback for user through Deal Cashback Repository.
func (dcs *DealCashbackService) GetUserDealCashbacksByDealGUID(userGUID, dealGUID, pageNumber, pageLimit, relations string) ([]*DealCashback, int) {
	dealCashbacks, totalDealCashbacks := dcs.DealCashbackRepository.GetByUserGUIDAndDealGUIDGroupByShoppingList(userGUID, dealGUID, pageNumber, pageLimit, relations)

	return dealCashbacks, totalDealCashbacks
}

// GetUserDealCashbacksByShoppingList function used to retrieve all deal cashbacks for specific shopping list.
func (dcs *DealCashbackService) GetUserDealCashbacksByShoppingList(dbTransaction *gorm.DB, userGUID, shoppingListGUID, transactionStatus, pageNumber,
	pageLimit, relations string) ([]*DealCashback, int, *systems.ErrorData) {

	userDealCashbacks, totalUserDealCashbacks := dcs.DealCashbackRepository.GetByUserGUIDShoppingListGUIDAndTransactionStatus(userGUID, shoppingListGUID,
		transactionStatus, pageNumber, pageLimit, "deals")

	if transactionStatus == "empty" {
		for _, userDealCashback := range userDealCashbacks {

			error := dcs.RemoveDealCashbackIfDealExpiredMoreThan7Days(dbTransaction, userGUID, shoppingListGUID,
				userDealCashback.Deals.GUID, userDealCashback.Deals.EndDate)

			if error != nil {
				dbTransaction.Rollback()
				return nil, 0, error
			}
		}
	}

	dbTransaction.Commit()

	relations = "deals"

	if relations != "" {
		relations = relations + ",deals"
	}

	userDealCashbacks, totalUserDealCashbacks = dcs.DealCashbackRepository.GetByUserGUIDShoppingListGUIDAndTransactionStatus(userGUID, shoppingListGUID,
		transactionStatus, pageNumber, pageLimit, relations)

	userDealCashbacks = dcs.SetDealsExpiredInfo(userDealCashbacks)

	return userDealCashbacks, totalUserDealCashbacks, nil
}

// RemoveDealCashbackIfDealExpiredMoreThan7Days function used to soft delete deal cashback that already expired.
func (dcs *DealCashbackService) RemoveDealCashbackIfDealExpiredMoreThan7Days(dbTransaction *gorm.DB, userGUID, shoppingListGUID,
	dealGUID string, dealEndDate time.Time) *systems.ErrorData {

	currentDateInGMT8 := time.Now().UTC().Add(time.Hour * 8)

	diffInDays := currentDateInGMT8.Sub(dealEndDate).Hours() / 24

	if diffInDays > 7 {
		error := dcs.DealCashbackRepository.DeleteByUserGUIDAndShoppingListGUIDAndDealGUID(dbTransaction, userGUID, shoppingListGUID, dealGUID)

		if error != nil {
			return error
		}
	}

	return nil
}

// SetDealsExpiredInfo function used to set expired to true and set remaining days the deal cashback
// will be remove.
func (dcs *DealCashbackService) SetDealsExpiredInfo(dealCashbacks []*DealCashback) []*DealCashback {
	currentDateInGMT8 := time.Now().UTC().Add(time.Hour * 8)

	for key, dealCashback := range dealCashbacks {

		diffInDays := currentDateInGMT8.Sub(dealCashback.Deals.EndDate).Hours() / 24

		// When the deal already expired
		if diffInDays > 0 {
			dealCashbacks[key].Expired = 1
			dealCashbacks[key].RemainingDaysToRemove = int(diffInDays)
		}
	}

	return dealCashbacks
}
