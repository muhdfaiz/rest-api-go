package v1

import (
	"bitbucket.org/cliqers/shoppermate-api/systems"
)

// AuthServiceInterface is a contract that defines the method needed for Auth Service.
type AuthServiceInterface interface {
	AuthenticateUserViaPhoneNumber(phoneNo string, debug string) (*User, *systems.ErrorData)
	AuthenticateUserViaFacebook(facebookID string, deviceUUID string) (*User, *systems.JwtToken, *systems.ErrorData)
	LogoutUser(deviceUUID string, userGUID string) *systems.ErrorData
	GenerateJWTTokenForUser(userGUID string, userPhoneNo string, deviceUUID string) (*systems.JwtToken, *systems.ErrorData)
}
