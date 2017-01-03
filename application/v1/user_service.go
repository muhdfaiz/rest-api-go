package v1

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/fatih/structs"
	"github.com/jinzhu/gorm"

	"os"

	"bitbucket.org/cliqers/shoppermate-api/services/facebook"
	"bitbucket.org/cliqers/shoppermate-api/services/filesystem"
	"bitbucket.org/cliqers/shoppermate-api/systems"
)

type UserService struct {
	UserRepository                     UserRepositoryInterface
	TransactionService                 TransactionServiceInterface
	DealCashbackService                DealCashbackServiceInterface
	FacebookService                    facebook.FacebookServiceInterface
	SmsService                         SmsServiceInterface
	DeviceService                      DeviceServiceInterface
	AmazonS3FileSystem                 *filesystem.AmazonS3Upload
	ReferralCashbackTransactionService ReferralCashbackTransactionServiceInterface
}

// ViewUser function used to view user details.
func (us *UserService) ViewUser(userGUID, relations string) *User {
	user := us.UserRepository.GetByGUID(userGUID, relations)

	user = us.CalculateAllTimeAmountAndPendingAmount(user)

	return user
}

// CreateUser function used to create new user and store in database.
func (us *UserService) CreateUser(dbTransaction *gorm.DB, userData CreateUser, profilePicture multipart.File,
	referralSettings map[string]string, debug string) (*User, *systems.ErrorData) {

	error := us.CheckUserPhoneNumberDuplicate(userData.PhoneNo)

	if error != nil {
		return nil, error
	}

	error = us.CheckUserFacebookIDValidOrNot(userData.FacebookID, userData.Debug)

	if error != nil {
		return nil, error
	}

	referentUser, error := us.CheckUserReferralCodeExistOrNot(userData.ReferralCode, referralSettings)

	if error != nil {
		return nil, error
	}

	uploadedProfilePicture, error := us.UploadUserProfilePicture(profilePicture)

	if error != nil {
		return nil, error
	}

	if uploadedProfilePicture != nil {
		userData.ProfilePicture = uploadedProfilePicture["path"]
	}

	userData.ReferralCode = us.GenerateReferralCode(userData.Name)

	newUser, error := us.UserRepository.Create(dbTransaction, userData)

	if error != nil {
		return nil, error
	}

	error = us.CreateReferralCashbackTransaction(dbTransaction, newUser, referentUser, referralSettings)

	if error != nil {
		return nil, error
	}

	if debug != "1" {
		_, error = us.SmsService.SendVerificationCode(dbTransaction, newUser.PhoneNo, newUser.GUID)

		if error != nil {
			return nil, error
		}
	}

	newUser = us.CalculateAllTimeAmountAndPendingAmount(newUser)

	return newUser, nil
}

// UpdateUser function used to update user data and update in database.
func (us *UserService) UpdateUser(dbTransaction *gorm.DB, userGUID, deviceUUID string, userData UpdateUser,
	profilePicture multipart.File) (*User, *systems.ErrorData) {

	user, error := us.CheckUserGUIDExistOrNot(userGUID)

	if error != nil {
		return nil, error
	}

	uploadedProfilePicture, error := us.UploadUserProfilePicture(profilePicture)

	if error != nil {
		return nil, error
	}

	if uploadedProfilePicture != nil {
		userData.ProfilePicture = uploadedProfilePicture["path"]
	}

	error = us.UserRepository.Update(dbTransaction, userGUID, structs.Map(&userData))

	if error != nil {
		return nil, error
	}

	if uploadedProfilePicture != nil && user.ProfilePicture != nil {
		error = us.DeleteProfilePicture(user.ProfilePicture)

		if error != nil {
			return nil, error
		}
	}

	// TODO: Need to think how to handle when user change phone number. For now not priority.
	// error = us.SendSMSAndSetUserStatusToUnverifyWhenPhoneNumberIsNew(user.PhoneNo, userData.PhoneNo, userGUID, deviceUUID)

	// if error != nil {
	// 	return nil, error
	// }

	return user, nil
}

// CheckUserGUIDExistOrNot function used to check user exist or not in database by checking the user GUID.
func (us *UserService) CheckUserGUIDExistOrNot(userGUID string) (*User, *systems.ErrorData) {
	user := us.UserRepository.GetByGUID(userGUID, "")

	if user.GUID == "" {
		return nil, Error.ResourceNotFoundError("User", "guid", userGUID)
	}

	return user, nil
}

