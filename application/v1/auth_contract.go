package v1

import (
	"bitbucket.org/cliqers/shoppermate-api/systems"
)

// AuthServiceInterface is a contract that defines the method needed for Auth Service.
type AuthServiceInterface interface {
	AuthenticateUserViaPhoneNumber(userGUID, phoneNo, debug string) *systems.ErrorData
	AuthenticateUserViaFacebook(userGUID, userPhoneNo, facebookID, deviceUUID string) (*systems.JwtToken, *systems.ErrorData)
	LogoutUser(deviceUUID, userGUID string) *systems.ErrorData
	GenerateJWTTokenForUser(userGUID, userPhoneNo, deviceUUID string) (*systems.JwtToken, *systems.ErrorData)
}
