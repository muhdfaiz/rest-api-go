package v1

import (
	"bitbucket.org/shoppermate/systems"

	"github.com/fatih/structs"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

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
		GUID:           uuid.NewV4().String(),
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
		uf.DB.Rollback()
		return nil, ErrorMesg.InternalServerError(result.Error, systems.DatabaseError)
	}

	return result.Value.(*User), nil
}

// Update function used to update user detail by certain field.
func (uf *UserFactory) Update(guid string, data UpdateUser) *systems.ErrorData {
	updateData := map[string]string{}
	for key, value := range structs.Map(data) {
		if value != "" {
			updateData[key] = value.(string)
		}
	}

	result := uf.DB.Model(&User{}).Where(&User{GUID: guid}).Updates(updateData)

	if result.Error != nil || result.RowsAffected == 0 {
		uf.DB.Rollback()
		return ErrorMesg.InternalServerError(result.Error, systems.DatabaseError)
	}

	return nil
}

func (uf *UserFactory) Delete(attribute string, value string) *systems.ErrorData {
	result := uf.DB.Where(attribute+" = ?", value).Delete(&User{})

	if result.Error != nil || result.RowsAffected == 0 {
		uf.DB.Rollback()
		return ErrorMesg.InternalServerError(result.Error, systems.DatabaseError)
	}

	return nil
}
