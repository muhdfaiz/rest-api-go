package v1

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"bitbucket.org/cliqers/shoppermate-api/systems"
)

// SmsHistoryService will handle all application logic related to Sms History resources.
type SmsHistoryService struct {
	SmsHistoryRepository SmsHistoryRepositoryInterface
}

// GetLatestSmsHistoryByPhoneNoAndEventName function used to retrieve latest Sms History by phone number and event name.
func (shs *SmsHistoryService) GetLatestSmsHistoryByPhoneNoAndEventName(phoneNo, eventName string) *SmsHistory {
	smsHistory := shs.SmsHistoryRepository.GetLatestByRecipientNoAndEventName(phoneNo, eventName)

	return smsHistory
}

// CheckIfUserReachSmsLimitForToday function used to check if user already sent 3 SMS for today.
func (shs *SmsHistoryService) CheckIfUserReachSmsLimitForToday(phoneNo, eventName string) *systems.ErrorData {
	totalNumberOfSmsHistory := shs.SmsHistoryRepository.CountByPhoneNoForTodayDate(phoneNo, eventName)

	if totalNumberOfSmsHistory >= 3 {
		return Error.GenericError(strconv.Itoa(http.StatusUnprocessableEntity), systems.ReachLimitSmsSentForToday,
			fmt.Sprintf(systems.TitleReachLimitSmsSentForToday, phoneNo), "message", systems.ErrorReachLimitSmsSentForToday)
	}

	return nil
}

// VerifyVerificationCode function used to verify sms verification code correct or not by
// checking phone number and sms verification in database through Sms History Repository.
func (shs *SmsHistoryService) VerifyVerificationCode(phoneNo, verificationCode, eventName string) *systems.ErrorData {
	smsHistory := shs.SmsHistoryRepository.GetByPhoneNoAndVerificationCodeAndEventName(phoneNo, verificationCode, eventName)

	if smsHistory.GUID == "" {
		error := Error.GenericError(strconv.Itoa(http.StatusBadRequest), systems.VerificationCodeInvalid,
			systems.TitleVerificationCodeInvalid, "", fmt.Sprintf(systems.ErrorVerificationCodeInvalid, verificationCode))

		return error
	}

	return nil
}

// CalculateIntervalBetweenCurrentTimeAndLastSmsSentTime function used to calculate interval
// between current time and last sms sent time.
// This to protect user continuous click send sms button with short interval.
// The inverval value must be more than 250 second before user can send sms again
func (shs *SmsHistoryService) CalculateIntervalBetweenCurrentTimeAndLastSmsSentTime(smsSentTime time.Time) *systems.ErrorData {
	currentTime := time.Now()

	interval := int(currentTime.Sub(smsSentTime).Seconds())

	// If time interval in second below 250 return error message
	if interval < 250 {
		durationUserMustWait := 250 - interval

		error := Error.GenericError("500", systems.FailedToSendSMS, systems.TitleSentSmsError,
			"", fmt.Sprintf(systems.ErrorSentSms, strconv.Itoa(durationUserMustWait)))

		return error
	}

	return nil
}
