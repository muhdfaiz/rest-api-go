package v11

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// UserHandler will handle all request related to User resource.
type UserHandler struct {
	UserService    UserServiceInterface
	SettingService SettingServiceInterface
}

// View function used to view user detail including user relations.
func (uh *UserHandler) View(context *gin.Context) {
	userGUID := context.Param("guid")

	userToken := context.MustGet("Token").(map[string]string)

	if userToken["user_guid"] != userGUID {
		context.JSON(http.StatusUnauthorized, Error.TokenIdentityNotMatchError("View User"))
		return
	}

	relations := context.DefaultQuery("include", "")

	user := uh.UserService.ViewUser(userGUID, relations)

	context.JSON(http.StatusOK, gin.H{"data": user})
}

// Create function used create new user and store in database.
// If profile picture exists in the request, API will upload to Amazon S3.
func (uh *UserHandler) Create(context *gin.Context) {
	userData := CreateUser{}

	if error := Binding.Bind(&userData, context); error != nil {
		context.JSON(http.StatusUnprocessableEntity, error)
		return
	}

	profilePicture, _, _ := context.Request.FormFile("profile_picture")

	debug := context.Query("debug")
	debugToken := context.Query("debug_token")
	debugFacebook := context.Query("debug_facebook")

	referralActive := uh.SettingService.GetSettingBySlug("referral_active").Value
	pricePerReferral := uh.SettingService.GetSettingBySlug("referral_price").Value
	maxReferralPerUser := uh.SettingService.GetSettingBySlug("max_referral_user").Value

	referralSettings := map[string]string{
		"referral_active":       referralActive,
		"price_per_referral":    pricePerReferral,
		"max_referral_per_user": maxReferralPerUser,
	}

	dbTransaction := context.MustGet("DB").(*gorm.DB).Begin()

	user, error := uh.UserService.CreateUser(dbTransaction, userData, profilePicture, referralSettings, debug, debugFacebook)

	if error != nil {
		dbTransaction.Rollback()
		errorCode, _ := strconv.Atoi(error.Error.Status)
		context.JSON(errorCode, error)
		return
	}

	jwtToken, error := JWT.GenerateToken(user.GUID, userData.PhoneNo, userData.DeviceUUID, debugToken)

	if error != nil {
		dbTransaction.Rollback()
		errorCode, _ := strconv.Atoi(error.Error.Status)
		context.JSON(errorCode, error)
		return
	}

	response := make(map[string]interface{})
	response["user"] = user
	response["access_token"] = jwtToken

	dbTransaction.Commit()

	context.JSON(http.StatusOK, gin.H{"data": response})
}

// Update function used to update user data
func (uh *UserHandler) Update(context *gin.Context) {
	userGUID := context.Param("guid")

	userToken := context.MustGet("Token").(map[string]string)

	if userToken["user_guid"] != userGUID {
		context.JSON(http.StatusUnauthorized, Error.TokenIdentityNotMatchError("Update User"))
		return
	}

	userData := UpdateUser{}

	if error := Binding.Bind(&userData, context); error != nil {
		context.JSON(http.StatusBadRequest, error)
		return
	}

	profilePicture, _, _ := context.Request.FormFile("profile_picture")

	dbTransaction := context.MustGet("DB").(*gorm.DB).Begin()

	updatedUser, error := uh.UserService.UpdateUser(dbTransaction, userGUID, userToken["device_uuid"], userData, profilePicture)

	if error != nil {
		dbTransaction.Rollback()
		errorCode, _ := strconv.Atoi(error.Error.Status)
		context.JSON(errorCode, error)
		return
	}

	dbTransaction.Commit()

	updatedUser = uh.UserService.ViewUser(userGUID, "")

	updatedUser = uh.UserService.CalculateAllTimeAmountAndPendingAmount(updatedUser)

	context.JSON(http.StatusOK, gin.H{"data": updatedUser})
}
