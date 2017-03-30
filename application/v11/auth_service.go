package v11

import (
	"bitbucket.org/cliqers/shoppermate-api/systems"
	"github.com/jinzhu/gorm"
)

// AuthService used to handle application logic related to Auth resource.
type AuthService struct {
	SmsService    SmsServiceInterface
	DeviceService DeviceServiceInterface
}

// AuthenticateUserViaPhoneNumber function used to authenticate user using phone number.
// It will send SMS to the phone no.
// It will skip send sms if debug parameter exist and equal to 1. Useful during development and testing.
func (as *AuthService) AuthenticateUserViaPhoneNumber(dbTransaction *gorm.DB, phoneNo, debug string) *systems.ErrorData {
	if debug != "1" {
		_, error := as.SmsService.SendVerificationCode(dbTransaction, phoneNo, "login")

		if error != nil {
			return error
		}
	}

	return nil
}

// AuthenticateUserViaFacebook function used to authenticate user using facebook.
// First, it will check if device UUID exist or not in database. If not exist, it will return an error.
// Then, if device user GUID equal to nil, it will used value for userGUID parameter to update device user GUID.
// Then, it will reactive device by setting `deleted_at` field to nil in device table.
// Lastly, it will generate JWT Token based on user GUID, user phone number and device uuid parameter passed to this function.
func (as *AuthService) AuthenticateUserViaFacebook(dbTransaction *gorm.DB, userGUID, userPhoneNo, facebookID, deviceUUID string) (*systems.JwtToken, *systems.ErrorData) {
	device := as.DeviceService.ViewDeviceByUUIDIncludingSoftDelete(deviceUUID)

	if device.UUID == "" {
		return nil, Error.ResourceNotFoundError("Device", "uuid", deviceUUID)
	}

	if device.UserGUID == nil {
		_, error := as.DeviceService.UpdateByDeviceUUID(dbTransaction, deviceUUID, UpdateDevice{UserGUID: userGUID})

		if error != nil {
			return nil, error
		}
	}

	error := as.DeviceService.ReactivateDevice(dbTransaction, device.GUID)

	if error != nil {
		return nil, error
	}

	jwtToken, error := as.GenerateJWTTokenForUser(userGUID, userPhoneNo, deviceUUID)

	if error != nil {
		return nil, error
	}

	return jwtToken, nil
}

// LogoutUser function used to logout user from application.
// First, it will check if device uuid exist or not in database. If not exist return an error.
// Then, it will logout user by soft delete the device based on device uuid.
// Soft delete means it's will set the current date and time to `deleted_at` field in device table.
// It will return any error encountered or return nil.
func (as *AuthService) LogoutUser(dbTransaction *gorm.DB, deviceUUID, userGUID string) *systems.ErrorData {
	device := as.DeviceService.ViewDeviceByUUIDandUserGUID(deviceUUID, userGUID)

	if device.UUID == "" {
		return Error.ResourceNotFoundError("Device", "uuid", deviceUUID)
	}

	error := as.DeviceService.DeleteDeviceByUUID(dbTransaction, deviceUUID)

	if error != nil {
		return error
	}

	return nil
}

// GenerateJWTTokenForUser function used to generate JWT Token for the user.
// It will assign user GUID, user phone number and device uuid to the token payload.
// When you decrypt the token, you can see user GUID assign to `aud` & `sub` keys, device uuid assign to `jti` key.
// Example payload when you decrypt the token:
// {
//  "phone_no": "60174862127",
//  "aud": "8c2e6ea5-5c56-5050-ae37-a44b88e612a7",
//  "exp": 1491365635,
//  "jti": "FFB5AD2D4FB2BFF2E63FEA05CE989995",
//  "iat": 1490760835,
//  "iss": "http://api.shoppermate-api.com",
//  "nbf": 1490760835,
//  "sub": "8c2e6ea5-5c56-5050-ae37-a44b88e612a7"
// }
func (as *AuthService) GenerateJWTTokenForUser(userGUID, userPhoneNo, deviceUUID string) (*systems.JwtToken, *systems.ErrorData) {
	jwt := &systems.Jwt{}

	jwtToken, error := jwt.GenerateToken(userGUID, userPhoneNo, deviceUUID, "")

	if error != nil {
		return nil, error
	}

	return jwtToken, nil
}
