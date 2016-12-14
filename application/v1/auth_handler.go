package v1

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	AuthService AuthServiceInterface
}

// LoginViaPhone function will handle user authentication using phone number.
func (ah *AuthHandler) LoginViaPhone(context *gin.Context) {
	authData := &LoginViaPhone{}

	if error := Binding.Bind(authData, context); error != nil {
		context.JSON(http.StatusUnprocessableEntity, error)
		return
	}

	user, error := ah.AuthService.AuthenticateUserViaPhoneNumber(authData.PhoneNo)

	if error != nil {
		errorCode, _ := strconv.Atoi(error.Error.Status)
		context.JSON(errorCode, error)
		return
	}

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

	user, jwtToken, error := ah.AuthService.AuthenticateUserViaFacebook(authData.FacebookID, authData.DeviceUUID)

	if error != nil {
		errorCode, _ := strconv.Atoi(error.Error.Status)
		context.JSON(errorCode, error)
		return
	}

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
	tokenData := context.MustGet("Token").(map[string]string)

	error := ah.AuthService.LogoutUser(tokenData["device_uuid"], tokenData["user_guid"])

	if error != nil {
		errorCode, _ := strconv.Atoi(error.Error.Status)
		context.JSON(errorCode, error)
		return
	}

	result := make(map[string]string)
	result["message"] = "Successfully logout"

	context.JSON(http.StatusOK, gin.H{"data": result})
}
