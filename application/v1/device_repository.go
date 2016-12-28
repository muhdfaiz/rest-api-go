package v1

import (
	"bitbucket.org/cliqers/shoppermate-api/systems"
	"github.com/fatih/structs"
	"github.com/jinzhu/gorm"
)

// DeviceRepository will handle all CRUD function for Device resource.
type DeviceRepository struct {
	DB *gorm.DB
}

// Create function used to create new device
// Optional UserGUID because app must register device first when app is loaded
func (dr *DeviceRepository) Create(data CreateDevice) (*Device, *systems.ErrorData) {
	device := &Device{
		GUID:       Helper.GenerateUUID(),
		UserGUID:   data.UserGUID,
		UUID:       data.UUID,
		Os:         data.Os,
		Model:      data.Model,
		PushToken:  data.PushToken,
		AppVersion: data.AppVersion,
	}

	result := dr.DB.Create(device)

	if result.Error != nil || result.RowsAffected == 0 {
		return nil, Error.InternalServerError(result.Error, systems.DatabaseError)
	}

	return result.Value.(*Device), nil
}

// Update function used to update device data. Require device uuid.
func (dr *DeviceRepository) Update(uuid string, data UpdateDevice) *systems.ErrorData {
	updateData := map[string]string{}
	for key, value := range structs.Map(data) {
		if value != "" {
			updateData[key] = value.(string)
		}
	}

	result := dr.DB.Unscoped().Model(&Device{}).Where(&Device{UUID: uuid}).Updates(updateData)

	if result.Error != nil {
		return Error.InternalServerError(result.Error, systems.DatabaseError)
	}

	return nil
}

func (dr *DeviceRepository) SetDeletedAtToNull(deviceGUID string) *systems.ErrorData {
	// Reactivate device by set null to deleted_at column in devices table
	result := dr.DB.Unscoped().Model(&Device{}).Where(&Device{GUID: deviceGUID}).Update("deleted_at", nil)

	if result.Error != nil {
		return Error.InternalServerError(result.Error, systems.DatabaseError)
	}

	return nil
}

// Delete function used to soft delete device.
func (dr *DeviceRepository) Delete(attribute string, value string) *systems.ErrorData {
	result := dr.DB.Where(attribute+" = ?", value).Delete(&Device{})

	if result.Error != nil {
		return Error.InternalServerError(result.Error, systems.DatabaseError)
	}

	return nil
}

// GetByUUID function used to retrieve device by device uuid
func (dr *DeviceRepository) GetByUUID(uuid string) *Device {
	result := dr.DB.Where(&Device{UUID: uuid}).First(&Device{})

	if result.RowsAffected == 0 {
		return &Device{}
	}

	return result.Value.(*Device)
}

// GetByUUIDAndUserGUID function used to retrieve device by device uuid and user guid
func (dr *DeviceRepository) GetByUUIDAndUserGUID(uuid string, userGUID string) *Device {
	result := dr.DB.Where(&Device{UUID: uuid, UserGUID: userGUID}).First(&Device{})

	if result.RowsAffected == 0 {
		return &Device{}
	}

	return result.Value.(*Device)
}

// GetByUUIDUnscoped function used to retrieve device by device uuid and user guid
func (dr *DeviceRepository) GetByUUIDUnscoped(uuid string) *Device {
	result := dr.DB.Unscoped().Where(&Device{UUID: uuid}).First(&Device{})

	if result.RowsAffected == 0 {
		return &Device{}
	}

	return result.Value.(*Device)
}
