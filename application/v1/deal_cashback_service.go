package v1

import "bitbucket.org/cliqers/shoppermate-api/systems"

type DealCashbackServiceInterface interface {
	CountTotalNumberOfDealUserAddToList(userGUID string, dealGUID string) int
	CreateDealCashbackAndShoppingListItem(userGUID string, dealCashbackData CreateDealCashback) *systems.ErrorData
	GetUserDealCashbackForUserShoppingList(userGUID string, shoppingListGUID string, pageNumber string,
		pageLimit string, relations string) ([]*DealCashback, int)
}

type DealCashbackService struct {
	DealRepository          DealRepositoryInterface
	ShoppingListItemFactory ShoppingListItemFactoryInterface
	DealCashbackRepository  DealCashbackRepositoryInterface
	DealCashbackFactory     DealCashbackFactoryInterface
}

// CreateDealCashbackAndShoppingListItem function used to create deal cashback and store new shopping list item based on deal item
func (dcs *DealCashbackService) CreateDealCashbackAndShoppingListItem(userGUID string, dealCashbackData CreateDealCashback) *systems.ErrorData {
	// Create New Deal Cashback
	_, err := dcs.DealCashbackFactory.Create(userGUID, dealCashbackData)

	// Output error if failed to create new device
	if err != nil {
		return err
	}

	deal := dcs.DealRepository.GetDealByGUID(dealCashbackData.DealGUID)

	shoppingListItemData := CreateShoppingListItem{
		UserGUID:         userGUID,
		ShoppingListGUID: dealCashbackData.ShoppingListGUID,
		Name:             deal.Name,
		Quantity:         1,
		AddedFromDeal:    1,
		DealGUID:         dealCashbackData.DealGUID,
	}

	_, err = dcs.ShoppingListItemFactory.CreateForDeal(shoppingListItemData)

	if err != nil {
		return err
	}

	return nil
}

func (dcs *DealCashbackService) CountTotalNumberOfDealUserAddToList(userGUID string, dealGUID string) int {
	total := dcs.DealCashbackRepository.CountByDealGUIDAndUserGUID(dealGUID, userGUID)

	return total
}

func (dcs *DealCashbackService) GetUserDealCashbackForUserShoppingList(userGUID string, shoppingListGUID string, pageNumber string,
	pageLimit string, relations string) ([]*DealCashback, int) {

	userDealCashbacks, totalUserDealCashbacks := dcs.DealCashbackRepository.GetByUserGUIDShoppingListGUIDAndTransactionGUIDEmpty(userGUID, shoppingListGUID, pageNumber, pageLimit, relations)

	return userDealCashbacks, totalUserDealCashbacks
}
