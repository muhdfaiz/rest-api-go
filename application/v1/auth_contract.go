package v1

import (
	"bitbucket.org/cliqers/shoppermate-api/systems"
	"github.com/jinzhu/gorm"
)

// AuthServiceInterface is a contract that defines the method needed for Auth Service.
type AuthServiceInterface interface {
	AuthenticateUserViaPhoneNumber(dbTransaction *gorm.DB, userGUID, phoneNo, debug string) *systems.ErrorData
	AuthenticateUserViaFacebook(dbTransaction *gorm.DB, userGUID, userPhoneNo, facebookID, deviceUUID string) (*systems.JwtToken, *systems.ErrorData)
	LogoutUser(dbTransaction *gorm.DB, deviceUUID, userGUID string) *systems.ErrorData
	GenerateJWTTokenForUser(userGUID, userPhoneNo, deviceUUID string) (*systems.JwtToken, *systems.ErrorData)
}
