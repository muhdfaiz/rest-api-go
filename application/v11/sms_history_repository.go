package v11

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
// Return newly created SMS History or internal server error if encountered.
func (shr *SmsHistoryRepository) Create(dbTransaction *gorm.DB, data map[string]string) (interface{}, *systems.ErrorData) {
	smsHistory := &SmsHistory{
		GUID:             data["guid"],
		Provider:         data["provider"],
		Text:             data["text"],
		SmsID:            data["sms_id"],
		RecipientNo:      data["recipient_no"],
		VerificationCode: data["verification_code"],
		Event:            data["event"],
	}

	result := dbTransaction.Create(smsHistory)

	if result.Error != nil || result.RowsAffected == 0 {
		return nil, Error.InternalServerError(result.Error, systems.DatabaseError)
	}

	return result.Value, nil
}

// CountByPhoneNoAndTodayDateAndEventName function used to count total number of SMS History for today date.
// Filter by:
// - recipient_no
// - event
// - created_at = today date
//
// Available events:
// - login
// - update
func (shr *SmsHistoryRepository) CountByPhoneNoAndTodayDateAndEventName(phoneNo, event string) int64 {
	todayDate := time.Now().UTC().Format("2006-01-02")

	var totalNumberOfSmsHistory int64

	shr.DB.Table("sms_histories").Select("count(*) AS total_number_of_sms_history").Where("recipient_no = ? AND date(created_at) = ? AND event = ?", phoneNo, todayDate, event).Count(&totalNumberOfSmsHistory)

	return totalNumberOfSmsHistory
}

// GetLatestByRecipientNoAndEventName function used to retrieve last latest sms history.
// Filter by:
// - recipient_no
// - event
//
// Available events:
// - login
// - update
func (shr *SmsHistoryRepository) GetLatestByRecipientNoAndEventName(recipientNo, eventName string) *SmsHistory {
	smsHistory := &SmsHistory{}

	shr.DB.Where(&SmsHistory{RecipientNo: recipientNo, Event: eventName}).Last(smsHistory)

	return smsHistory
}

// GetByPhoneNoAndVerificationCodeAndEventName function used to retrieve multiple sms history.
// Filter by:
// - recipient_no
// - verification_code
// - event
func (shr *SmsHistoryRepository) GetByPhoneNoAndVerificationCodeAndEventName(phoneNo, verificationCode, eventName string) *SmsHistory {
	smsHistory := &SmsHistory{}

	shr.DB.Where(&SmsHistory{RecipientNo: phoneNo, VerificationCode: verificationCode, Event: eventName}).
		Find(smsHistory)

	return smsHistory
}
