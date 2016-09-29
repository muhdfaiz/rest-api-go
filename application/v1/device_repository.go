package v1

import "github.com/jinzhu/gorm"

type DeviceRepositoryInterface interface {
	GetByUUID(DB *gorm.DB, uuid string) *Device
	GetByUUIDAndUserGUID(DB *gorm.DB, uuid string, userGUID string) *Device
	GetByUUIDUnscoped(DB *gorm.DB, uuid string) *Device
}

type DeviceRepository struct{}

// GetByUUID function used to retrieve device by device uuid
func (dr *DeviceRepository) GetByUUID(DB *gorm.DB, uuid string) *Device {
	result := DB.Where(&Device{UUID: uuid}).First(&Device{})

	if result.RowsAffected == 0 {
		return &Device{}
	}

	return result.Value.(*Device)
}

// GetByUUIDAndUserGUID function used to retrieve device by device uuid and user guid
func (dr *DeviceRepository) GetByUUIDAndUserGUID(DB *gorm.DB, uuid string, userGUID string) *Device {
	result := DB.Where(&Device{UUID: uuid, UserGUID: userGUID}).First(&Device{})

	if result.RowsAffected == 0 {
		return &Device{}
	}

	return result.Value.(*Device)
}

// GetByUUIDUnscoped function used to retrieve device by device uuid and user guid
func (dr *DeviceRepository) GetByUUIDUnscoped(DB *gorm.DB, uuid string) *Device {
	result := DB.Unscoped().Where(&Device{UUID: uuid}).First(&Device{})

	if result.RowsAffected == 0 {
		return &Device{}
	}

	return result.Value.(*Device)
}
