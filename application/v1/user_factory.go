package v1

import (
	"bitbucket.org/cliqers/shoppermate-api/systems"
	"github.com/jinzhu/gorm"
)

type UserFactoryInterface interface {
	Create(DB *gorm.DB, data CreateUser) (*User, *systems.ErrorData)
	Update(DB *gorm.DB, guid string, data map[string]interface{}) *systems.ErrorData
	Delete(DB *gorm.DB, attribute string, value string) *systems.ErrorData
}

type UserFactory struct{}

// Create function will create new user and store in database
func (uf *UserFactory) Create(DB *gorm.DB, data CreateUser) (*User, *systems.ErrorData) {
	registerBy := "phone_no"

	// Set registerBy equal to facebook if register using Facebook
	if data.FacebookID != "" {
		registerBy = "facebook"
	}

	userService := &UserService{}
	user := &User{
		GUID:           Helper.GenerateUUID(),
		FacebookID:     data.FacebookID,
		Name:           data.Name,
		Email:          data.Email,
		PhoneNo:        data.PhoneNo,
		ProfilePicture: data.ProfilePicture,
		RegisterBy:     registerBy,
		ReferralCode:   userService.GenerateReferralCode(DB, data.Name),
		Verified:       0,
	}

	// Store new user in database
	result := DB.Create(user)

	if result.Error != nil || result.RowsAffected == 0 {
		return nil, Error.InternalServerError(result.Error, systems.DatabaseError)
	}

	return result.Value.(*User), nil
}

// Update function used to update user detail by certain field.
func (uf *UserFactory) Update(DB *gorm.DB, guid string, data map[string]interface{}) *systems.ErrorData {
	updateData := map[string]interface{}{}
	for key, value := range data {
		if data, ok := value.(string); ok && value.(string) != "" {
			updateData[key] = data
		}
		if data, ok := value.(int); ok && value.(int) != 0 {
			updateData[key] = data
		}
	}

	result := DB.Model(&User{}).Where(&User{GUID: guid}).Updates(updateData)

	if result.Error != nil {
		return Error.InternalServerError(result.Error, systems.DatabaseError)
	}

	return nil
}

func (uf *UserFactory) Delete(DB *gorm.DB, attribute string, value string) *systems.ErrorData {
	result := DB.Where(attribute+" = ?", value).Delete(&User{})

	if result.Error != nil {
		return Error.InternalServerError(result.Error, systems.DatabaseError)
	}

	return nil
}
