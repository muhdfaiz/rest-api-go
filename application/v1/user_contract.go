package v1

import (
	"mime/multipart"

	"bitbucket.org/cliqers/shoppermate-api/systems"
)

// UserServiceInterface is a contract that defines the method needed for User Service.
type UserServiceInterface interface {
	CheckUserExistOrNot(userGUID string) *systems.ErrorData
	CheckUserPhoneNumberValidOrNot(phoneNo string) (*User, *systems.ErrorData)
	CheckUserFacebookIDValidOrNot(facebookID string) (*User, *systems.ErrorData)
	UploadProfileImage(file multipart.File) (map[string]string, *systems.ErrorData)
	GiveReferralCashback(referrerGUID string, referentGUID string) (interface{}, *systems.ErrorData)
	GenerateReferralCode(name string) string
	DeleteImage(ImageURL string) *systems.ErrorData
}

// UserRepositoryInterface is a contract that defines the method needed for User Repository
type UserRepositoryInterface interface {
	Create(data CreateUser) (*User, *systems.ErrorData)
	Update(guid string, data map[string]interface{}) *systems.ErrorData
	UpdateUserWallet(userGUID string, amount float64) *systems.ErrorData
	Delete(attribute string, value string) *systems.ErrorData
	GetByGUID(guid string, relations string) *User
	GetByPhoneNo(phoneNo string, relations string) *User
	GetByFacebookID(facebookID string, relations string) *User
	SearchByReferralCode(referralCode string, relations string) *User
}
