package v11

import (
	"time"

	"github.com/jinzhu/gorm"

	"bitbucket.org/cliqers/shoppermate-api/systems"
)

// SmsHistoryServiceInterface is a contract that defines the method needed for SmsHistoryService.
type SmsHistoryServiceInterface interface {
	GetLatestSmsHistoryByPhoneNoAndEventName(phoneNo, eventName string) *SmsHistory
	CheckIfUserReachSmsLimitForToday(phoneNo, eventName string) *systems.ErrorData
	VerifyVerificationCode(phoneNo, verificationCode, eventName string) *systems.ErrorData
	CalculateIntervalBetweenCurrentTimeAndLastSmsSentTime(smsSentTime time.Time) *systems.ErrorData
}

// SmsHistoryRepositoryInterface is a contract that defines the method needed for SmsHistoryRepository.
type SmsHistoryRepositoryInterface interface {
	Create(dbTransaction *gorm.DB, data map[string]string) (interface{}, *systems.ErrorData)
	CountByPhoneNoForTodayDate(phoneNo, eventName string) int64
	GetLatestByRecipientNoAndEventName(recipientNo, eventName string) *SmsHistory
	GetByPhoneNoAndVerificationCodeAndEventName(phoneNo, verificationCode, eventName string) *SmsHistory
}