// CheckUserPhoneNumberDuplicate function used to check if user phone no already exist in the database.
func (us *UserService) CheckUserPhoneNumberDuplicate(phoneNo string) *systems.ErrorData {

	user := us.UserRepository.GetByPhoneNo(phoneNo, "")

	if user.PhoneNo != "" {
		return Error.DuplicateValueErrors("User", "phone_no", phoneNo)
	}

	return nil
}

// CheckUserPhoneNumberExistOrNot function used to check user phone number exist or not in database.
func (us *UserService) CheckUserPhoneNumberExistOrNot(phoneNo string) (*User, *systems.ErrorData) {
	user := us.UserRepository.GetByPhoneNo(phoneNo, "")

	if user.PhoneNo == "" {
		return nil, Error.ResourceNotFoundError("User", "phone_no", phoneNo)
	}

	return user, nil
}

// CheckUserReferralCodeExistOrNot function used to check user referral code exist or not in database.
func (us *UserService) CheckUserReferralCodeExistOrNot(referralCode string, referralSettings map[string]string) (*User, *systems.ErrorData) {
	if referralCode != "" && referralSettings["referral_active"] == "true" {
		user := us.UserRepository.SearchByReferralCode(referralCode, "")

		if user.ReferralCode == "" {
			return nil, Error.ResourceNotFoundError("User", "referral_code", referralCode)
		}

		return user, nil
	}

	return nil, nil
}

// CheckUserFacebookIDExistOrNot function used to check user facebook ID exist or not in database.
func (us *UserService) CheckUserFacebookIDExistOrNot(facebookID string) (*User, *systems.ErrorData) {
	user := us.UserRepository.GetByFacebookID(facebookID, "")

	if user.PhoneNo == "" {
		return nil, Error.ResourceNotFoundError("User", "facebook_id", facebookID)
	}

	return user, nil
}

// CheckUserFacebookIDValidOrNot function used to check user facebook ID valid or not by querying Facebook API.
func (us *UserService) CheckUserFacebookIDValidOrNot(facebookID string, debug int) *systems.ErrorData {
	if facebookID != "" {
		valid := us.FacebookService.IDIsValid(facebookID, debug)

		if !valid {
			mesg := fmt.Sprintf(systems.ErrorFacebookIDNotValid, facebookID)

			return Error.GenericError(strconv.Itoa(http.StatusBadRequest), systems.FacebookIDNotValid,
				systems.TitleFacebookIDNotValidError, "facebook_id", mesg)
		}
		return nil
	}
	return nil
}

// CreateReferralCashbackTransaction function used to create referral cashback transaction for the person
// who has been refer by another person using account registration.
// Referrer User is a person that give the referral means the person that use other person referral code during
// account registration.
// Referent User is a person that got the referral.
func (us *UserService) CreateReferralCashbackTransaction(dbTransaction *gorm.DB, referrerUser *User, referentUser *User,
	referralSettings map[string]string) *systems.ErrorData {

	if referentUser != nil {
		totalNumberOfUserReferralCashbackTransactions := us.ReferralCashbackTransactionService.CountTotalNumberOfUserReferralCashbackTransaction(referentUser.GUID)

		maxReferralPerUserInInt, _ := strconv.ParseInt(referralSettings["max_referral_per_user"], 10, 64)

		if referralSettings["referral_active"] == "true" && referrerUser.GUID != "" && totalNumberOfUserReferralCashbackTransactions <= maxReferralPerUserInInt {

			pricePerReferralInFloat64, _ := strconv.ParseFloat(referralSettings["price_per_referral"], 64)

			transaction, error := us.TransactionService.CreateTransaction(dbTransaction, referentUser.GUID, "a606113b-fb22-59f3-876f-dd05da7befc7", "669e13c0-eaea-5aef-a25f-6ba54b529e33", pricePerReferralInFloat64)

			if error != nil {
				return error
			}

			_, error = us.ReferralCashbackTransactionService.CreateReferralCashbackTransaction(dbTransaction, referentUser.GUID, referrerUser.GUID, transaction.GUID)

			if error != nil {
				return error
			}

			error = us.UserRepository.UpdateUserWallet(dbTransaction, referentUser.GUID, referentUser.Wallet+pricePerReferralInFloat64)

			if error != nil {
				return error
			}
		}

		return nil
	}

	return nil
}

