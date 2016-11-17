package v1

import "github.com/jinzhu/gorm"

type DealCashbackRepositoryInterface interface {
	GetByGUID(GUID string) *DealCashback
	GetByDealGUIDAndUserGUID(dealGUID string, userGUID string) *DealCashback
	CountByDealGUIDAndUserGUID(dealGUID string, userGUID string) int
	CountByDealGUID(dealGUID string) int
	GetByUserGUIDShoppingListGUIDAndTransactionGUIDEmpty(userGUID string, shoppingListGUID string, pageNumber string,
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

func (dcr *DealCashbackRepository) GetByUserGUIDShoppingListGUIDAndTransactionGUIDEmpty(userGUID string, shoppingListGUID string, pageNumber string,
	pageLimit string, relations string) ([]*DealCashback, int) {

	dealCashbacks := []*DealCashback{}

	offset := SetOffsetValue(pageNumber, pageLimit)

	DB := dcr.DB.Model(&DealCashback{})

	if relations != "" {
		DB = LoadRelations(DB, relations)
	}

	DB.Where(DealCashback{UserGUID: userGUID, ShoppingListGUID: shoppingListGUID}).Where("deal_cashback_transaction_guid IS NULL").Offset(offset).Limit(pageLimit).Find(&dealCashbacks)

	var totalDealCashback *int

	dcr.DB.Model(&DealCashback{}).Where(DealCashback{UserGUID: userGUID, ShoppingListGUID: shoppingListGUID}).Where("deal_cashback_transaction_guid IS NULL").Count(&totalDealCashback)

	return dealCashbacks, *totalDealCashback
}
