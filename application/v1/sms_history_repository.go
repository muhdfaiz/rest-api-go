package v1

import (
	"time"

	"bitbucket.org/cliqers/shoppermate-api/systems"

	"github.com/jinzhu/gorm"
)

// SmsHistoryRepository will handle all CRUD function for Sms History resources.
type SmsHistoryRepository struct {
	DB *gorm.DB
}

// Create function used to create new sms history and store in database.
func (shr *SmsHistoryRepository) Create(data map[string]string) (interface{}, *systems.ErrorData) {
	smsHistory := &SmsHistory{
		GUID:             data["guid"],
		UserGUID:         data["user_guid"],
		Provider:         data["provider"],
		Text:             data["text"],
		SmsID:            data["sms_id"],
		RecipientNo:      data["recipient_no"],
		VerificationCode: data["verification_code"],
		Status:           data["status"],
	}

	result := shr.DB.Create(smsHistory)

	if result.Error != nil || result.RowsAffected == 0 {
		return nil, Error.InternalServerError(result.Error, systems.DatabaseError)
	}

	return result.Value, nil
}

func (shr *SmsHistoryRepository) CountAll(conditionAttribute string, conditionValue string) int64 {
	var count int64
	shr.DB.Model(&SmsHistory{}).Where(conditionAttribute+" = ?", conditionValue).Count(&count)
	return count
}

// GetByRecipientNo function used to retrieve sms history by recipientNo
// Return sms history row data if found
// Return nil if not found
func (shr *SmsHistoryRepository) GetByRecipientNo(recipientNo string) *SmsHistory {
	smsHistory := shr.DB.Where(&SmsHistory{RecipientNo: recipientNo}).Last(&SmsHistory{})
	if smsHistory.RowsAffected != 0 {
		return smsHistory.Value.(*SmsHistory)
	}

	return nil
}

// CalculateIntervalBetweenCurrentTimeAndLastSmsSentTime function used to calculate interval
// between current time and last sms sent time.
// This to protect user continuous click send sms button with short interval.
// The inverval value must be more than 250 second before user can send sms again
func (shr *SmsHistoryRepository) CalculateIntervalBetweenCurrentTimeAndLastSmsSentTime(smsSentTime time.Time) int {
	currentTime := time.Now()
	interval := int(currentTime.Sub(smsSentTime).Seconds())

	return interval
}

// VerifyVerificationCode function used to verify verification code user enter during login & registration
func (shr *SmsHistoryRepository) VerifyVerificationCode(phoneNo string, verificationCode string) *SmsHistory {

	result := shr.DB.Where(&SmsHistory{RecipientNo: phoneNo, VerificationCode: verificationCode}).
		Find(&SmsHistory{})

	if result.RowsAffected == 1 {
		return result.Value.(*SmsHistory)
	}
	return nil
}
