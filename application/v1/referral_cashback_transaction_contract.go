package v1

import (
	"bitbucket.org/cliqers/shoppermate-api/systems"
)

// ReferralCashbackTransactionServiceInterface is a contract that defines the method needed for
// Referral Cashback Transaction Service
type ReferralCashbackTransactionServiceInterface interface {
	CreateReferralCashbackTransaction(userGUID, referrerGUID,
		transactionGUID string) (*ReferralCashbackTransaction, *systems.ErrorData)
	CountTotalNumberOfUserReferralCashbackTransaction(userGUID string) int64
}

// ReferralCashbackTransactionRepositoryInterface is a contract that defines the method needed for
// Referral Cashback Transaction Repository
type ReferralCashbackTransactionRepositoryInterface interface {
	Create(userGUID, referrerGUID, transactionGUID string) (*ReferralCashbackTransaction, *systems.ErrorData)
	CountTotalNumberOfReferralCashbackByUserGUID(userGUID string) int64
}
