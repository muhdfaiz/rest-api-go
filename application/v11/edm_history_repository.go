package v11

import (
	"time"

	"bitbucket.org/cliqers/shoppermate-api/systems"
	"github.com/jinzhu/gorm"
)

type EdmHistoryRepository struct {
	DB *gorm.DB
}

// Create function used to create new edm history and store in database.
func (ehr *EdmHistoryRepository) Create(dbTransaction *gorm.DB, data map[string]string) (*EdmHistory, *systems.ErrorData) {
	edmHistory := &EdmHistory{
		GUID:     data["guid"],
		UserGUID: data["user_guid"],
		Event:    data["event"],
	}

	result := dbTransaction.Create(edmHistory)

	if result.Error != nil || result.RowsAffected == 0 {
		return nil, Error.InternalServerError(result.Error, systems.DatabaseError)
	}

	return result.Value.(*EdmHistory), nil
}

func (ehr *EdmHistoryRepository) GetByUserGUIDAndEventAndCreatedAt(userGUID, event string) *EdmHistory {
	edmHistory := &EdmHistory{}

	todayDate := time.Now().UTC().Format("2006-01-02")

	ehr.DB.Model(&EdmHistory{}).Where("user_guid = ? AND event = ? AND date(created_at) = ?", userGUID, event, todayDate).
		Find(edmHistory)

	return edmHistory
}
