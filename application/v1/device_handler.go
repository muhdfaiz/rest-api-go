package v1

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// DeviceHandler will handle all request related to device endpoint.
type DeviceHandler struct {
	DeviceService DeviceServiceInterface
}

// Create function used to create new device and store in database.
func (dh *DeviceHandler) Create(context *gin.Context) {
	deviceData := CreateDevice{}

	if error := Binding.Bind(&deviceData, context); error != nil {
		context.JSON(http.StatusUnprocessableEntity, error)
		return
	}

	device, error := dh.DeviceService.CreateDevice(deviceData)

	if error != nil {
		errorCode, _ := strconv.Atoi(error.Error.Status)
		context.JSON(errorCode, error)
		return
	}

	context.JSON(http.StatusOK, gin.H{"data": device})
}

// Update function used to update device with new data.
func (dh *DeviceHandler) Update(context *gin.Context) {
	deviceData := UpdateDevice{}

	if err := Binding.Bind(&deviceData, context); err != nil {
		context.JSON(http.StatusUnprocessableEntity, err)
		return
	}

	deviceUUID := context.Param("uuid")

	device, error := dh.DeviceService.UpdateDevice(deviceUUID, deviceData)

	if error != nil {
		errorCode, _ := strconv.Atoi(error.Error.Status)
		context.JSON(errorCode, error)
		return
	}

	context.JSON(http.StatusOK, gin.H{"data": device})
}

// Delete function used to soft delete device by setting current date and time as a value
// for deleted_at column.
func (dh *DeviceHandler) Delete(context *gin.Context) {
	deviceUUID := context.Param("uuid")

	error := dh.DeviceService.DeleteDeviceByUUID(deviceUUID)

	if error != nil {
		errorCode, _ := strconv.Atoi(error.Error.Status)
		context.JSON(errorCode, error)
		return
	}

	result := make(map[string]string)
	result["message"] = "Successfully deleted device with uuid " + deviceUUID

	context.JSON(http.StatusOK, gin.H{"data": result})
}
