package v1

import "bitbucket.org/cliqers/shoppermate-api/systems"

// SmsServiceInterface is a contract that defines the methods needed for Sms Service.
type SmsServiceInterface interface {
	SendVerificationCode(phoneNo string, userGUID string) (interface{}, *systems.ErrorData)
	send(message string, recipientNumber string) (map[string]string, *systems.ErrorData)
}
