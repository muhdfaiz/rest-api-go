package v1

import "bitbucket.org/cliqers/shoppermate-api/systems"

// AuthServiceInterface is a contract that defines the method needed for Auth Service.
type AuthServiceInterface interface {
	AuthenticateUserViaPhoneNumber(phoneNo string) (*User, *systems.ErrorData)
	AuthenticateUserViaFacebook(facebookID string, deviceUUID string) (*User, *systems.JwtToken, *systems.ErrorData)
	LogoutUser(deviceUUID string, userGUID string) *systems.ErrorData
	GenerateJWTTokenForUser(userGUID string, userPhoneNo string, deviceUUID string) (*systems.JwtToken, *systems.ErrorData)
}

type AuthService struct {
	UserService   UserServiceInterface
	SmsService    SmsServiceInterface
	DeviceService DeviceServiceInterface
}

// AuthenticateUserViaPhoneNumber function used to login user using phone number.
func (as *AuthService) AuthenticateUserViaPhoneNumber(phoneNo string) (*User, *systems.ErrorData) {
	user, error := as.UserService.CheckUserPhoneNumberValidOrNot(phoneNo)

	if error != nil {
		return nil, error
	}

	_, error = as.SmsService.SendVerificationCode(phoneNo, user.GUID)

	if error != nil {
		return nil, error
	}

	return user, nil
}

// AuthenticateUserViaFacebook function used to login user using facebook.
func (as *AuthService) AuthenticateUserViaFacebook(facebookID string, deviceUUID string) (*User, *systems.JwtToken, *systems.ErrorData) {
	user, error := as.UserService.CheckUserFacebookIDValidOrNot(facebookID)

	if error != nil {
		return nil, nil, error
	}

	device := as.DeviceService.ViewDeviceByUUIDIncludingSoftDelete(deviceUUID)

	if device.UserGUID == "" {
		device, error = as.DeviceService.UpdateByDeviceUUID(deviceUUID, UpdateDevice{UserGUID: user.GUID})

		if error != nil {
			return nil, nil, error
		}
	}

	error = as.DeviceService.ReactivateDevice(device.GUID)

	if error != nil {
		return nil, nil, error
	}

	jwtToken, error := as.GenerateJWTTokenForUser(user.GUID, user.PhoneNo, deviceUUID)

	if error != nil {
		return nil, nil, error
	}

	return user, jwtToken, nil
}

// LogoutUser function used to logout user from application by soft delete user device.
func (as *AuthService) LogoutUser(deviceUUID string, userGUID string) *systems.ErrorData {
	device := as.DeviceService.ViewDeviceByUUIDandUserGUID(deviceUUID, userGUID)

	if device.UUID == "" {
		return Error.ResourceNotFoundError("Device", "uuid", deviceUUID)
	}

	error := as.DeviceService.DeleteDeviceByUUID(deviceUUID)

	if error != nil {
		return error
	}

	return nil
}

// GenerateJWTTokenForUser function used to generate JWT Token for the user.
func (as *AuthService) GenerateJWTTokenForUser(userGUID string, userPhoneNo string, deviceUUID string) (*systems.JwtToken, *systems.ErrorData) {
	jwt := &systems.Jwt{}

	jwtToken, error := jwt.GenerateToken(userGUID, userPhoneNo, deviceUUID)

	if error != nil {
		return nil, error
	}

	return jwtToken, nil
}
