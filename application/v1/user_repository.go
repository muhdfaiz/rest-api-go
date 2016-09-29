package v1

import "github.com/jinzhu/gorm"

type UserRepositoryInterface interface {
	GetByGUID(DB *gorm.DB, guid string) *User
	GetByPhoneNo(DB *gorm.DB, phoneNo string) *User
	GetFacebookID(DB *gorm.DB, facebookID string) *User
	SearchReferralCode(DB *gorm.DB, referralCode string) *User
}

type UserRepository struct{}

// GetByGUID function used to retrieve user by guid.
// Return user data if found and return empty user if not found
func (ur *UserRepository) GetByGUID(DB *gorm.DB, guid string) *User {
	result := DB.Where(&User{GUID: guid}).First(&User{})

	if result.RowsAffected == 0 {
		return &User{}
	}

	return result.Value.(*User)
}

// GetByPhoneNo function used to retrieve user by phone no.
// Return user data if found and return empty user if not found
func (ur *UserRepository) GetByPhoneNo(DB *gorm.DB, phoneNo string) *User {
	result := DB.Where(&User{PhoneNo: phoneNo}).First(&User{})

	if result.RowsAffected == 0 {
		return &User{}
	}

	return result.Value.(*User)
}

// GetFacebookID function used to retrieve user by facebook id.
// Return user data if found and return empty user if not found
func (ur *UserRepository) GetFacebookID(DB *gorm.DB, facebookID string) *User {
	result := DB.Where(&User{FacebookID: facebookID}).First(&User{})

	if result.RowsAffected == 0 {
		return &User{}
	}

	return result.Value.(*User)
}

// SearchReferralCode function used to search user by referral code
// Return user data if foun and return empty user if not found
func (ur *UserRepository) SearchReferralCode(DB *gorm.DB, referralCode string) *User {
	result := DB.Where("referral_code LIKE ?", "%"+referralCode+"%").First(&User{})

	if result.RowsAffected == 0 {
		return &User{}
	}

	return result.Value.(*User)
}
