package v1

import "github.com/jinzhu/gorm"

type DealCashbackTransactionRepositoryInterface interface {
	CountByUserGUID(userGUID string) *int
}

type DealCashbackTransactionRepository struct {
	DB *gorm.DB
}

func (dctr *DealCashbackTransactionRepository) CountByUserGUID(userGUID string) *int {
	var totalDealCashback *int

	dctr.DB.Model(&DealCashbackTransaction{}).Where("user_guid = ?", userGUID).Count(&totalDealCashback)

	return totalDealCashback
}
