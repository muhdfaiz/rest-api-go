package v1

import (
	"fmt"
	"net/http"
	"strconv"

	"bitbucket.org/cliqers/shoppermate-api/services/facebook"
	"bitbucket.org/cliqers/shoppermate-api/systems"

	"github.com/fatih/structs"
	"github.com/gin-gonic/gin"
)

// UserHandler will handle all request related to User
type UserHandler struct {
	UserRepository             UserRepositoryInterface
	UserService                UserServiceInterface
	ReferralCashbackRepository ReferralCashbackRepositoryInterface
	SmsService                 SmsServiceInterface
	FacebookService            facebook.FacebookServiceInterface
	DeviceService              DeviceServiceInterface
	TransactionService         TransactionServiceInterface
	DealCashbackService        DealCashbackServiceInterface
}

// View function used to view user detail
func (uh *UserHandler) View(c *gin.Context) {
	// Retrieve user guid in url
	userGUID := c.Param("guid")

	// Retrieve query string for relations
	relations := c.DefaultQuery("include", "")

	// Retrieve user by GUID
	user := uh.UserRepository.GetByGUID(userGUID, relations)

	// If user GUID empty return error message
	if user.GUID == "" {
		c.JSON(http.StatusNotFound, Error.ResourceNotFoundError("User", "guid", userGUID))
		return
	}

	totalAmountOfPendingDealCashbackTransactions := uh.TransactionService.SumTotalAmountOfUserPendingTransaction(userGUID)

	totalAmountOfDealCashbackAddedToList := uh.DealCashbackService.SumTotalAmountOfDealAddedTolistByUser(userGUID)

	totalPendingAmount := totalAmountOfPendingDealCashbackTransactions + totalAmountOfDealCashbackAddedToList

	user.PendingAmount = &totalPendingAmount

	totalCashoutAmount := uh.TransactionService.SumTotalAmountOfUserCashoutTransaction(userGUID)

	user.AllTimeAmount = &totalCashoutAmount

	c.JSON(http.StatusOK, gin.H{"data": user})
}

// Create function will create new user
func (uh *UserHandler) Create(c *gin.Context) {
	userData := CreateUser{}

	// Bind request based on content type and validate request data
	if err := Binding.Bind(&userData, c); err != nil {
		c.JSON(http.StatusUnprocessableEntity, err)
		return
	}

	// Retrieve user by phone_no
	user := uh.UserRepository.GetByPhoneNo(userData.PhoneNo, "")

	// If user phone_no not empty return error message
	if user.PhoneNo != "" {
		c.JSON(http.StatusConflict, Error.DuplicateValueErrors("User", "phone_no", userData.PhoneNo))
		return
	}

	// If facebook_id exist in request data
	if userData.FacebookID != "" {
		fbIDValid := uh.FacebookService.IDIsValid(userData.FacebookID, userData.Debug)

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
		user = uh.UserRepository.SearchByReferralCode(userData.ReferralCode, "")

		// If referral code not found return error message
		if user.ReferralCode == "" {
			c.JSON(http.StatusBadRequest, Error.GenericError(strconv.Itoa(http.StatusBadRequest),
				systems.ReferralCodeNotExist, systems.TitleReferralCodeNotExist, "referral_code", systems.ErrorReferralCodeNotExist))
			return
		}

		// Count total referral user got
		totalPreviousReferral := uh.ReferralCashbackRepository.Count("referent_guid", user.GUID)

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
		profileImage, err = uh.UserService.UploadProfileImage(file)

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
	result, err := uh.UserRepository.Create(userData)

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

	debug := c.Query("debug")

	if debug != "1" {
		// Send SMS verification code
		_, err = uh.SmsService.SendVerificationCode(createdUser.PhoneNo, createdUser.GUID)

		if err != nil {
			errorCode, _ := strconv.Atoi(err.Error.Status)
			c.JSON(errorCode, err)
			return
		}
	}

	// Give cashback to user if referral code validate
	if user.ReferralCode != "" {
		referentUserGUID := user.GUID
		_, err := uh.UserService.GiveReferralCashback(createdUser.GUID, referentUserGUID)

		if err != nil {
			errorCode, _ := strconv.Atoi(err.Error.Status)
			c.JSON(errorCode, err)
			return
		}
	}

	// userTransformer := UserTransformer{}
	// userTransformer.TransformCreateData(createdUser)

	c.JSON(http.StatusOK, gin.H{"data": createdUser})
}

// Update function used to update user data
func (uh *UserHandler) Update(c *gin.Context) {
	// Retrieve user guid in url
	userGUID := c.Param("guid")

	// Retrieve User Token
	userToken := c.MustGet("Token").(map[string]string)

	if userToken["user_guid"] != userGUID {
		c.JSON(http.StatusUnauthorized, Error.TokenIdentityNotMatchError("Update User"))
		return
	}

	// Retrieve user by guid
	user := uh.UserRepository.GetByGUID(userGUID, "")

	// If user guid empty return error message
	if user.GUID == "" {
		c.JSON(http.StatusNotFound, Error.ResourceNotFoundError("User", "guid", userGUID))
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
		profileImage, err = uh.UserService.UploadProfileImage(file)

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
	err := uh.UserRepository.Update(userGUID, structs.Map(&userData))

	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	if user.ProfilePicture != "" {
		// Delete shopping list item image from Amazon S3
		err = uh.UserService.DeleteImage(user.ProfilePicture)

		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}
	}

	// Retrieve latest user data
	updatedUser := uh.UserRepository.GetByGUID(userGUID, "")

	totalAmountOfPendingDealCashbackTransactions := uh.TransactionService.SumTotalAmountOfUserPendingTransaction(userGUID)

	totalAmountOfDealCashbackAddedToList := uh.DealCashbackService.SumTotalAmountOfDealAddedTolistByUser(userGUID)

	totalPendingAmount := totalAmountOfPendingDealCashbackTransactions + totalAmountOfDealCashbackAddedToList

	updatedUser.PendingAmount = &totalPendingAmount

	totalCashoutAmount := uh.TransactionService.SumTotalAmountOfUserCashoutTransaction(userGUID)

	updatedUser.AllTimeAmount = &totalCashoutAmount

	// Send SMS verification code
	if user.PhoneNo != updatedUser.PhoneNo && updatedUser.PhoneNo != "" {
		_, err = uh.SmsService.SendVerificationCode(updatedUser.PhoneNo, updatedUser.GUID)

		if err != nil {
			errorCode, _ := strconv.Atoi(err.Error.Status)
			c.JSON(errorCode, err)
			return
		}

		// Soft delete device by set current time to deleted_at column
		err := uh.DeviceService.DeleteDeviceByUUID(userToken["device_uuid"])

		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"data": updatedUser})
}
