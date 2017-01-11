package v1

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// AuthHandler will handle all request related to Auth Resources.
type AuthHandler struct {
	AuthService AuthServiceInterface
	UserService UserServiceInterface
}

// LoginViaPhone function will handle user authentication using phone number.
func (ah *AuthHandler) LoginViaPhone(context *gin.Context) {
	authData := &LoginViaPhone{}

	if error := Binding.Bind(authData, context); error != nil {
		context.JSON(http.StatusUnprocessableEntity, error)
		return
	}

	debug := context.Query("debug")

	user, error := ah.UserService.CheckUserPhoneNumberExistOrNot(authData.PhoneNo)

	if error != nil {
		errorCode, _ := strconv.Atoi(error.Error.Status)
		context.JSON(errorCode, error)
		return
	}

	dbTransaction := context.MustGet("DB").(*gorm.DB).Begin()

	error = ah.AuthService.AuthenticateUserViaPhoneNumber(dbTransaction, authData.PhoneNo, debug)

	if error != nil {
		dbTransaction.Rollback()
		errorCode, _ := strconv.Atoi(error.Error.Status)
		context.JSON(errorCode, error)
		return
	}

	dbTransaction.Commit()

	result := make(map[string]string)
	result["user_guid"] = user.GUID
	result["message"] = "Successfully sent sms to " + user.PhoneNo

	context.JSON(http.StatusOK, gin.H{"data": result})
}

// LoginViaFacebook function used to login user via facebook
func (ah *AuthHandler) LoginViaFacebook(context *gin.Context) {
	authData := &LoginViaFacebook{}

	if error := Binding.Bind(authData, context); error != nil {
		context.JSON(http.StatusUnprocessableEntity, error)
		return
	}

	user, error := ah.UserService.CheckUserFacebookIDExistOrNot(authData.FacebookID)

	if error != nil {
		errorCode, _ := strconv.Atoi(error.Error.Status)
		context.JSON(errorCode, error)
		return
	}

	dbTransaction := context.MustGet("DB").(*gorm.DB).Begin()

	jwtToken, error := ah.AuthService.AuthenticateUserViaFacebook(dbTransaction, user.GUID, user.PhoneNo, authData.FacebookID, authData.DeviceUUID)

	if error != nil {
		dbTransaction.Rollback()
		errorCode, _ := strconv.Atoi(error.Error.Status)
		context.JSON(errorCode, error)
		return
	}

	dbTransaction.Commit()

	response := make(map[string]interface{})
	response["user"] = user
	response["access_token"] = jwtToken

	context.JSON(http.StatusOK, gin.H{"data": response})
}

// Refresh function used to refresh device token. Example when user close the app and open the app again,
// app must request this endpoint to avoid token expired
func (ah *AuthHandler) Refresh(context *gin.Context) {
	tokenData := context.MustGet("Token").(map[string]string)

	jwtToken, error := ah.AuthService.GenerateJWTTokenForUser(tokenData["user_guid"], tokenData["user_phone_no"], tokenData["device_uuid"])

	if error != nil {
		errorCode, _ := strconv.Atoi(error.Error.Status)
		context.JSON(errorCode, error)
		return
	}

	context.JSON(http.StatusOK, gin.H{"data": jwtToken})
}

// Logout function used to logout user from application.
// System will soft delete device by set deleted_at column to the current date & time.
func (ah *AuthHandler) Logout(context *gin.Context) {
	dbTransaction := context.MustGet("DB").(*gorm.DB).Begin()

	tokenData := context.MustGet("Token").(map[string]string)

	error := ah.AuthService.LogoutUser(dbTransaction, tokenData["device_uuid"], tokenData["user_guid"])

	if error != nil {
		dbTransaction.Rollback()
		errorCode, _ := strconv.Atoi(error.Error.Status)
		context.JSON(errorCode, error)
		return
	}

	dbTransaction.Commit()

	result := make(map[string]string)
	result["message"] = "Successfully logout"

	context.JSON(http.StatusOK, gin.H{"data": result})
}
