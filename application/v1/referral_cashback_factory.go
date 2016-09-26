package v1

import (
	"bitbucket.org/cliqers/shoppermate-api/systems"
	"github.com/jinzhu/gorm"
)

type ReferralCashbackFactorysInterface interface {
	CreateReferralCashbackFactory(referrerGUID string, referentGUID string) (interface{}, *systems.ErrorData)
}

type ReferralCashbackFactory struct {
	DB *gorm.DB
}

// CreateReferralCashbackFactory function used to store referral cashback history in database after registration
func (rcf *ReferralCashbackFactory) CreateReferralCashbackFactory(referrerGUID string, referentGUID string) (interface{}, *systems.ErrorData) {
	referralCashback := &ReferralCashback{
		GUID:           Helper.GenerateUUID(),
		ReferrerGUID:   referrerGUID,
		ReferentGUID:   referentGUID,
		CashbackAmount: 5,
	}

	result := rcf.DB.Create(referralCashback)

	if result.Error != nil || result.RowsAffected == 0 {
		rcf.DB.Rollback()
		return nil, Error.InternalServerError(result.Error, systems.DatabaseError)
	}

	return result.Value, nil
}
