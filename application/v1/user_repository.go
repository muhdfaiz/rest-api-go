package v1

import "github.com/jinzhu/gorm"

type UserRepositoryInterface interface {
	GetByGUID(guid string, relations string) *User
	GetByPhoneNo(phoneNo string, relations string) *User
	GetFacebookID(facebookID string, relations string) *User
	SearchReferralCode(referralCode string, relations string) *User
}

type UserRepository struct {
	DB *gorm.DB
}

// GetByGUID function used to retrieve user by guid.
// Return user data if found and return empty user if not found
func (ur *UserRepository) GetByGUID(guid string, relations string) *User {
	user := &User{}

	DB := ur.DB.Model(&User{})

	if relations != "" {
		DB = LoadRelations(DB, relations)
	}

	DB.Where(&User{GUID: guid}).First(&user)

	return user
}

// GetByPhoneNo function used to retrieve user by phone no.
// Return user data if found and return empty user if not found
func (ur *UserRepository) GetByPhoneNo(phoneNo string, relations string) *User {
	user := &User{}

	DB := ur.DB.Model(&User{})

	if relations != "" {
		DB = LoadRelations(DB, relations)
	}

	DB.Where(&User{PhoneNo: phoneNo}).First(&user)

	return user
}

// GetFacebookID function used to retrieve user by facebook id.
// Return user data if found and return empty user if not found
func (ur *UserRepository) GetFacebookID(facebookID string, relations string) *User {
	user := &User{}

	DB := ur.DB.Model(&User{})

	if relations != "" {
		DB = LoadRelations(DB, relations)
	}

	DB.Where(&User{FacebookID: facebookID}).First(&user)

	return user
}

// SearchReferralCode function used to search user by referral code
// Return user data if foun and return empty user if not found
func (ur *UserRepository) SearchReferralCode(referralCode string, relations string) *User {
	user := &User{}

	DB := ur.DB.Model(&User{})

	if relations != "" {
		DB = LoadRelations(DB, relations)
	}

	DB.Where("referral_code LIKE ?", "%"+referralCode+"%").First(&User{})

	return user
}
