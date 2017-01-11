package v1_1

import (
	"bitbucket.org/cliqers/shoppermate-api/systems"
	"github.com/jinzhu/gorm"
)

// SmsServiceInterface is a contract that defines the methods needed for Sms Service.
type SmsServiceInterface interface {
	SendVerificationCode(dbTransaction *gorm.DB, phoneNo, eventName string) (interface{}, *systems.ErrorData)
	send(message, recipientNumber string) (map[string]string, *systems.ErrorData)
}
