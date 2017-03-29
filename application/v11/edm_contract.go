package v11

import (
	"bitbucket.org/cliqers/shoppermate-api/systems"
	"github.com/jinzhu/gorm"
)

// EdmServiceInterface is a contract that defines the method needed for EDM Service.
type EdmServiceInterface interface {
	SendEdmForInsufficientFunds(dbTransaction *gorm.DB, userGUID string, data SendEdmInsufficientFunds) *systems.ErrorData
}
