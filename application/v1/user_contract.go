package v1

import (
	"mime/multipart"

	"github.com/jinzhu/gorm"

	"bitbucket.org/cliqers/shoppermate-api/systems"
)

// UserServiceInterface is a contract that defines the method needed for User Service.
type UserServiceInterface interface {
	ViewUser(userGUID, relations string) *User
	CreateUser(dbTransaction *gorm.DB, userData CreateUser, profilePicture multipart.File,
		referralSettings map[string]string, debug string) (*User, *systems.ErrorData)
	UpdateUser(dbTransaction *gorm.DB, userGUID, deviceUUID string, userData UpdateUser,
		profilePicture multipart.File) (*User, *systems.ErrorData)
	CheckUserGUIDExistOrNot(userGUID string) (*User, *systems.ErrorData)
	CheckUserPhoneNumberDuplicate(phoneNo string) *systems.ErrorData
	CheckUserPhoneNumberExistOrNot(phoneNo string) (*User, *systems.ErrorData)
	CheckUserReferralCodeExistOrNot(referralCode string, referralSettings map[string]string) (*User, *systems.ErrorData)
	CheckUserFacebookIDExistOrNot(facebookID string) (*User, *systems.ErrorData)
	CheckUserFacebookIDValidOrNot(facebookID string, debug int) *systems.ErrorData
	UploadUserProfilePicture(file multipart.File) (map[string]string, *systems.ErrorData)
	GenerateReferralCode(name string) string
	DeleteProfilePicture(profilePictureURL *string) *systems.ErrorData
	CalculateAllTimeAmountAndPendingAmount(user *User) *User
}

// UserRepositoryInterface is a contract that defines the method needed for User Repository
type UserRepositoryInterface interface {
	Create(dbTransaction *gorm.DB, data CreateUser) (*User, *systems.ErrorData)
	Update(dbTransaction *gorm.DB, guid string, data map[string]interface{}) *systems.ErrorData
	UpdateUserWallet(dbTransaction *gorm.DB, userGUID string, amount float64) *systems.ErrorData
	Delete(dbTransaction *gorm.DB, attribute string, value string) *systems.ErrorData
	GetByGUID(guid string, relations string) *User
	GetByPhoneNo(phoneNo string, relations string) *User
	GetByFacebookID(facebookID string, relations string) *User
	SearchByReferralCode(referralCode string, relations string) *User
}
