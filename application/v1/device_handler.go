package v1

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// DeviceHandler will handle all request related to device endpoint.
type DeviceHandler struct {
	DeviceService DeviceServiceInterface
	UserService   UserServiceInterface
}

// Create function used to create new device and store in database.
func (dh *DeviceHandler) Create(context *gin.Context) {
	dbTransaction := context.MustGet("DB").(*gorm.DB).Begin()

	deviceData := CreateDevice{}

	if error := Binding.Bind(&deviceData, context); error != nil {
		context.JSON(http.StatusUnprocessableEntity, error)
		return
	}

	if deviceData.UserGUID != "" {
		_, error := dh.UserService.CheckUserGUIDExistOrNot(deviceData.UserGUID)

		if error != nil {
			errorCode, _ := strconv.Atoi(error.Error.Status)
			context.JSON(errorCode, error)
			return
		}
	}

	device, error := dh.DeviceService.CreateDevice(dbTransaction, deviceData)

	if error != nil {
		dbTransaction.Rollback()
		errorCode, _ := strconv.Atoi(error.Error.Status)
		context.JSON(errorCode, error)
		return
	}

	dbTransaction.Commit()

	context.JSON(http.StatusOK, gin.H{"data": device})
}

// Update function used to update device with new data.
func (dh *DeviceHandler) Update(context *gin.Context) {
	dbTransaction := context.MustGet("DB").(*gorm.DB).Begin()

	deviceData := UpdateDevice{}

	if error := Binding.Bind(&deviceData, context); error != nil {
		context.JSON(http.StatusUnprocessableEntity, error)
		return
	}

	if deviceData.UserGUID != "" {
		_, error := dh.UserService.CheckUserGUIDExistOrNot(deviceData.UserGUID)

		if error != nil {
			errorCode, _ := strconv.Atoi(error.Error.Status)
			context.JSON(errorCode, error)
			return
		}
	}

	deviceUUID := context.Param("uuid")

	device, error := dh.DeviceService.UpdateDevice(dbTransaction, deviceUUID, deviceData)

	if error != nil {
		dbTransaction.Rollback()
		errorCode, _ := strconv.Atoi(error.Error.Status)
		context.JSON(errorCode, error)
		return
	}

	dbTransaction.Commit()

	device = dh.DeviceService.ViewDeviceByUUID(deviceUUID)

	context.JSON(http.StatusOK, gin.H{"data": device})
}

// Delete function used to soft delete device by setting current date and time as a value
// for deleted_at column.
func (dh *DeviceHandler) Delete(context *gin.Context) {
	dbTransaction := context.MustGet("DB").(*gorm.DB).Begin()

	deviceUUID := context.Param("uuid")

	error := dh.DeviceService.DeleteDeviceByUUID(dbTransaction, deviceUUID)

	if error != nil {
		dbTransaction.Rollback()
		errorCode, _ := strconv.Atoi(error.Error.Status)
		context.JSON(errorCode, error)
		return
	}

	dbTransaction.Commit()
	result := make(map[string]string)
	result["message"] = "Successfully deleted device with uuid " + deviceUUID

	context.JSON(http.StatusOK, gin.H{"data": result})
}
