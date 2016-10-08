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
	UploadProfileImage(file multipart.File) (map[string]string, *systems.ErrorData)
	GiveReferralCashback(referrerGUID string, referentGUID string) (interface{}, *systems.ErrorData)
	GenerateReferralCode(name string) string
	DeleteImage(ImageURL string) *systems.ErrorData
}

type UserService struct {
	DB                 *gorm.DB
	AmazonS3FileSystem *filesystem.AmazonS3Upload
}

// UploadProfileImage function used to upload profile image to Amazon S3 if profile_image exist in the request
func (us *UserService) UploadProfileImage(file multipart.File) (map[string]string, *systems.ErrorData) {

	// If profile image file type other than jpg, jpeg, png, gif return error message
	err := FileValidation.ValidateFileType([]string{"jpg", "jpeg", "png", "gif"}, file)
	if err != nil {
		return nil, err
	}

	// If profile image size equal or more than 1MB return error message
	_, err = FileValidation.ValidateFileSize(file, 1000000, "profile_picture")
	if err != nil {
		return nil, err
	}

	localUploadPath := os.Getenv("GOPATH") + Config.Get("app.yaml", "storage_path", "src/bitbucket.org/cliqers/shoppermate-api/storages/")
	amazonS3UploadPath := "/profile_images/"
	uploadedFile, err := us.AmazonS3FileSystem.Upload(file, localUploadPath, amazonS3UploadPath)

	if err != nil {
		return nil, err
	}

	return uploadedFile, nil
}

// GiveReferralCashback function used to give cashback to user that refer by another user during registration
func (us *UserService) GiveReferralCashback(referrerGUID string, referentGUID string) (interface{}, *systems.ErrorData) {
	ReferralCashbackFactory := &ReferralCashbackFactory{}
	referralCashbackCreated, err := ReferralCashbackFactory.CreateReferralCashbackFactory(referrerGUID, referentGUID)

	if err != nil {
		return nil, err
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

	url, _ := url.Parse(ImageURL)

	uriSegments := strings.SplitN(url.Path, "/", 3)

	imageURLs[0] = uriSegments[2]

	err := us.AmazonS3FileSystem.Delete(imageURLs)

	if err != nil {
		return err
	}

	return nil
}
