package v1

import "github.com/jinzhu/gorm"

var (
	ReferralCashbackModel = &ReferralCashback{}
)

type ReferralCashbackRepositoryInterface interface {
	Count(conditionAttribute string, conditionValue string) int64
}

type ReferralCashbackRepository struct {
	DB *gorm.DB
}

func (rcr *ReferralCashbackRepository) Count(conditionAttribute string, conditionValue string) int64 {
	var count int64
	rcr.DB.Model(ReferralCashbackModel).Where(conditionAttribute+" = ?", conditionValue).Count(&count)
	return count
}
