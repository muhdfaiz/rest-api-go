package v11

import (
	"bitbucket.org/cliqers/shoppermate-api/systems"
	"github.com/jinzhu/gorm"
)

// DeviceService used to handle application logic related to Device resource.
type DeviceService struct {
	DeviceRepository DeviceRepositoryInterface
}

// CheckDuplicateDevice function used to check if the device already exist in database.
func (ds *DeviceService) CheckDuplicateDevice(deviceUUID string) *systems.ErrorData {
	device := ds.DeviceRepository.GetByUUIDUnscoped(deviceUUID)

	if device.UUID != "" {
		return Error.DuplicateValueErrors("Device", "uuid", deviceUUID)
	}

	return nil
}

// CheckDeviceExistOrNot function used to check if the device exist or not in database by checking
// device UUID.
func (ds *DeviceService) CheckDeviceExistOrNot(deviceUUID string) (*Device, *systems.ErrorData) {
	device := ds.DeviceRepository.GetByUUIDUnscoped(deviceUUID)

	if device.UUID == "" {
		return nil, Error.ResourceNotFoundError("Device", "uuid", deviceUUID)
	}

	return device, nil
}

// CreateDevice function used to create new device and store in database.
func (ds *DeviceService) CreateDevice(dbTransaction *gorm.DB, deviceData CreateDevice) (*Device, *systems.ErrorData) {
	error := ds.CheckDuplicateDevice(deviceData.UUID)

	if error != nil {
		return nil, error
	}

	device, error := ds.DeviceRepository.Create(dbTransaction, deviceData)

	if error != nil {
		return nil, error
	}

	return device, nil
}

// UpdateDevice function used to update device information in database.
func (ds *DeviceService) UpdateDevice(dbTransaction *gorm.DB, deviceUUID string, deviceData UpdateDevice) (*Device, *systems.ErrorData) {
	_, error := ds.CheckDeviceExistOrNot(deviceUUID)

	if error != nil {
		return nil, error
	}

	device, error := ds.UpdateByDeviceUUID(dbTransaction, deviceUUID, deviceData)

	if error != nil {
		return nil, error
	}

	return device, nil
}

// UpdateByDeviceUUID function used to update device by Device UUID.
func (ds *DeviceService) UpdateByDeviceUUID(dbTransaction *gorm.DB, deviceUUID string, deviceData UpdateDevice) (*Device, *systems.ErrorData) {
	error := ds.DeviceRepository.Update(dbTransaction, deviceUUID, deviceData)

	if error != nil {
		return nil, error
	}

	device := ds.ViewDeviceByUUID(deviceUUID)

	return device, nil
}

// ReactivateDevice function used to reactivate device by set deleted_at column to NULL.
func (ds *DeviceService) ReactivateDevice(dbTransaction *gorm.DB, deviceGUID string) *systems.ErrorData {
	error := ds.DeviceRepository.SetDeletedAtToNull(dbTransaction, deviceGUID)

	if error != nil {
		return error
	}

	return nil
}

// DeleteDeviceByUUID function used to soft delete device by setting the current date and time
// to deleted_at column.
func (ds *DeviceService) DeleteDeviceByUUID(dbTransaction *gorm.DB, deviceUUID string) *systems.ErrorData {
	_, error := ds.CheckDeviceExistOrNot(deviceUUID)

	if error != nil {
		return error
	}

	error = ds.DeviceRepository.Delete(dbTransaction, "uuid", deviceUUID)

	if error != nil {
		return error
	}

	return nil
}

// ViewDeviceByUUID function used to retrieve device by device UUID.
func (ds *DeviceService) ViewDeviceByUUID(deviceUUID string) *Device {
	device := ds.DeviceRepository.GetByUUID(deviceUUID)

	return device
}

// ViewDeviceByUUIDandUserGUID function used to retrieve device by device UUID and user GUID.
func (ds *DeviceService) ViewDeviceByUUIDandUserGUID(deviceUUID string, userGUID string) *Device {
	device := ds.DeviceRepository.GetByUUIDAndUserGUID(deviceUUID, userGUID)

	return device
}

// ViewDeviceByUUIDIncludingSoftDelete function used to retrieve device by device UUID including
// device that already soft delete.
func (ds *DeviceService) ViewDeviceByUUIDIncludingSoftDelete(deviceUUID string) *Device {
	device := ds.DeviceRepository.GetByUUIDUnscoped(deviceUUID)

	return device
}
