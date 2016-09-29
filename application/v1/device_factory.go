package v1

import (
	"bitbucket.org/cliqers/shoppermate-api/systems"
	"github.com/fatih/structs"
	"github.com/jinzhu/gorm"
)

type DeviceFactoryInterface interface {
	Create(DB *gorm.DB, data CreateDevice) (*Device, *systems.ErrorData)
	Update(DB *gorm.DB, uuid string, data UpdateDevice) *systems.ErrorData
	Delete(DB *gorm.DB, attribute string, value string) *systems.ErrorData
}

// DeviceFactory will handle all function to create, update and delete device
type DeviceFactory struct{}

// Create function used to create new device
// Optional UserGUID because app must register device first when app is loaded
func (df *DeviceFactory) Create(DB *gorm.DB, data CreateDevice) (*Device, *systems.ErrorData) {
	device := &Device{
		GUID:       Helper.GenerateUUID(),
		UserGUID:   data.UserGUID,
		UUID:       data.UUID,
		Os:         data.Os,
		Model:      data.Model,
		PushToken:  data.PushToken,
		AppVersion: data.AppVersion,
	}

	result := DB.Create(device)

	if result.Error != nil || result.RowsAffected == 0 {
		return nil, Error.InternalServerError(result.Error, systems.DatabaseError)
	}

	return result.Value.(*Device), nil
}

// Update function used to update device data
// Require device uuid. Must provide in url
func (df *DeviceFactory) Update(DB *gorm.DB, uuid string, data UpdateDevice) *systems.ErrorData {
	updateData := map[string]string{}
	for key, value := range structs.Map(data) {
		if value != "" {
			updateData[key] = value.(string)
		}
	}

	result := DB.Model(&Device{}).Where(&Device{UUID: uuid}).Updates(updateData)

	if result.Error != nil || result.RowsAffected == 0 {
		return Error.InternalServerError(result.Error, systems.DatabaseError)
	}

	return nil
}

func (df *DeviceFactory) Delete(DB *gorm.DB, attribute string, value string) *systems.ErrorData {
	result := DB.Where(attribute+" = ?", value).Delete(&Device{})

	if result.Error != nil || result.RowsAffected == 0 {
		return Error.InternalServerError(result.Error, systems.DatabaseError)
	}

	return nil
}
