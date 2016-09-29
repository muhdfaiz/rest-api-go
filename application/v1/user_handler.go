package v1

import (
	"fmt"
	"net/http"
	"strconv"

	"bitbucket.org/cliqers/shoppermate-api/services/facebook"
	"bitbucket.org/cliqers/shoppermate-api/systems"

	"github.com/fatih/structs"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// UserHandler will handle all request related to User
type UserHandler struct {
	DB                         *gorm.DB
	UserRepository             UserRepositoryInterface
	UserService                UserServiceInterface
	UserFactory                UserFactoryInterface
	ReferralCashbackRepository ReferralCashbackRepositoryInterface
	SmsService                 SmsServiceInterface
	FacebookService            facebook.FacebookServiceInterface
	DeviceFactory              DeviceFactoryInterface
}

// View function used to view user detail
func (uh *UserHandler) View(c *gin.Context) {
	db := c.MustGet("DB").(*gorm.DB).Begin()

	// Retrieve user guid in url
	userGUID := c.Param("guid")

	// Retrieve user by GUID
	user := uh.UserRepository.GetByGUID(db, userGUID)

	// If user GUID empty return error message
	if user.GUID == "" {
		db.Rollback().Close()
		c.JSON(http.StatusBadRequest, Error.ResourceNotFoundError("User", "guid", userGUID))
		return
	}

	db.Close()
	c.JSON(http.StatusOK, gin.H{"data": user})

}

// Create function will create new user
func (uh *UserHandler) Create(c *gin.Context) {
	db := c.MustGet("DB").(*gorm.DB).Begin()

	userData := CreateUser{}

	// Bind request based on content type and validate request data
	if err := Binding.Bind(&userData, c); err != nil {
		db.Rollback().Close()
		c.JSON(http.StatusUnprocessableEntity, err)
		return
	}

	// Retrieve user by phone_no
	user := uh.UserRepository.GetByPhoneNo(db, userData.PhoneNo)

	// If user phone_no not empty return error message
	if user.PhoneNo != "" {
		db.Rollback().Close()
		c.JSON(http.StatusConflict, Error.DuplicateValueErrors("User", "phone_no", userData.PhoneNo))
		return
	}

	// If facebook_id exist in request data
	if userData.FacebookID != "" {
		// Validate facebook_id valid or not
		fbIDValid := uh.FacebookService.IDIsValid(userData.FacebookID)

		// If facebook_id not valid return error message
		if !fbIDValid {
			db.Rollback().Close()
			mesg := fmt.Sprintf(systems.ErrorFacebookIDNotValid, userData.FacebookID)
			c.JSON(http.StatusBadRequest, Error.GenericError(strconv.Itoa(http.StatusBadRequest),
				systems.FacebookIDNotValid, systems.TitleFacebookIDNotValidError, "facebook_id", mesg))
			return
		}
	}

	user = &User{}
	// If referral_code exist in request data
	if userData.ReferralCode != "" {
		// Search referral code
		user = uh.UserRepository.SearchReferralCode(db, userData.ReferralCode)

		// If referral code not found return error message
		if user.ReferralCode == "" {
			db.Rollback().Close()
			c.JSON(http.StatusBadRequest, Error.GenericError(strconv.Itoa(http.StatusBadRequest),
				systems.ReferralCodeNotExist, systems.TitleReferralCodeNotExist, "referral_code", systems.ErrorReferralCodeNotExist))
			return
		}

		// Count total referral user got
		totalPreviousReferral := uh.ReferralCashbackRepository.Count(db, "referent_guid", user.GUID)

		// If total referral more than 3 return error message
		if totalPreviousReferral > 3 {
			db.Rollback().Close()
			c.JSON(http.StatusBadRequest, Error.GenericError(strconv.Itoa(http.StatusBadRequest),
				systems.ReferralCodeExceedLimit, systems.TitleReferralCodeExceedLimit, "referral_code", systems.ErrorReferralCodeExceedLimit))
			return
		}
	}

	// Retrieve profile_picture in the request
	file, _, _ := c.Request.FormFile("profile_picture")

	profileImage := map[string]string{}

	// If profile_picture exist in the request
	if file != nil {
		err := &systems.ErrorData{}

		// Upload profile picture
		profileImage, err = uh.UserService.UploadProfileImage(file)

		if err != nil {
			db.Rollback().Close()
			errorCode, _ := strconv.Atoi(err.Error.Status)
			c.JSON(errorCode, err)
			return
		}
	}

	// Set profile_picture to the user data
	if profileImage != nil {
		userData.ProfilePicture = profileImage["path"]
	}

	// Store new user in database
	result, err := uh.UserFactory.Create(db, userData)

	if err != nil {
		db.Rollback().Close()
		errorCode, _ := strconv.Atoi(err.Error.Status)
		c.JSON(errorCode, err)
		return
	}

	createdUser := result

	// Set Profile Image to Amazon S3 URL
	if profileImage != nil {
		createdUser.ProfilePicture = profileImage["path"]
	}

	// Send SMS verification code
	_, err = uh.SmsService.SendVerificationCode(db, createdUser.PhoneNo, createdUser.GUID)

	if err != nil {
		db.Rollback().Close()
		errorCode, _ := strconv.Atoi(err.Error.Status)
		c.JSON(errorCode, err)
		return
	}

	// Give cashback to user if referral code validate
	if user.ReferralCode != "" {
		referentUserGUID := user.GUID
		_, err := uh.UserService.GiveReferralCashback(db, createdUser.GUID, referentUserGUID)

		if err != nil {
			db.Rollback().Close()
			errorCode, _ := strconv.Atoi(err.Error.Status)
			c.JSON(errorCode, err)
			return
		}
	}

	// userTransformer := UserTransformer{}
	// userTransformer.TransformCreateData(createdUser)

	db.Commit().Close()
	c.JSON(http.StatusOK, gin.H{"data": createdUser})
}

// Update function used to update user data
func (uh *UserHandler) Update(c *gin.Context) {
	db := c.MustGet("DB").(*gorm.DB).Begin()

	// Retrieve user guid in url
	userGUID := c.Param("guid")

	// Retrieve User Token
	userToken := c.MustGet("Token").(map[string]string)

	if userToken["user_guid"] != userGUID {
		db.Rollback().Close()
		c.JSON(http.StatusUnauthorized, Error.TokenIdentityNotMatchError("Update User"))
		return
	}

	// Retrieve user by guid
	user := uh.UserRepository.GetByGUID(db, userGUID)

	// If user guid empty return error message
	if user.GUID == "" {
		db.Rollback().Close()
		c.JSON(http.StatusBadRequest, Error.ResourceNotFoundError("User", "guid", userGUID))
		return
	}

	userData := UpdateUser{}

	// Bind request based on content type and validate request data
	if err := Binding.Bind(&userData, c); err != nil {
		db.Rollback().Close()
		c.JSON(http.StatusBadRequest, err)
		return
	}

	// Upload profile image if exists
	file, _, _ := c.Request.FormFile("profile_picture")
	profileImage := map[string]string{}

	if file != nil {
		err := &systems.ErrorData{}
		userService := &UserService{}
		profileImage, err = userService.UploadProfileImage(file)

		if err != nil {
			db.Rollback().Close()
			errorCode, _ := strconv.Atoi(err.Error.Status)
			c.JSON(errorCode, err)
			return
		}
	}

	if profileImage != nil {
		userData.ProfilePicture = profileImage["path"]
	}

	// Update User
	userFactory := &UserFactory{}
	err := userFactory.Update(db, userGUID, structs.Map(&userData))

	if err != nil {
		db.Rollback().Close()
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	// Retrieve latest user data
	updatedUser := uh.UserRepository.GetByGUID(db, userGUID)

	// Send SMS verification code
	if user.PhoneNo != updatedUser.PhoneNo && updatedUser.PhoneNo != "" {
		_, err = uh.SmsService.SendVerificationCode(db, updatedUser.PhoneNo, updatedUser.GUID)

		if err != nil {
			db.Rollback().Close()
			errorCode, _ := strconv.Atoi(err.Error.Status)
			c.JSON(errorCode, err)
			return
		}

		// Soft delete device by set current time to deleted_at column
		err := uh.DeviceFactory.Delete(db, "uuid", userToken["device_uuid"])

		if err != nil {
			db.Rollback().Close()
			c.JSON(http.StatusInternalServerError, err)
			return
		}
	}

	db.Commit().Close()
	c.JSON(http.StatusOK, gin.H{"data": updatedUser})

}
