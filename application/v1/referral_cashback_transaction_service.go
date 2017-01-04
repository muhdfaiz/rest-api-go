package v1

import (
	"bitbucket.org/cliqers/shoppermate-api/systems"
	"github.com/jinzhu/gorm"
)

type ReferralCashbackTransactionService struct {
	ReferralCashbackTransactionRepository ReferralCashbackTransactionRepositoryInterface
}

// CreateReferralCashbackTransaction function used to create referral cashback transaction when someone register new
// account and apply other user referral code.
func (rcts *ReferralCashbackTransactionService) CreateReferralCashbackTransaction(dbTransaction *gorm.DB, userGUID, referrerGUID,
	transactionGUID string) (*ReferralCashbackTransaction, *systems.ErrorData) {

	referralCashbackTransaction, error := rcts.ReferralCashbackTransactionRepository.Create(dbTransaction, userGUID, referrerGUID, transactionGUID)

	if error != nil {
		return nil, error
	}

	return referralCashbackTransaction, nil
}

// CountTotalNumberOfUserReferralCashbackTransaction function to count how many times user has been refer by other user during
// account registration.
func (rcts *ReferralCashbackTransactionService) CountTotalNumberOfUserReferralCashbackTransaction(userGUID string) int64 {
	totalNumberUserReferralCashbackTransaction := rcts.ReferralCashbackTransactionRepository.CountTotalNumberOfReferralCashbackByUserGUID(userGUID)

	return totalNumberUserReferralCashbackTransaction
}
