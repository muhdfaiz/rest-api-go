package v1

import (
	"time"

	"github.com/jinzhu/gorm"

	"bitbucket.org/cliqers/shoppermate-api/systems"
)

// SmsHistoryRepositoryInterface is a contract that defines the method needed for SmsHistoryRepository.
type SmsHistoryRepositoryInterface interface {
	Create(dbTransaction *gorm.DB, data map[string]string) (interface{}, *systems.ErrorData)
	CountAll(conditionAttribute string, conditionValue string) int64
	GetByRecipientNo(recipientNo string) *SmsHistory
	CalculateIntervalBetweenCurrentTimeAndLastSmsSentTime(smsSentTime time.Time) int
	VerifyVerificationCode(phoneNo string, verificationCode string) *SmsHistory
}
