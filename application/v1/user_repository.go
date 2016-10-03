package v1

import "github.com/jinzhu/gorm"

type UserRepositoryInterface interface {
	GetByGUID(guid string) *User
	GetByPhoneNo(phoneNo string) *User
	GetFacebookID(facebookID string) *User
	SearchReferralCode(referralCode string) *User
}

type UserRepository struct {
	DB *gorm.DB
}

// GetByGUID function used to retrieve user by guid.
// Return user data if found and return empty user if not found
func (ur *UserRepository) GetByGUID(guid string) *User {
	result := ur.DB.Where(&User{GUID: guid}).First(&User{})

	if result.RowsAffected == 0 {
		return &User{}
	}

	return result.Value.(*User)
}

// GetByPhoneNo function used to retrieve user by phone no.
// Return user data if found and return empty user if not found
func (ur *UserRepository) GetByPhoneNo(phoneNo string) *User {
	result := ur.DB.Where(&User{PhoneNo: phoneNo}).First(&User{})

	if result.RowsAffected == 0 {
		return &User{}
	}

	return result.Value.(*User)
}

// GetFacebookID function used to retrieve user by facebook id.
// Return user data if found and return empty user if not found
func (ur *UserRepository) GetFacebookID(facebookID string) *User {
	result := ur.DB.Where(&User{FacebookID: facebookID}).First(&User{})

	if result.RowsAffected == 0 {
		return &User{}
	}

	return result.Value.(*User)
}

// SearchReferralCode function used to search user by referral code
// Return user data if foun and return empty user if not found
func (ur *UserRepository) SearchReferralCode(referralCode string) *User {
	result := ur.DB.Where("referral_code LIKE ?", "%"+referralCode+"%").First(&User{})

	if result.RowsAffected == 0 {
		return &User{}
	}

	return result.Value.(*User)
}
