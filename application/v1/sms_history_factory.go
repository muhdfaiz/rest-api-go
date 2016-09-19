package v1

import (
	"time"

	"github.com/jinzhu/gorm"

	"bitbucket.org/shoppermate/systems"
)

type SmsFactory struct {
	DB *gorm.DB
}

// CreateSmsHistory function used to store Sms History in database after registration & login
func (sf *SmsFactory) CreateSmsHistory(data map[string]string) (interface{}, *systems.ErrorData) {
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

	result := sf.DB.Create(smsHistory)

	if result.Error != nil || result.RowsAffected == 0 {
		sf.DB.Rollback()
		return nil, ErrorMesg.InternalServerError(result.Error, systems.DatabaseError)
	}

	return result.Value, nil
}
