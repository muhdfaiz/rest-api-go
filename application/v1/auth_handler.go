package v1

import (
	"net/http"
	"strconv"

	"bitbucket.org/cliqers/shoppermate-api/systems"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type AuthHandler struct {
	UserRepository   UserRepositoryInterface
	DeviceRepository DeviceRepositoryInterface
	DeviceFactory    DeviceFactoryInterface
	SmsService       SmsServiceInterface
}

// LoginViaPhone used to login user via phone no.
func (ah *AuthHandler) LoginViaPhone(c *gin.Context) {
	DB := c.MustGet("DB").(*gorm.DB).Begin()
	authData := &LoginViaPhone{}

	// Bind request based on content type and validate request data.
	if err := Binding.Bind(authData, c); err != nil {
		c.JSON(http.StatusBadRequest, err)
		DB.Close()
		return
	}

	// Retrieve user by phone_no.
	user := ah.UserRepository.GetByPhoneNo(authData.PhoneNo, "")

	// If user phone_no empty return error message.
	if user.PhoneNo == "" {
		DB.Rollback().Close()
		c.JSON(http.StatusNotFound, Error.ResourceNotFoundError("User", "phone_no", user.PhoneNo))
		return
	}

	// Send SMS verification code and soft delete device if user change the phone no.
	_, err := ah.SmsService.SendVerificationCode(user.PhoneNo, user.GUID)

	if err != nil {
		DB.Rollback().Close()
		errorCode, _ := strconv.Atoi(err.Error.Status)
		c.JSON(errorCode, err)
		return
	}

	result := make(map[string]string)
	result["user_guid"] = user.GUID
	result["message"] = "Successfully sent sms to " + user.PhoneNo

	DB.Commit().Close()
	c.JSON(http.StatusOK, gin.H{"data": result})
}

// LoginViaFacebook function used to login user via facebook
func (ah *AuthHandler) LoginViaFacebook(c *gin.Context) {
	DB := c.MustGet("DB").(*gorm.DB).Begin()

	authData := &LoginViaFacebook{}

	// Bind request based on content type and validate request data
	if err := Binding.Bind(authData, c); err != nil {
		DB.Close()
		c.JSON(http.StatusBadRequest, err)
		return
	}

	// Retrieve user facebook_id
	user := ah.UserRepository.GetByFacebookID(authData.FacebookID, "")

	// If facebook_id empty return error message
	if user.FacebookID == "" {
		DB.Close()
		c.JSON(http.StatusNotFound, Error.ResourceNotFoundError("User", "facebook_id", authData.FacebookID))
		return
	}

	// Retrieve device by UUID and ignored deleted_at column
	device := ah.DeviceRepository.GetByUUIDUnscoped(authData.DeviceUUID)

	// If Device User GUID empty, update device with User GUID
	if device.UserGUID == "" {
		err := ah.DeviceFactory.Update(authData.DeviceUUID, UpdateDevice{UserGUID: user.GUID})

		if err != nil {
			DB.Rollback().Close()
			c.JSON(http.StatusInternalServerError, err)
			return
		}
	}

	// Reactivate device by set null to deleted_at column in devices table
	result := DB.Unscoped().Model(&Device{}).Update("deleted_at", nil)
	if result.Error != nil {
		DB.Rollback().Close()
		c.JSON(http.StatusInternalServerError, Error.InternalServerError(result.Error, systems.DatabaseError))
		return
	}

	// Generate new JWT Token
	jwt := &systems.Jwt{}
	jwtToken, err := jwt.GenerateToken(user.GUID, user.PhoneNo, device.UUID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
	}

	DB.Commit().Close()

	response := make(map[string]interface{})
	response["user"] = user
	response["access_token"] = jwtToken
	c.JSON(http.StatusOK, gin.H{"data": response})

}

// Refresh function used to refresh device token. Example when user close the app and open the app again,
// app must request this endpoint to avoid token expired
func (ah *AuthHandler) Refresh(c *gin.Context) {
	tokenData := c.MustGet("Token").(map[string]string)

	// Generate new JWT Token
	jwt := &systems.Jwt{}
	jwtToken, err := jwt.GenerateToken(tokenData["user_guid"], tokenData["user_phone_no"], tokenData["device_uuid"])

	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
	}

	c.JSON(http.StatusOK, gin.H{"data": jwtToken})
}

// Logout function used to logout user from application.
// System will soft delete device by set the time value to column deleted_at
func (ah *AuthHandler) Logout(c *gin.Context) {
	DB := c.MustGet("DB").(*gorm.DB).Begin()

	tokenData := c.MustGet("Token").(map[string]string)

	// Retrieve device by UUID and User GUID
	device := ah.DeviceRepository.GetByUUIDAndUserGUID(tokenData["device_uuid"], tokenData["user_guid"])

	// If device uuid empty return error message
	if device.UUID == "" {
		DB.Close()
		c.JSON(http.StatusNotFound, Error.ResourceNotFoundError("Device", "uuid", device.UUID))
		return
	}

	// Soft delete device by set current time to deleted_at columns
	err := ah.DeviceFactory.Delete("uuid", device.UUID)

	if err != nil {
		DB.Rollback().Close()
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	// Response data
	result := make(map[string]string)
	result["message"] = "Successfully logout"

	DB.Commit().Close()
	c.JSON(http.StatusOK, gin.H{"data": result})
}
