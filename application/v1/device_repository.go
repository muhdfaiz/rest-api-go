package v1

import "github.com/jinzhu/gorm"

type DeviceRepository struct {
	DB *gorm.DB
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

// GetByUUIDAndUserGUIDUnscoped function used to retrieve device by device uuid and user guid
func (dr *DeviceRepository) GetByUUIDAndUserGUIDUnscoped(uuid string, userGUID string) *Device {
	result := dr.DB.Unscoped().Where(&Device{UUID: uuid, UserGUID: userGUID}).First(&Device{})

	if result.RowsAffected == 0 {
		return &Device{}
	}

	return result.Value.(*Device)
}
