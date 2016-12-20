package v1

import (
	"bitbucket.org/cliqers/shoppermate-api/systems"
	"github.com/jinzhu/gorm"
)

// ReferralCashbackRepository will handle all CRUD functions for Referral Cashback resource.
type ReferralCashbackRepository struct {
	DB *gorm.DB
}

// Create function used to store referral cashback history in database after registration
func (rcr *ReferralCashbackRepository) Create(referrerGUID string, referentGUID string) (interface{}, *systems.ErrorData) {
	referralCashback := &ReferralCashback{
		GUID:           Helper.GenerateUUID(),
		ReferrerGUID:   referrerGUID,
		ReferentGUID:   referentGUID,
		CashbackAmount: 5,
	}

	result := rcr.DB.Create(referralCashback)

	if result.Error != nil || result.RowsAffected == 0 {
		return nil, Error.InternalServerError(result.Error, systems.DatabaseError)
	}

	return result.Value, nil
}

func (rcr *ReferralCashbackRepository) Count(conditionAttribute string, conditionValue string) int64 {
	var count int64

	rcr.DB.Model(&ReferralCashback{}).Where(conditionAttribute+" = ?", conditionValue).Count(&count)

	return count
}
