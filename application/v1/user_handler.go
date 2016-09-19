package v1

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"

	"bitbucket.org/shoppermate-api/systems"
)

// UserHandler will handle all request related to User
type UserHandler struct{}

func (uh *UserHandler) View(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	db := c.MustGet("DB").(*gorm.DB)
	tx := db.Begin()

	// Retrieve user guid in url
	userGUID := c.Param("guid")

	// Check User GUID valid
	userRepository := &UserRepository{DB: tx}
	user := userRepository.GetByGUID(userGUID)

	if user.GUID == "" {
		c.JSON(http.StatusBadRequest, ErrorMesg.GenericError(strconv.Itoa(http.StatusBadRequest), systems.ResourceNotFound,
			fmt.Sprintf(systems.TitleResourceNotFoundError, "User"), "",
			fmt.Sprintf(systems.ErrorResourceNotFound, "User", "guid", userGUID)))
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
		statusCode, _ := strconv.Atoi(err.Error.Status)
		c.JSON(statusCode, err)
		return
	}

	// Check User Exist
	userRepository := &UserRepository{DB: tx}
	user := userRepository.GetByPhoneNo(userData.PhoneNo)

	if user.PhoneNo != "" {
		c.JSON(http.StatusConflict, ErrorMesg.DuplicateValueErrors("User", "phone_no", userData.PhoneNo))
		return
	}

	// Check if user register using facebook
	// If true then validate facebook_id field
	if userData.FacebookID != "" {
		fbIDValid := FacebookService.IDIsValid(userData.FacebookID)

		if !fbIDValid {
			mesg := fmt.Sprintf(systems.ErrorFacebookIDNotValid, userData.FacebookID)
			c.JSON(http.StatusBadRequest, ErrorMesg.GenericError(strconv.Itoa(http.StatusBadRequest),
				systems.FacebookIDNotValid, systems.TitleFacebookIDNotValidError, "facebook_id", mesg))
			return
		}
	}

	// If user registration data contain referral code, check if referral code exists.
	// If exist, count how many referral for referent user.
	// If more than 3 return error message
	user = &User{}
	if userData.ReferralCode != "" {
		user = userRepository.SearchReferralCode(userData.ReferralCode)

		if user.ReferralCode == "" {
			c.JSON(http.StatusBadRequest, ErrorMesg.GenericError(strconv.Itoa(http.StatusBadRequest),
				systems.ReferralCodeNotExist, systems.TitleReferralCodeNotExist, "referral_code", systems.ErrorReferralCodeNotExist))
			return
		}

		referralCashbackRepository := &ReferralCashbackRepository{DB: tx}
		totalPreviousReferral := referralCashbackRepository.Count("referent_guid", user.GUID)

		if totalPreviousReferral > 3 {
			c.JSON(http.StatusBadRequest, ErrorMesg.GenericError(strconv.Itoa(http.StatusBadRequest),
				systems.ReferralCodeExceedLimit, systems.TitleReferralCodeExceedLimit, "referral_code", systems.ErrorReferralCodeExceedLimit))
			return
		}
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

// Create function will create new user
func (uh *UserHandler) Update(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	db := c.MustGet("DB").(*gorm.DB)
	tx := db.Begin()

	// Retrieve user guid in url
	userGUID := c.Param("guid")

	// Check User GUID valid
	userRepository := &UserRepository{DB: tx}
	user := userRepository.GetByGUID(userGUID)

	if user.GUID == "" {
		c.JSON(http.StatusBadRequest, ErrorMesg.GenericError(strconv.Itoa(http.StatusBadRequest), systems.ResourceNotFound,
			fmt.Sprintf(systems.TitleResourceNotFoundError, "User"), "message",
			fmt.Sprintf(systems.ErrorResourceNotFound, "User", "guid", userGUID)))
		return
	}

	userData := UpdateUser{}

	// Bind request based on content type and validate request data
	if err := Binding.Bind(&userData, c); err != nil {
		statusCode, _ := strconv.Atoi(err.Error.Status)
		c.JSON(statusCode, err)
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
	fmt.Println(userData)
	// Update User
	userFactory := &UserFactory{DB: tx}
	err := userFactory.Update(userGUID, userData)

	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	// Retrieve user latest update
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
