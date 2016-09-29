package v1

import (
	"bitbucket.org/cliqers/shoppermate-api/systems"
	"github.com/jinzhu/gorm"
)

type ReferralCashbackFactorysInterface interface {
	CreateReferralCashbackFactory(DB *gorm.DB, referrerGUID string, referentGUID string) (interface{}, *systems.ErrorData)
}

type ReferralCashbackFactory struct{}

// CreateReferralCashbackFactory function used to store referral cashback history in database after registration
func (rcf *ReferralCashbackFactory) CreateReferralCashbackFactory(DB *gorm.DB, referrerGUID string, referentGUID string) (interface{}, *systems.ErrorData) {
	referralCashback := &ReferralCashback{
		GUID:           Helper.GenerateUUID(),
		ReferrerGUID:   referrerGUID,
		ReferentGUID:   referentGUID,
		CashbackAmount: 5,
	}

	result := DB.Create(referralCashback)

	if result.Error != nil || result.RowsAffected == 0 {
		return nil, Error.InternalServerError(result.Error, systems.DatabaseError)
	}

	return result.Value, nil
}
