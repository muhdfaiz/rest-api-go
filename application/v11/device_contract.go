package v11

import (
	"bitbucket.org/cliqers/shoppermate-api/systems"
	"github.com/jinzhu/gorm"
)

// DeviceServiceInterface is a contract that defines the method needed for Device Service.
type DeviceServiceInterface interface {
	CheckDuplicateDevice(deviceUUID string) *systems.ErrorData
	CheckDeviceExistOrNot(deviceUUID string) (*Device, *systems.ErrorData)
	CreateDevice(dbTransaction *gorm.DB, deviceData CreateDevice) (*Device, *systems.ErrorData)
	UpdateDevice(dbTransaction *gorm.DB, deviceUUID string, deviceData UpdateDevice) (*Device, *systems.ErrorData)
	UpdateByDeviceUUID(dbTransaction *gorm.DB, deviceUUID string, deviceData UpdateDevice) (*Device, *systems.ErrorData)
	ReactivateDevice(dbTransaction *gorm.DB, deviceGUID string) *systems.ErrorData
	DeleteDeviceByUUID(dbTransaction *gorm.DB, deviceUUID string) *systems.ErrorData
	ViewDeviceByUUID(deviceUUID string) *Device
	ViewDeviceByUUIDandUserGUID(deviceUUID string, userGUID string) *Device
	ViewDeviceByUUIDIncludingSoftDelete(deviceUUID string) *Device
}

// DeviceRepositoryInterface is a contract that defines the method needed for Device Repository.ååå
type DeviceRepositoryInterface interface {
	Create(dbTransaction *gorm.DB, data CreateDevice) (*Device, *systems.ErrorData)
	Update(dbTransaction *gorm.DB, uuid string, data UpdateDevice) *systems.ErrorData
	SetDeletedAtToNull(dbTransaction *gorm.DB, deviceGUID string) *systems.ErrorData
	Delete(dbTransaction *gorm.DB, attribute string, value string) *systems.ErrorData
	GetByUUID(uuid string) *Device
	GetByUUIDAndUserGUID(uuid string, userGUID string) *Device
	GetByUUIDUnscoped(uuid string) *Device
}
