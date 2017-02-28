package v1_1

import (
	"github.com/jinzhu/gorm"

	"bitbucket.org/cliqers/shoppermate-api/systems"
)

type EdmHistoryRepositoryInterface interface {
	Create(dbTransaction *gorm.DB, data map[string]string) (*EdmHistory, *systems.ErrorData)
	GetByUserGUIDAndEventAndCreatedAt(userGUID, event string) *EdmHistory
}
