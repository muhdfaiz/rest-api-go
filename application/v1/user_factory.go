package v1

import (
	"bitbucket.org/cliqers/shoppermate-api/systems"
	"github.com/jinzhu/gorm"
	"github.com/serenize/snaker"
)

type UserFactoryInterface interface {
	Create(data CreateUser) (*User, *systems.ErrorData)
	Update(guid string, data map[string]interface{}) *systems.ErrorData
	UpdateUserWallet(userGUID string, amount float64) *systems.ErrorData
	Delete(attribute string, value string) *systems.ErrorData
}

type UserFactory struct {
	DB *gorm.DB
}

// Create function will create new user and store in database
func (uf *UserFactory) Create(data CreateUser) (*User, *systems.ErrorData) {
	registerBy := "phone_no"

	// Set registerBy equal to facebook if register using Facebook
	if data.FacebookID != "" {
		registerBy = "facebook"
	}

	userService := &UserService{DB: uf.DB}
	user := &User{
		GUID:           Helper.GenerateUUID(),
		FacebookID:     data.FacebookID,
		Name:           data.Name,
		Email:          data.Email,
		PhoneNo:        data.PhoneNo,
		ProfilePicture: data.ProfilePicture,
		RegisterBy:     registerBy,
		ReferralCode:   userService.GenerateReferralCode(data.Name),
		Verified:       0,
	}

	// Store new user in database
	result := uf.DB.Create(user)

	if result.Error != nil || result.RowsAffected == 0 {
		return nil, Error.InternalServerError(result.Error, systems.DatabaseError)
	}

	return result.Value.(*User), nil
}

// Update function used to update user detail by certain field.
func (uf *UserFactory) Update(guid string, data map[string]interface{}) *systems.ErrorData {
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

	result := uf.DB.Model(&User{}).Where(&User{GUID: guid}).Updates(updateData)

	if result.Error != nil {
		return Error.InternalServerError(result.Error, systems.DatabaseError)
	}

	return nil
}

func (uf *UserFactory) UpdateUserWallet(userGUID string, amount float64) *systems.ErrorData {
	result := uf.DB.Model(&User{}).Where(&User{GUID: userGUID}).Updates(map[string]interface{}{
		"wallet": amount,
	})

	if result.Error != nil {
		return Error.InternalServerError(result.Error, systems.DatabaseError)
	}

	return nil
}

func (uf *UserFactory) Delete(attribute string, value string) *systems.ErrorData {
	result := uf.DB.Where(attribute+" = ?", value).Delete(&User{})

	if result.Error != nil {
		return Error.InternalServerError(result.Error, systems.DatabaseError)
	}

	return nil
}
