package v1

import (
	"mime/multipart"
	"net/url"
	"os"
	"strings"

	"github.com/jinzhu/gorm"

	"bitbucket.org/cliqers/shoppermate-api/services/filesystem"
	"bitbucket.org/cliqers/shoppermate-api/systems"
)

type UserServiceInterface interface {
	CheckUserExistOrNot(userGUID string) *systems.ErrorData
	CheckUserPhoneNumberValidOrNot(phoneNo string) (*User, *systems.ErrorData)
	CheckUserFacebookIDValidOrNot(facebookID string) (*User, *systems.ErrorData)
	UploadProfileImage(file multipart.File) (map[string]string, *systems.ErrorData)
	GiveReferralCashback(referrerGUID string, referentGUID string) (interface{}, *systems.ErrorData)
	GenerateReferralCode(name string) string
	DeleteImage(ImageURL string) *systems.ErrorData
}

type UserService struct {
	DB                 *gorm.DB
	AmazonS3FileSystem *filesystem.AmazonS3Upload
	UserRepository     UserRepositoryInterface
}

// CheckUserExistOrNot function used to check user exist or not in database by checking the user GUID.
func (us *UserService) CheckUserExistOrNot(userGUID string) *systems.ErrorData {
	user := us.UserRepository.GetByGUID(userGUID, "")

	if user.GUID == "" {
		return Error.ResourceNotFoundError("User", "guid", userGUID)
	}

	return nil
}

// CheckUserPhoneNumberValidOrNot function used to check user phone number valid or not.
func (us *UserService) CheckUserPhoneNumberValidOrNot(phoneNo string) (*User, *systems.ErrorData) {
	user := us.UserRepository.GetByPhoneNo(phoneNo, "")

	if user.PhoneNo == "" {
		return nil, Error.ResourceNotFoundError("User", "phone_no", phoneNo)
	}

	return user, nil
}

// CheckUserFacebookIDValidOrNot function used to check user facebook ID valid or not.
func (us *UserService) CheckUserFacebookIDValidOrNot(facebookID string) (*User, *systems.ErrorData) {
	user := us.UserRepository.GetByFacebookID(facebookID, "")

	if user.PhoneNo == "" {
		return nil, Error.ResourceNotFoundError("User", "facebook_id", facebookID)
	}

	return user, nil
}

// UploadProfileImage function used to upload profile image to Amazon S3 if profile_image exist in the request
func (us *UserService) UploadProfileImage(file multipart.File) (map[string]string, *systems.ErrorData) {

	// If profile image file type other than jpg, jpeg, png, gif return error message
	error := FileValidation.ValidateFileType([]string{"jpg", "jpeg", "png", "gif"}, file)
	if error != nil {
		return nil, error
	}

	// If profile image size equal or more than 1MB return error message
	_, error = FileValidation.ValidateFileSize(file, 1000000, "profile_picture")
	if error != nil {
		return nil, error
	}

	localUploadPath := os.Getenv("GOPATH") + Config.Get("app.yaml", "storage_path", "src/bitbucket.org/cliqers/shoppermate-api/storages/")
	amazonS3UploadPath := "/profile_images/"
	uploadedFile, error := us.AmazonS3FileSystem.Upload(file, localUploadPath, amazonS3UploadPath)

	if error != nil {
		return nil, error
	}

	return uploadedFile, nil
}

// GiveReferralCashback function used to give cashback to user that refer by another user during registration
func (us *UserService) GiveReferralCashback(referrerGUID string, referentGUID string) (interface{}, *systems.ErrorData) {
	ReferralCashbackFactory := &ReferralCashbackFactory{DB: us.DB}
	referralCashbackCreated, error := ReferralCashbackFactory.CreateReferralCashbackFactory(referrerGUID, referentGUID)

	if error != nil {
		return nil, error
	}

	return referralCashbackCreated, nil
}

// GenerateReferralCode function used to generate referral code (first 3 letter(UPPERCASE) combine with 5 numeric)
func (us *UserService) GenerateReferralCode(name string) string {

	// Retrieve email name only from full email string
	SplittedName := strings.Split(name, " ")
	firstName := SplittedName[0]

	// Split firstName to letter
	nameLetters := strings.SplitAfter(firstName, "")

	// Grab First 3 Letter from nameLetters
	firstThreeLetterInFirstname := strings.ToUpper(nameLetters[0] + nameLetters[1] + nameLetters[2])
	var referralCode string
	//var counter int = 1

	// Loop until generated referralCode not exist in database
	for {
		randomNumeric := Helper.RandomString("Digit", 5, "", "")

		referralCode = firstThreeLetterInFirstname + randomNumeric

		// Use For Debugging. Don't Enable it on production
		// if counter == 1 {
		// 	referralCode = "MUH84606"
		// }

		// Check referralCode exist in database
		referralCodeExist := us.DB.Where(&User{ReferralCode: referralCode}).First(&User{})
		//fmt.Println(referralCodeExist.RowsAffected)
		//counter++

		if referralCodeExist.RowsAffected == 0 {
			break
		}
	}

	return referralCode
}

// DeleteImage function used to delete profile picture from Amazon S3
func (us *UserService) DeleteImage(ImageURL string) *systems.ErrorData {
	imageURLs := make([]string, 1)

	// Example URI: `https://s3-ap-southeast-1.amazonaws.com/shoppermate-test/profile_images/f83617cd-2b17-3c59-81a5-78c9cfbe7c4f.png`
	url, _ := url.Parse(ImageURL)

	uriSegments := strings.SplitN(url.Path, "/", 3)

	// Retrieve image path after bucket name
	// Example: `profile_images/f83617cd-2b17-3c59-81a5-78c9cfbe7c4f.png`
	imageURLs[0] = uriSegments[2]

	error := us.AmazonS3FileSystem.Delete(imageURLs)

	if error != nil {
		return error
	}

	return nil
}
