package v1

import (
	"bitbucket.org/cliqers/shoppermate-api/systems"
	"github.com/jinzhu/gorm"
	"github.com/serenize/snaker"
)

type UserRepository struct {
	DB *gorm.DB
}

// Create function will create new user and store in database
func (us *UserRepository) Create(data CreateUser) (*User, *systems.ErrorData) {
	registerBy := "phone_no"

	// Set registerBy equal to facebook if register using Facebook
	if data.FacebookID != "" {
		registerBy = "facebook"
	}

	user := &User{
		GUID:           Helper.GenerateUUID(),
		FacebookID:     data.FacebookID,
		Name:           data.Name,
		Email:          data.Email,
		PhoneNo:        data.PhoneNo,
		ProfilePicture: data.ProfilePicture,
		RegisterBy:     registerBy,
		ReferralCode:   data.ReferralCode,
		Verified:       0,
	}

	// Store new user in database
	result := us.DB.Create(user)

	if result.Error != nil || result.RowsAffected == 0 {
		return nil, Error.InternalServerError(result.Error, systems.DatabaseError)
	}

	return result.Value.(*User), nil
}

// Update function used to update user detail by certain field.
func (us *UserRepository) Update(guid string, data map[string]interface{}) *systems.ErrorData {
	updateData := map[string]interface{}{}

	for key, value := range data {
		if data, ok := value.(string); ok && value.(string) != "" {
			updateData[snaker.CamelToSnake(key)] = data
		}
		if data, ok := value.(int); ok {
			updateData[snaker.CamelToSnake(key)] = data
		}
		if data, ok := value.(float64); ok {
			updateData[snaker.CamelToSnake(key)] = data
		}
	}

	result := us.DB.Model(&User{}).Where(&User{GUID: guid}).Updates(updateData)

	if result.Error != nil {
		return Error.InternalServerError(result.Error, systems.DatabaseError)
	}

	return nil
}

func (us *UserRepository) UpdateUserWallet(userGUID string, amount float64) *systems.ErrorData {
	result := us.DB.Model(&User{}).Where(&User{GUID: userGUID}).Updates(map[string]interface{}{
		"wallet": amount,
	})

	if result.Error != nil {
		return Error.InternalServerError(result.Error, systems.DatabaseError)
	}

	return nil
}

func (us *UserRepository) Delete(attribute string, value string) *systems.ErrorData {
	result := us.DB.Where(attribute+" = ?", value).Delete(&User{})

	if result.Error != nil {
		return Error.InternalServerError(result.Error, systems.DatabaseError)
	}

	return nil
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

// GetByFacebookID function used to retrieve user by facebook id.
// Return user data if found and return empty user if not found
func (ur *UserRepository) GetByFacebookID(facebookID string, relations string) *User {
	user := &User{}

	DB := ur.DB.Model(&User{})

	if relations != "" {
		DB = LoadRelations(DB, relations)
	}

	DB.Where(&User{FacebookID: facebookID}).First(&user)

	return user
}

// SearchByReferralCode function used to search user by referral code
// Return user data if foun and return empty user if not found
func (ur *UserRepository) SearchByReferralCode(referralCode string, relations string) *User {
	user := &User{}

	DB := ur.DB.Model(&User{})

	if relations != "" {
		DB = LoadRelations(DB, relations)
	}

	DB.Where(&User{ReferralCode: referralCode}).First(&user)

	return user
}
