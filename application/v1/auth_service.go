package v1

import (
	"bitbucket.org/cliqers/shoppermate-api/systems"
	"github.com/jinzhu/gorm"
)

type AuthService struct {
	SmsService    SmsServiceInterface
	DeviceService DeviceServiceInterface
}

// AuthenticateUserViaPhoneNumber function used to login user using phone number.
func (as *AuthService) AuthenticateUserViaPhoneNumber(dbTransaction *gorm.DB, userGUID, phoneNo, debug string) *systems.ErrorData {
	if debug != "1" {
		_, error := as.SmsService.SendVerificationCode(dbTransaction, phoneNo, userGUID)

		if error != nil {
			return error
		}
	}

	return nil
}

// AuthenticateUserViaFacebook function used to login user using facebook.
func (as *AuthService) AuthenticateUserViaFacebook(dbTransaction *gorm.DB, userGUID, userPhoneNo, facebookID, deviceUUID string) (*systems.JwtToken, *systems.ErrorData) {
	device := as.DeviceService.ViewDeviceByUUIDIncludingSoftDelete(deviceUUID)

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

// LogoutUser function used to logout user from application by soft delete user device.
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
func (as *AuthService) GenerateJWTTokenForUser(userGUID, userPhoneNo, deviceUUID string) (*systems.JwtToken, *systems.ErrorData) {
	jwt := &systems.Jwt{}

	jwtToken, error := jwt.GenerateToken(userGUID, userPhoneNo, deviceUUID, "")

	if error != nil {
		return nil, error
	}

	return jwtToken, nil
}
