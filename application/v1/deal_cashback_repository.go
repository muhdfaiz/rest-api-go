package v1

import "github.com/jinzhu/gorm"

type DealCashbackRepositoryInterface interface {
	GetByGUID(GUID string) *DealCashback
	GetByDealGUIDAndUserGUID(dealGUID string, userGUID string) *DealCashback
	GetByDealCashbackTransactionGUIDAndGroupByShoppingListGUID(dealCashbackTransactionGUID *string) []*DealCashback
	GetByDealCashbackTransactionGUIDAndShoppingListGUID(dealCashbackTransactionGUID *string,
		shoppingListGUID string, relations string) []*DealCashback
	CountByDealGUIDAndUserGUID(dealGUID string, userGUID string) int
	CountByDealGUID(dealGUID string) int
	CalculateTotalCashbackAmountFromDealCashbackAddedTolist(userGUID string) float64
	GetByUserGUIDShoppingListGUIDAndTransactionStatus(userGUID string, shoppingListGUID string, transactionStatus string, pageNumber string,
		pageLimit string, relations string) ([]*DealCashback, int)
}

type DealCashbackRepository struct {
	DB *gorm.DB
}

func (dcr *DealCashbackRepository) GetByGUID(GUID string) *DealCashback {
	dealCashback := &DealCashback{}

	dcr.DB.Model(&DealCashback{}).Where(DealCashback{GUID: GUID}).Find(&dealCashback)

	return dealCashback
}

func (dcr *DealCashbackRepository) GetByDealGUIDAndUserGUID(dealGUID string, userGUID string) *DealCashback {
	dealCashback := &DealCashback{}

	dcr.DB.Model(&DealCashback{}).Where(DealCashback{DealGUID: dealGUID, UserGUID: userGUID}).Find(&dealCashback)

	return dealCashback
}

func (dcr *DealCashbackRepository) GetByDealCashbackTransactionGUIDAndGroupByShoppingListGUID(dealCashbackTransactionGUID *string) []*DealCashback {
	dealCashbacks := []*DealCashback{}

	dcr.DB.Model(&DealCashback{}).Where(DealCashback{DealCashbackTransactionGUID: dealCashbackTransactionGUID}).Group("shopping_list_guid").Find(&dealCashbacks)

	return dealCashbacks
}

func (dcr *DealCashbackRepository) GetByDealCashbackTransactionGUIDAndShoppingListGUID(dealCashbackTransactionGUID *string,
	shoppingListGUID string, relations string) []*DealCashback {

	dealCashbacks := []*DealCashback{}

	dcr.DB.Model(&DealCashback{}).Where(DealCashback{DealCashbackTransactionGUID: dealCashbackTransactionGUID, ShoppingListGUID: shoppingListGUID}).
		Group("shopping_list_guid").Find(&dealCashbacks)

	return dealCashbacks
}

func (dcr *DealCashbackRepository) CountByDealGUIDAndUserGUID(dealGUID string, userGUID string) int {
	var totalNumberOfUserDealCashback int

	dcr.DB.Model(&DealCashback{}).Where(DealCashback{DealGUID: dealGUID, UserGUID: userGUID}).Count(&totalNumberOfUserDealCashback)

	return totalNumberOfUserDealCashback
}

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
		Joins("left join ads on ads.guid = deal_cashbacks.deal_guid").
		Where("deal_cashbacks.deal_cashback_transaction_guid IS NULL").Scan(dealCashback)

	return dealCashback.TotalAmountOfCashback
}

func (dcr *DealCashbackRepository) GetByUserGUIDShoppingListGUIDAndTransactionStatus(userGUID string, shoppingListGUID string,
	transactionStatus string, pageNumber string, pageLimit string, relations string) ([]*DealCashback, int) {

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
