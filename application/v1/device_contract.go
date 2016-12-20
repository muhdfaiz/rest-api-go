package v1

import "bitbucket.org/cliqers/shoppermate-api/systems"

// DeviceServiceInterface is a contract that defines the method needed for Device Service.
type DeviceServiceInterface interface {
	CheckDuplicateDevice(deviceUUID string) *systems.ErrorData
	CreateDevice(deviceData CreateDevice) (*Device, *systems.ErrorData)
	UpdateDevice(deviceUUID string, deviceData UpdateDevice) (*Device, *systems.ErrorData)
	UpdateByDeviceUUID(deviceUUID string, deviceData UpdateDevice) (*Device, *systems.ErrorData)
	ReactivateDevice(deviceGUID string) *systems.ErrorData
	DeleteDeviceByUUID(deviceUUID string) *systems.ErrorData
	ViewDeviceByUUID(deviceUUID string) *Device
	ViewDeviceByUUIDandUserGUID(deviceUUID string, userGUID string) *Device
	ViewDeviceByUUIDIncludingSoftDelete(deviceUUID string) *Device
}

// DeviceRepositoryInterface is a contract that defines the method needed for Device Repository.ååå
type DeviceRepositoryInterface interface {
	Create(data CreateDevice) (*Device, *systems.ErrorData)
	Update(uuid string, data UpdateDevice) *systems.ErrorData
	SetDeletedAtToNull(deviceGUID string) *systems.ErrorData
	Delete(attribute string, value string) *systems.ErrorData
	GetByUUID(uuid string) *Device
	GetByUUIDAndUserGUID(uuid string, userGUID string) *Device
	GetByUUIDUnscoped(uuid string) *Device
}
