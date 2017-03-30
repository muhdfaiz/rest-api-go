package v11

import (
	"bitbucket.org/cliqers/shoppermate-api/systems"
	"github.com/jinzhu/gorm"
	"github.com/serenize/snaker"
)

// UserRepository will handle all CRUD task related to User resource.
type UserRepository struct {
	BaseRepository
	DB *gorm.DB
}

// Create function used to create new user and store in database.
// Use database transaction to create new user. Don't forgot to commit the transaction after used this function.
// Return new user and error if encountered.
func (ur *UserRepository) Create(dbTransansaction *gorm.DB, data CreateUser) (*User, *systems.ErrorData) {
	registerBy := "phone_no"

	if data.FacebookID != "" {
		registerBy = "facebook"
	}

	user := &User{
		GUID:           Helper.GenerateUUID(),
		Name:           data.Name,
		Email:          data.Email,
		FacebookID:     &data.FacebookID,
		PhoneNo:        data.PhoneNo,
		ProfilePicture: &data.ProfilePicture,
		RegisterBy:     registerBy,
		ReferralCode:   data.ReferralCode,
	}

	// Set facebook ID to nil if empty.
	if data.FacebookID == "" {
		user.FacebookID = nil
	}

	// Set profile picture to nil if empty.
	if data.ProfilePicture == "" {
		user.ProfilePicture = nil
	}

	result := dbTransansaction.Create(user)

	if result.Error != nil || result.RowsAffected == 0 {
		return nil, Error.InternalServerError(result.Error, systems.DatabaseError)
	}

	return result.Value.(*User), nil
}

// Update function used to update existing user information using user GUID.
// It's not update all fields available in user table but only update fields that exist in data parameter.
// Use database transaction to update existing user info. Don't forget to commit the transaction after used this function.
// Return updated user and error if encountered.
func (ur *UserRepository) Update(dbTransaction *gorm.DB, guid string, data map[string]interface{}) *systems.ErrorData {
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

	result := dbTransaction.Model(&User{}).Where(&User{GUID: guid}).Updates(updateData)

	if result.Error != nil {
		return Error.InternalServerError(result.Error, systems.DatabaseError)
	}

	return nil
}

// UpdateUserWallet function used to update wallet amount for specific user using user GUID.
// Return nil or internal server error.
func (ur *UserRepository) UpdateUserWallet(dbTransaction *gorm.DB, userGUID string, amount float64) *systems.ErrorData {
	result := dbTransaction.Model(&User{}).Where(&User{GUID: userGUID}).Updates(map[string]interface{}{
		"wallet": amount,
	})

	if result.Error != nil {
		return Error.InternalServerError(result.Error, systems.DatabaseError)
	}

	return nil
}

// Delete function used to soft delete user in database.
// Soft delete means, API will set current date and time to `deleted_at` field in user table.
// Return nil or internal server error.
func (ur *UserRepository) Delete(dbTransaction *gorm.DB, attribute string, value string) *systems.ErrorData {
	result := dbTransaction.Where(attribute+" = ?", value).Delete(&User{})

	if result.Error != nil {
		return Error.InternalServerError(result.Error, systems.DatabaseError)
	}

	return nil
}

// GetByGUID function used to retrieve first user filter by guid.
// Return user data if found and return empty user if not found.
func (ur *UserRepository) GetByGUID(guid string, relations string) *User {
	user := &User{}

	DB := ur.DB.Model(&User{})

	if relations != "" {
		DB = ur.LoadRelations(DB, relations)
	}

	DB.Where(&User{GUID: guid}).First(&user)

	return user
}

// GetByPhoneNo function used to retrieve first user filter by phone no.
// Return user data if found and return empty user if not found
func (ur *UserRepository) GetByPhoneNo(phoneNo string, relations string) *User {
	user := &User{}

	DB := ur.DB.Model(&User{})

	if relations != "" {
		DB = ur.LoadRelations(DB, relations)
	}

	DB.Where(&User{PhoneNo: phoneNo}).First(&user)

	return user
}

// GetByFacebookID function used to retrieve first user filter by facebook id.
// Return user data if found and return empty user if not found
func (ur *UserRepository) GetByFacebookID(facebookID string, relations string) *User {
	user := &User{}

	DB := ur.DB.Model(&User{})

	if relations != "" {
		DB = ur.LoadRelations(DB, relations)
	}

	DB.Where(&User{FacebookID: &facebookID}).First(&user)

	return user
}

// SearchByReferralCode function used to retrieve first user filter by referral code
// Return user data if found and return empty user if not found
func (ur *UserRepository) SearchByReferralCode(referralCode string, relations string) *User {
	user := &User{}

	DB := ur.DB.Model(&User{})

	if relations != "" {
		DB = ur.LoadRelations(DB, relations)
	}

	DB.Where(&User{ReferralCode: referralCode}).First(&user)

	return user
}
