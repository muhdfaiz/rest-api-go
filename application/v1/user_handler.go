package v1

import (
	"fmt"
	"net/http"
	"strconv"

	"bitbucket.org/shoppermate-api/systems"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// UserHandler will handle all request related to User
type UserHandler struct{}

// View function used to view user detail
func (uh *UserHandler) View(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "application/json; charset=utf-8")

	db := c.MustGet("DB").(*gorm.DB)
	tx := db.Begin()

	// Retrieve user guid in url
	userGUID := c.Param("guid")

	// Retrieve user by GUID
	userRepository := &UserRepository{DB: tx}
	user := userRepository.GetByGUID(userGUID)

	// If user GUID empty return error message
	if user.GUID == "" {
		c.JSON(http.StatusBadRequest, Error.ResourceNotFoundError("User", "guid", userGUID))
		return
	}

	tx.Commit()
	c.JSON(http.StatusOK, gin.H{"data": user})

}

// Create function will create new user
func (uh *UserHandler) Create(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "application/json; charset=utf-8")

	db := c.MustGet("DB").(*gorm.DB)
	tx := db.Begin()

	userData := CreateUser{}

	// Bind request based on content type and validate request data
	if err := Binding.Bind(&userData, c); err != nil {
		c.JSON(http.StatusUnprocessableEntity, err)
		return
	}

	// Retrieve user by phone_no
	userRepository := &UserRepository{DB: tx}
	user := userRepository.GetByPhoneNo(userData.PhoneNo)

	// If user phone_no not empty return error message
	if user.PhoneNo != "" {
		c.JSON(http.StatusConflict, Error.DuplicateValueErrors("User", "phone_no", userData.PhoneNo))
		return
	}

	// If facebook_id exist in request data
	if userData.FacebookID != "" {
		// Validate facebook_id valid or not
		fbIDValid := FacebookService.IDIsValid(userData.FacebookID)

		// If facebook_id not valid return error message
		if !fbIDValid {
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
		user = userRepository.SearchReferralCode(userData.ReferralCode)

		// If referral code not found return error message
		if user.ReferralCode == "" {
			c.JSON(http.StatusBadRequest, Error.GenericError(strconv.Itoa(http.StatusBadRequest),
				systems.ReferralCodeNotExist, systems.TitleReferralCodeNotExist, "referral_code", systems.ErrorReferralCodeNotExist))
			return
		}

		// Count total referral user got
		referralCashbackRepository := &ReferralCashbackRepository{DB: tx}
		totalPreviousReferral := referralCashbackRepository.Count("referent_guid", user.GUID)

		// If total referral more than 3 return error message
		if totalPreviousReferral > 3 {
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
		userService := &UserService{DB: tx}
		profileImage, err = userService.UploadProfileImage(file)

		if err != nil {
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
	userFactory := &UserFactory{DB: tx}
	result, err := userFactory.Create(userData)

	if err != nil {
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
	smsService := &SmsService{DB: tx}
	_, err = smsService.SendVerificationCode(createdUser.PhoneNo, createdUser.GUID)

	if err != nil {
		errorCode, _ := strconv.Atoi(err.Error.Status)
		c.JSON(errorCode, err)
		return
	}

	// Give cashback to user if referral code validate
	if user.ReferralCode != "" {
		referentUserGUID := user.GUID
		userService := UserService{DB: tx}
		_, err := userService.GiveReferralCashback(createdUser.GUID, referentUserGUID)

		if err != nil {
			errorCode, _ := strconv.Atoi(err.Error.Status)
			c.JSON(errorCode, err)
			return
		}
	}

	// userTransformer := UserTransformer{}
	// userTransformer.TransformCreateData(createdUser)

	tx.Commit()
	c.JSON(http.StatusOK, gin.H{"data": createdUser})
}

// Update function used to update user data
func (uh *UserHandler) Update(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "application/json; charset=utf-8")

	db := c.MustGet("DB").(*gorm.DB)
	tx := db.Begin()

	// Retrieve user guid in url
	userGUID := c.Param("guid")

	// Retrieve User Token
	userToken := c.MustGet("Token").(map[string]string)

	if userToken["user_guid"] != userGUID {
		c.JSON(http.StatusBadRequest, Error.TokenIdentityNotMatchError("Update User"))
		return
	}

	// Retrieve user by guid
	userRepository := &UserRepository{DB: tx}
	user := userRepository.GetByGUID(userGUID)

	// If user guid empty return error message
	if user.GUID == "" {
		c.JSON(http.StatusBadRequest, Error.ResourceNotFoundError("User", "guid", userGUID))
		return
	}

	userData := UpdateUser{}

	// Bind request based on content type and validate request data
	if err := Binding.Bind(&userData, c); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	// Upload profile image if exists
	file, _, _ := c.Request.FormFile("profile_picture")
	profileImage := map[string]string{}

	if file != nil {
		err := &systems.ErrorData{}
		userService := &UserService{DB: tx}
		profileImage, err = userService.UploadProfileImage(file)

		if err != nil {
			errorCode, _ := strconv.Atoi(err.Error.Status)
			c.JSON(errorCode, err)
			return
		}
	}

	if profileImage != nil {
		userData.ProfilePicture = profileImage["path"]
	}

	// Update User
	userFactory := &UserFactory{DB: tx}
	err := userFactory.Update(userGUID, userData)

	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	// Retrieve latest user data
	updatedUser := userRepository.GetByGUID(userGUID)

	// Send SMS verification code and soft delete device if user change the phone no
	if user.PhoneNo != updatedUser.PhoneNo {
		smsService := &SmsService{DB: tx}
		_, err = smsService.SendVerificationCode(updatedUser.PhoneNo, updatedUser.GUID)

		if err != nil {
			errorCode, _ := strconv.Atoi(err.Error.Status)
			c.JSON(errorCode, err)
			return
		}
	}

	tx.Commit()
	c.JSON(http.StatusOK, gin.H{"data": updatedUser})

}
