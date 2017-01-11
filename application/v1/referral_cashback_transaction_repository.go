package v1

import (
	"bitbucket.org/cliqers/shoppermate-api/systems"
	"github.com/jinzhu/gorm"
)

// ReferralCashbackTransactionRepository will handle all CRUD functions for Referral Cashback Transaction Resources.
type ReferralCashbackTransactionRepository struct {
	DB *gorm.DB
}

// Create function used to store referral cashback history in database after registration
func (rctr *ReferralCashbackTransactionRepository) Create(dbTransaction *gorm.DB, userGUID, referrerGUID, transactionGUID string) (*ReferralCashbackTransaction, *systems.ErrorData) {
	referralCashback := &ReferralCashbackTransaction{
		GUID:            Helper.GenerateUUID(),
		UserGUID:        userGUID,
		ReferrerGUID:    referrerGUID,
		TransactionGUID: transactionGUID,
	}

	result := dbTransaction.Create(referralCashback)

	if result.Error != nil || result.RowsAffected == 0 {
		return nil, Error.InternalServerError(result.Error, systems.DatabaseError)
	}

	return result.Value.(*ReferralCashbackTransaction), nil
}

// CountTotalNumberOfReferralCashbackByUserGUID function used to count total number of referral cashback transaction for userGUID
// by user GUID.
func (rctr *ReferralCashbackTransactionRepository) CountTotalNumberOfReferralCashbackByUserGUID(userGUID string) int64 {
	var count int64

	rctr.DB.Model(&ReferralCashbackTransaction{}).Where(&ReferralCashbackTransaction{UserGUID: userGUID}).Count(&count)

	return count
}
