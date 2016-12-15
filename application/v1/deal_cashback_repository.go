package v1

import "github.com/jinzhu/gorm"

// DealCashbackRepositoryInterface is a contract that defines the method needed for deal cashback repository.
type DealCashbackRepositoryInterface interface {
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

type DealCashbackRepository struct {
	DB *gorm.DB
}

// GetByGUID function used to retrieve deal cashback by deal cashback GUID.
func (dcr *DealCashbackRepository) GetByGUID(GUID string) *DealCashback {
	dealCashback := &DealCashback{}

	dcr.DB.Model(&DealCashback{}).Where(DealCashback{GUID: GUID}).Find(&dealCashback)

	return dealCashback
}

// GetByDealGUIDAndUserGUID function used to retrieve deal cashback by deal cashback GUID and user GUID.
func (dcr *DealCashbackRepository) GetByDealGUIDAndUserGUID(dealGUID, userGUID string) *DealCashback {
	dealCashback := &DealCashback{}

	dcr.DB.Model(&DealCashback{}).Where(DealCashback{DealGUID: dealGUID, UserGUID: userGUID}).Find(&dealCashback)

	return dealCashback
}

// GetByDealCashbackTransactionGUIDAndGroupByShoppingListGUID function used to retrieve deal cashback by
// deal cashback transaction GUID and group the result by shopping list.
func (dcr *DealCashbackRepository) GetByDealCashbackTransactionGUIDAndGroupByShoppingListGUID(dealCashbackTransactionGUID *string) []*DealCashback {
	dealCashbacks := []*DealCashback{}

	dcr.DB.Model(&DealCashback{}).Where(DealCashback{DealCashbackTransactionGUID: dealCashbackTransactionGUID}).Group("shopping_list_guid").Find(&dealCashbacks)

	return dealCashbacks
}

// GetByDealCashbackTransactionGUIDAndShoppingListGUID function used to retrieve deal cashback by transaction
// GUID and shopping list GUID.
func (dcr *DealCashbackRepository) GetByDealCashbackTransactionGUIDAndShoppingListGUID(dealCashbackTransactionGUID *string,
	shoppingListGUID, relations string) []*DealCashback {

	dealCashbacks := []*DealCashback{}

	dcr.DB.Model(&DealCashback{}).Where(DealCashback{DealCashbackTransactionGUID: dealCashbackTransactionGUID, ShoppingListGUID: shoppingListGUID}).
		Group("shopping_list_guid").Find(&dealCashbacks)

	return dealCashbacks
}

// GetByUserGUIDAndShoppingListGUIDAndDealGUID function used to retrieve deal cashback by user GUID, shopping list
// GUID and deal GUID.
func (dcr *DealCashbackRepository) GetByUserGUIDAndShoppingListGUIDAndDealGUID(userGUID, shoppingListGUID, dealGUID string) *DealCashback {
	dealCashback := &DealCashback{}

	dcr.DB.Model(&DealCashback{}).Where(DealCashback{UserGUID: userGUID, ShoppingListGUID: shoppingListGUID, DealGUID: dealGUID}).First(&dealCashback)

	return dealCashback
}

// CountByDealGUIDAndUserGUID function used to count total number of deal cashback by deal GUID and user GUID.
func (dcr *DealCashbackRepository) CountByDealGUIDAndUserGUID(dealGUID, userGUID string) int {
	var totalNumberOfUserDealCashback int

	dcr.DB.Model(&DealCashback{}).Where(DealCashback{DealGUID: dealGUID, UserGUID: userGUID}).Count(&totalNumberOfUserDealCashback)

	return totalNumberOfUserDealCashback
}

// CountByDealGUID function used to count total number of deal cashback by deal GUID.
func (dcr *DealCashbackRepository) CountByDealGUID(dealGUID string) int {
	var totalDealCashback int

	dcr.DB.Model(&DealCashback{}).Where(&DealCashback{DealGUID: dealGUID}).Count(&totalDealCashback)

	return totalDealCashback
}

// CalculateTotalCashbackAmountFromDealCashbackAddedTolist function used to sum all of cashback amount for deal cashback already
// added to list.
func (dcr *DealCashbackRepository) CalculateTotalCashbackAmountFromDealCashbackAddedTolist(userGUID string) float64 {

	type DealCashback struct {
		TotalAmountOfCashback float64 `json:"total_amount_of_cashback"`
	}

	dealCashback := &DealCashback{}

	dcr.DB.Model(&DealCashback{}).Select("sum(ads.cashback_amount) as total_amount_of_cashback").
		Joins("left join ads on ads.guid = deal_cashbacks.deal_guid").Where("user_guid = ?", userGUID).
		Where("deal_cashbacks.deal_cashback_transaction_guid IS NULL").Scan(dealCashback)

	return dealCashback.TotalAmountOfCashback
}

func (dcr *DealCashbackRepository) GetByUserGUIDShoppingListGUIDAndTransactionStatus(userGUID, shoppingListGUID, transactionStatus,
	pageNumber, pageLimit, relations string) ([]*DealCashback, int) {

	dealCashbacks := []*DealCashback{}

	offset := SetOffsetValue(pageNumber, pageLimit)

	DB := dcr.DB.Model(&DealCashback{})

	if relations != "" {
		DB = LoadRelations(DB, relations)
	}

	DB = DB.Where(DealCashback{UserGUID: userGUID, ShoppingListGUID: shoppingListGUID})

	if transactionStatus == "pending" || transactionStatus == "approved" || transactionStatus == "reject" || transactionStatus == "partial_success" {
		DB = DB.Joins("LEFT JOIN deal_cashback_transactions ON deal_cashback_transactions.guid = deal_cashback_transaction_guid").
			Joins("LEFT JOIN transactions ON transactions.GUID = deal_cashback_transactions.transaction_guid").
			Joins("LEFT JOIN transaction_statuses ON transaction_statuses.guid = transactions.transaction_status_guid").
			Where("transaction_statuses.slug = ?", transactionStatus)
	} else if transactionStatus == "notempty" {
		DB = DB.Where("deal_cashback_transaction_guid IS NOT NULL")
	} else if transactionStatus == "empty" {
		DB = DB.Where("deal_cashback_transaction_guid IS NULL")
	}

	var totalDealCashback *int

	DB.Offset(offset).Limit(pageLimit).Find(&dealCashbacks)

	DB.Count(&totalDealCashback)

	return dealCashbacks, *totalDealCashback
}
