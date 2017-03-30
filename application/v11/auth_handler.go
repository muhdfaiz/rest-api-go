package v11

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// AuthHandler will handle all request related to Auth resource.
type AuthHandler struct {
	AuthService AuthServiceInterface
	UserService UserServiceInterface
}

// LoginViaPhone function will handle user authentication using phone number.
// First, it will bind request body to struct and validate the request data.
// Then, it will check if phone number exist or not in database. If not exist return an error.
// Then, it will start database transaction.
// Then, it will continue authenticate user through user service.
// It will return an error encountered from user service and rollback database transaction.
// Lastly, if user service not return any error, it will commit database transaction and
// return the response.
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

// LoginViaFacebook function will handle user authentication using facebook.
// First, it will bind request body to struct and validate the request data.
// Then, it will check if facebook id exist or not through user service. If not exist return an error.
// Then, it will start database transaction.
// Then, it will continue authenticate user through user service.
// It will return an error encountered from user service and rollback database transaction.
// Lastly, if user service not return any error, it will commit database transaction and
// return the response contain user info and access token.
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

// Refresh function used to refresh access token. Useful to avoid access token expired issue.
// Example when user close the app and open the app again, app must request this endpoint to avoid access token expired.
// First, it will retrieve access token in context. Incoming requests to a server should create a Context.
// See auth middleware (middlewre/auth.go) how API store the access token in context when the request coming.
// Then, it will generate new token through auth service. It will return an error encountered from auth service.
// Lastly, it will return response contain new access token if auth service not return an error.
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
// It will soft delete device by setting current date & time as a value for deleted_at field.
// First, it will retrieve access token and database connection in context. Incoming requests to a server should create a Context.
// Then, it will logout user through auth service based on device uuid and user guid in token payload.
// It will return an error encountered from auth service and rollback database transaction.
// Lastly, if auth service not return any error, it will commit database transaction and
// return the response.
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