// SendSMSAndSetUserStatusToUnverifyWhenPhoneNumberIsNew function used to send SMS to user when the phone number
// user enter is new during update user.
func (us *UserService) SendSMSAndSetUserStatusToUnverifyWhenPhoneNumberIsNew(dbTransaction *gorm.DB, oldPhoneNo, newPhoneNo, userGUID,
	deviceUUID string) *systems.ErrorData {

	if oldPhoneNo != newPhoneNo && newPhoneNo != "" {
		_, error := us.SmsService.SendVerificationCode(dbTransaction, newPhoneNo, userGUID)

		if error != nil {
			return error
		}

		// Soft delete device by set current time to deleted_at column
		error = us.DeviceService.DeleteDeviceByUUID(dbTransaction, deviceUUID)

		if error != nil {
			return error
		}

		return nil
	}

	return nil
}

// UploadUserProfilePicture function used to upload profile picture to Amazon S3 if profile_picture exist in the request
// Max profile picture size allowed is 1MB.
func (us *UserService) UploadUserProfilePicture(profilePicture multipart.File) (map[string]string, *systems.ErrorData) {
	if profilePicture != nil {
		error := FileValidation.ValidateFileType([]string{"jpg", "jpeg", "png", "gif"}, profilePicture)

		if error != nil {
			return nil, error
		}

		_, error = FileValidation.ValidateFileSize(profilePicture, 1000000, "profile_picture")

		if error != nil {
			return nil, error
		}

		localUploadPath := os.Getenv("GOPATH") + os.Getenv("STORAGE_PATH")

		amazonS3UploadPath := "/profile_images/"

		uploadedFile, error := us.AmazonS3FileSystem.Upload(profilePicture, localUploadPath, amazonS3UploadPath)

		if error != nil {
			return nil, error
		}

		return uploadedFile, nil
	}

	return nil, nil
}

// GenerateReferralCode function used to generate referral code (first 3 letter(UPPERCASE) combine with 5 numeric)
func (us *UserService) GenerateReferralCode(name string) string {

	SplittedName := strings.Split(name, " ")

	firstName := SplittedName[0]

	// Split firstName to letter
	nameLetters := strings.SplitAfter(firstName, "")

	// Grab First 3 Letter from nameLetters
	firstThreeLetterInFirstname := strings.ToUpper(nameLetters[0] + nameLetters[1] + nameLetters[2])

	var generatedReferralCode string

	for {
		randomNumeric := Helper.RandomString("Digit", 5, "", "")

		generatedReferralCode = firstThreeLetterInFirstname + randomNumeric

		referralCode := us.UserRepository.SearchByReferralCode(generatedReferralCode, "")

		if referralCode.GUID == "" {
			break
		}
	}

	return generatedReferralCode
}

// DeleteProfilePicture function used to delete profile picture from Amazon S3
func (us *UserService) DeleteProfilePicture(profilePictureURL *string) *systems.ErrorData {
	profilePicturesToDelete := make([]string, 1)

	// Example URI: `https://s3-ap-southeast-1.amazonaws.com/shoppermate-test/profile_images/f83617cd-2b17-3c59-81a5-78c9cfbe7c4f.png`
	url, _ := url.Parse(*profilePictureURL)

	uriSegments := strings.SplitN(url.Path, "/", 3)

	// Retrieve image path after bucket name
	// Example: `profile_images/f83617cd-2b17-3c59-81a5-78c9cfbe7c4f.png`
	profilePicturesToDelete[0] = uriSegments[2]

	error := us.AmazonS3FileSystem.Delete(profilePicturesToDelete)

	if error != nil {
		return error
	}

	return nil
}

// CalculateAllTimeAmountAndPendingAmount function used to calculate total amount of all type transactions (deal cashback, referral cashback)
// and calculate pending amount of deal cashbacks.
func (us *UserService) CalculateAllTimeAmountAndPendingAmount(user *User) *User {
	totalAmountOfPendingDealCashbackTransactions := us.TransactionService.SumTotalAmountOfUserPendingTransaction(user.GUID)

	totalAmountOfDealCashbackAddedToList := us.DealCashbackService.SumTotalAmountOfDealAddedTolistByUser(user.GUID)

	totalPendingAmount := totalAmountOfPendingDealCashbackTransactions + totalAmountOfDealCashbackAddedToList

	user.PendingAmount = &totalPendingAmount

	totalCashoutAmount := us.TransactionService.SumTotalAmountOfUserCashoutTransaction(user.GUID)

	user.AllTimeAmount = &totalCashoutAmount

	return user
}
