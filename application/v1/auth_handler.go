package v1

import (
	"fmt"
	"net/http"
	"strconv"

	"bitbucket.org/shoppermate-api/systems"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type AuthHandler struct{}

func (ah *AuthHandler) LoginViaPhone(c *gin.Context) {
	db := c.MustGet("DB").(*gorm.DB)
	tx := db.Begin()

	c.Writer.Header().Set("Content-Type", "application/json; charset=utf-8")

	authData := &LoginViaPhone{}

	// Bind request based on content type and validate request data
	if err := Binding.Bind(authData, c); err != nil {
		statusCode, _ := strconv.Atoi(err.Error.Status)
		c.JSON(statusCode, err)
		return
	}

	userRepository := &UserRepository{DB: tx}
	user := userRepository.GetByPhoneNo(authData.PhoneNo)

	if user.PhoneNo == "" {
		c.JSON(http.StatusBadRequest, ErrorMesg.GenericError(strconv.Itoa(http.StatusBadRequest), systems.ResourceNotFound,
			fmt.Sprintf(systems.TitleResourceNotFoundError, "User"), "message",
			fmt.Sprintf(systems.ErrorResourceNotFound, "User", "phone_no", user.PhoneNo)))
		return
	}

	// Send SMS verification code and soft delete device if user change the phone no
	smsService := &SmsService{DB: tx}
	_, err := smsService.SendVerificationCode(user.PhoneNo, user.GUID)

	if err != nil {
		errorCode, _ := strconv.Atoi(err.Error.Status)
		c.JSON(errorCode, err)
		return
	}
	tx.Commit()
	result := make(map[string]string)
	result["message"] = "Successfully sent sms to " + user.PhoneNo
	c.JSON(http.StatusOK, gin.H{"data": result})
}

func (ah *AuthHandler) LoginViaFacebook(c *gin.Context) {
	db := c.MustGet("DB").(*gorm.DB)
	tx := db.Begin()

	c.Writer.Header().Set("Content-Type", "application/json; charset=utf-8")

	authData := &LoginViaFacebook{}

	// Bind request based on content type and validate request data
	if err := Binding.Bind(authData, c); err != nil {
		statusCode, _ := strconv.Atoi(err.Error.Status)
		c.JSON(statusCode, err)
		return
	}

	userRepository := &UserRepository{DB: tx}
	user := userRepository.GetFacebookID(authData.FacebookID)

	if user.GUID == "" {
		c.JSON(http.StatusBadRequest, ErrorMesg.GenericError(strconv.Itoa(http.StatusBadRequest), systems.ResourceNotFound,
			fmt.Sprintf(systems.TitleResourceNotFoundError, "User"), "message",
			fmt.Sprintf(systems.ErrorResourceNotFound, "User", "facebook_id", authData.FacebookID)))
		return
	}

	deviceRepository := DeviceRepository{DB: tx}
	device := deviceRepository.GetByUUIDAndUserGUIDUnscoped(authData.DeviceUUID, user.GUID)

	// Return error message if device uuid not exist
	if device.UUID == "" {
		c.JSON(http.StatusBadRequest, ErrorMesg.GenericError(strconv.Itoa(http.StatusBadRequest), systems.ResourceNotFound,
			fmt.Sprintf(systems.TitleResourceNotFoundError, "Device"), "message",
			fmt.Sprintf(systems.ErrorResourceNotFound, "Device", "uuid", device.UUID)))
		return
	}

	result := tx.Unscoped().Model(&Device{}).Update("deleted_at", nil)
	if result.Error != nil || result.RowsAffected == 0 {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, ErrorMesg.InternalServerError(result.Error, systems.DatabaseError))
		return
	}

	jwt := &systems.Jwt{}
	jwtToken, err := jwt.GenerateJWTToken(user.GUID, user.PhoneNo, device.UUID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
	}

	tx.Commit()
	c.JSON(http.StatusOK, gin.H{"data": jwtToken})

}

func (ah *AuthHandler) Logout(c *gin.Context) {
	db := c.MustGet("DB").(*gorm.DB)
	tx := db.Begin()

	tokenData := c.MustGet("Token").(map[string]string)

	deviceRepository := DeviceRepository{DB: tx}
	device := deviceRepository.GetByUUIDAndUserGUID(tokenData["device_uuid"], tokenData["user_guid"])

	// Return error message if device uuid not exist
	if device.UUID == "" {
		c.JSON(http.StatusBadRequest, ErrorMesg.GenericError(strconv.Itoa(http.StatusBadRequest), systems.ResourceNotFound,
			fmt.Sprintf(systems.TitleResourceNotFoundError, "Device"), "message",
			fmt.Sprintf(systems.ErrorResourceNotFound, "Device", "uuid", device.UUID)))
		return
	}

	deviceFactory := &DeviceFactory{DB: tx}
	err := deviceFactory.Delete("uuid", device.UUID)

	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	tx.Commit()
	result := make(map[string]string)
	result["message"] = "Successfully logout"
	c.JSON(http.StatusOK, gin.H{"data": result})
}
