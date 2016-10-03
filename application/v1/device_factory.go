package v1

import (
	"bitbucket.org/cliqers/shoppermate-api/systems"
	"github.com/fatih/structs"
	"github.com/jinzhu/gorm"
)

type DeviceFactoryInterface interface {
	Create(data CreateDevice) (*Device, *systems.ErrorData)
	Update(uuid string, data UpdateDevice) *systems.ErrorData
	Delete(attribute string, value string) *systems.ErrorData
}

// DeviceFactory will handle all function to create, update and delete device
type DeviceFactory struct {
	DB *gorm.DB
}

// Create function used to create new device
// Optional UserGUID because app must register device first when app is loaded
func (df *DeviceFactory) Create(data CreateDevice) (*Device, *systems.ErrorData) {
	device := &Device{
		GUID:       Helper.GenerateUUID(),
		UserGUID:   data.UserGUID,
		UUID:       data.UUID,
		Os:         data.Os,
		Model:      data.Model,
		PushToken:  data.PushToken,
		AppVersion: data.AppVersion,
	}

	result := df.DB.Create(device)

	if result.Error != nil || result.RowsAffected == 0 {
		return nil, Error.InternalServerError(result.Error, systems.DatabaseError)
	}

	return result.Value.(*Device), nil
}

// Update function used to update device data
// Require device uuid. Must provide in url
func (df *DeviceFactory) Update(uuid string, data UpdateDevice) *systems.ErrorData {
	updateData := map[string]string{}
	for key, value := range structs.Map(data) {
		if value != "" {
			updateData[key] = value.(string)
		}
	}

	result := df.DB.Model(&Device{}).Where(&Device{UUID: uuid}).Updates(updateData)

	if result.Error != nil {
		return Error.InternalServerError(result.Error, systems.DatabaseError)
	}

	return nil
}

func (df *DeviceFactory) Delete(attribute string, value string) *systems.ErrorData {
	result := df.DB.Where(attribute+" = ?", value).Delete(&Device{})

	if result.Error != nil {
		return Error.InternalServerError(result.Error, systems.DatabaseError)
	}

	return nil
}
