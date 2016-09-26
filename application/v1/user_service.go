package v1

import (
	"mime/multipart"
	"os"
	"strings"

	"bitbucket.org/cliqers/shoppermate-api/services/filesystem"
	"bitbucket.org/cliqers/shoppermate-api/systems"
	"github.com/jinzhu/gorm"
)

type UserServiceInterface interface {
	UploadProfileImage(file multipart.File) (map[string]string, *systems.ErrorData)
	GiveReferralCashback(referrerGUID string, referentGUID string) (interface{}, *systems.ErrorData)
	GenerateReferralCode(name string) string
}

type UserService struct {
	DB           *gorm.DB
	ReferralCode string
}

type AmazonS3UploadConfig struct{}

func (auc *AmazonS3UploadConfig) SetAmazonS3UploadPath() string {
	return "/profile_image/"
}

func (auc *AmazonS3UploadConfig) SetLocalUploadPath() string {
	return os.Getenv("GOPATH") + "src/bitbucket.org/cliqers/shoppermate-api/storages/"
}

func (auc *AmazonS3UploadConfig) SetBucketName() string {
	return Config.Get("app.yaml", "aws_bucket_name", "shoppermate-api")
}

// UploadProfileImage function used to upload profile image to Amazon S3 if profile_image exist in the request
func (us *UserService) UploadProfileImage(file multipart.File) (map[string]string, *systems.ErrorData) {
	// Validate file type is image
	_, err := FileValidation.ValidateFileType([]string{"jpg", "jpeg", "png", "gif"}, file)
	if err != nil {
		return nil, err
	}

	// Validate file size
	_, err = FileValidation.ValidateFileSize(file, 1000000, "profile_picture")
	if err != nil {
		return nil, err
	}

	fileSystem := FileSystem.Driver("amazonS3").(*filesystem.AmazonS3ServiceUpload)
	fileSystem.AccessKey = Config.Get("app.yaml", "aws_access_key_id", "")
	fileSystem.SecretKey = Config.Get("app.yaml", "aws_secret_access_key", "")
	fileSystem.Region = Config.Get("app.yaml", "aws_region_name", "")

	amazonS3UploadConfig := &AmazonS3UploadConfig{}
	uploadedFile, err := fileSystem.Upload(amazonS3UploadConfig, file)
	if err != nil {
		return nil, err
	}

	return uploadedFile, nil
}

// GiveReferralCashback function used to give cashback to user that refer by another user during registration
func (us *UserService) GiveReferralCashback(referrerGUID string, referentGUID string) (interface{}, *systems.ErrorData) {
	ReferralCashbackFactory := &ReferralCashbackFactory{DB: us.DB}
	referralCashbackCreated, err := ReferralCashbackFactory.CreateReferralCashbackFactory(referrerGUID, referentGUID)

	if err != nil {
		us.DB.Rollback()
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
