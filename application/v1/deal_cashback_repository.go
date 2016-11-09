package v1

import "github.com/jinzhu/gorm"

type DealCashbackRepositoryInterface interface {
	GetByDealGUIDAndUserGUID(dealGUID string, userGUID string) *DealCashback
	CountByDealGUIDAndUserGUID(dealGUID string, userGUID string) int
	CountByDealGUID(dealGUID string) int
}

type DealCashbackRepository struct {
	DB *gorm.DB
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
