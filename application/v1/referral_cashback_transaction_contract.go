package v1

import (
	"bitbucket.org/cliqers/shoppermate-api/systems"
	"github.com/jinzhu/gorm"
)

// ReferralCashbackTransactionServiceInterface is a contract that defines the method needed for
// Referral Cashback Transaction Service
type ReferralCashbackTransactionServiceInterface interface {
	CreateReferralCashbackTransaction(dbTransaction *gorm.DB, userGUID, referrerGUID,
		transactionGUID string) (*ReferralCashbackTransaction, *systems.ErrorData)
	CountTotalNumberOfUserReferralCashbackTransaction(userGUID string) int64
}

// ReferralCashbackTransactionRepositoryInterface is a contract that defines the method needed for
// Referral Cashback Transaction Repository
type ReferralCashbackTransactionRepositoryInterface interface {
	Create(dbTransaction *gorm.DB, userGUID, referrerGUID, transactionGUID string) (*ReferralCashbackTransaction, *systems.ErrorData)
	CountTotalNumberOfReferralCashbackByUserGUID(userGUID string) int64
}
