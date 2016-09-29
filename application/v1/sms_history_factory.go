package v1

import (
	"time"

	"github.com/jinzhu/gorm"

	"bitbucket.org/cliqers/shoppermate-api/systems"
)

type SmsFactory struct{}

// CreateSmsHistory function used to store Sms History in database after registration & login
func (sf *SmsFactory) CreateSmsHistory(DB *gorm.DB, data map[string]string) (interface{}, *systems.ErrorData) {
	currentTime := time.Now().UTC()
	smsHistory := &SmsHistory{
		GUID:             data["guid"],
		UserGUID:         data["user_guid"],
		Provider:         data["provider"],
		Text:             data["text"],
		SmsID:            data["sms_id"],
		RecipientNo:      data["recipient_no"],
		VerificationCode: data["verification_code"],
		Status:           data["status"],
		CreatedAt:        currentTime,
		UpdatedAt:        currentTime,
	}

	result := DB.Create(smsHistory)

	if result.Error != nil || result.RowsAffected == 0 {
		return nil, Error.InternalServerError(result.Error, systems.DatabaseError)
	}

	return result.Value, nil
}
