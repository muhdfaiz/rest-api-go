package v1

import (
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"

	"bitbucket.org/shoppermate-api/systems"
)

type ReferralCashbackFactory struct {
	DB *gorm.DB
}

// CreateSmsHistory function used to store Sms History in database after registration & login
func (rcf *ReferralCashbackFactory) CreateReferralCashbackFactory(referrerGUID string, referentGUID string) (interface{}, *systems.ErrorData) {
	referralCashback := &ReferralCashback{
		GUID:           uuid.NewV4().String(),
		ReferrerGUID:   referrerGUID,
		ReferentGUID:   referentGUID,
		CashbackAmount: 5,
	}

	result := rcf.DB.Create(referralCashback)

	if result.Error != nil || result.RowsAffected == 0 {
		rcf.DB.Rollback()
		return nil, ErrorMesg.InternalServerError(result.Error, systems.DatabaseError)
	}

	return result.Value, nil
}
