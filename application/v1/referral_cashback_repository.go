package v1

import "github.com/jinzhu/gorm"

var (
	ReferralCashbackModel = &ReferralCashback{}
)

type ReferralCashbackRepositoryInterface interface {
	Count(DB *gorm.DB, conditionAttribute string, conditionValue string) int64
}

type ReferralCashbackRepository struct{}

func (rcr *ReferralCashbackRepository) Count(DB *gorm.DB, conditionAttribute string, conditionValue string) int64 {
	var count int64
	DB.Model(ReferralCashbackModel).Where(conditionAttribute+" = ?", conditionValue).Count(&count)
	return count
}
