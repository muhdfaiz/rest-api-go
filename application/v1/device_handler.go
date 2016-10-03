package v1

import (
	"net/http"
	"strconv"

	validator "gopkg.in/go-playground/validator.v8"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// DeviceHandler will handle all request related to device endpoint
type DeviceHandler struct {
	UserRepository   UserRepositoryInterface
	DeviceRepository DeviceRepositoryInterface
	DeviceFactory    DeviceFactoryInterface
}

// Create function used to create new device and store inside database
func (dh *DeviceHandler) Create(c *gin.Context) {
	DB := c.MustGet("DB").(*gorm.DB).Begin()

	// Bind request data based on header content type
	deviceData := CreateDevice{}
	if err := c.Bind(&deviceData); err != nil {
		DB.Rollback().Close()
		c.JSON(http.StatusBadRequest, Error.ValidationErrors(err.(validator.ValidationErrors)))
		return
	}

	// Retrieve device by UUID
	device := dh.DeviceRepository.GetByUUID(deviceData.UUID)

	// If device UUID not empty return error message
	if device.UUID != "" {
		DB.Rollback().Close()
		c.JSON(http.StatusConflict, Error.DuplicateValueErrors("Device", "uuid", device.UUID))
		return
	}

	// If user GUID exist in the request
	if deviceData.UserGUID != "" {
		// Retrieve user by GUID
		user := dh.UserRepository.GetByGUID(deviceData.UserGUID)

		// If user GUID empty return error message
		if user.GUID == "" {
			DB.Rollback().Close()
			c.JSON(http.StatusBadRequest, Error.ResourceNotFoundError("User", "guid", deviceData.UserGUID))
			return
		}
	}

	// Create new device
	result, err := dh.DeviceFactory.Create(deviceData)

	// Output error if failed to create new device
	if err != nil {
		DB.Rollback().Close()
		errorCode, _ := strconv.Atoi(err.Error.Status)
		c.JSON(errorCode, err)
		return
	}

	DB.Commit().Close()
	c.JSON(http.StatusOK, gin.H{"data": result})
}

// Update function used to update device with new data.
func (dh *DeviceHandler) Update(c *gin.Context) {
	DB := c.MustGet("DB").(*gorm.DB).Begin()

	// Retrieve device UUID in url
	deviceUUID := c.Param("uuid")

	// Retrieve device by UUID
	device := dh.DeviceRepository.GetByUUID(deviceUUID)

	// If device UUID empty return error message
	if device.UUID == "" {
		DB.Rollback().Close()
		c.JSON(http.StatusNotFound, Error.ResourceNotFoundError("Device", "uuid", deviceUUID))
		return
	}

	// Bind Device data
	deviceData := UpdateDevice{}

	// Bind request based on content type and validate request data
	if err := Binding.Bind(&deviceData, c); err != nil {
		DB.Rollback().Close()
		c.JSON(http.StatusUnprocessableEntity, err)
		return
	}

	// If user GUID exist in the request
	if deviceData.UserGUID != "" {
		// Retrieve user by GUID
		user := dh.UserRepository.GetByGUID(deviceData.UserGUID)

		// If user GUID empty return error message
		if user.GUID == "" {
			DB.Rollback().Close()
			c.JSON(http.StatusNotFound, Error.ResourceNotFoundError("User", "guid", deviceData.UserGUID))
			return
		}
	}

	// Update Device data
	err := dh.DeviceFactory.Update(deviceUUID, deviceData)

	if err != nil {
		DB.Rollback().Close()
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	// Retrieve device latest data
	device = dh.DeviceRepository.GetByUUID(deviceUUID)

	DB.Commit().Close()
	c.JSON(http.StatusOK, gin.H{"data": device})
}

// Delete function used to soft delete device by setting current timeo the deleted_at column
func (dh *DeviceHandler) Delete(c *gin.Context) {
	DB := c.MustGet("DB").(*gorm.DB).Begin()

	// Retrieve device uuid in url
	deviceUUID := c.Param("uuid")

	// Retrieve device by UUID
	device := dh.DeviceRepository.GetByUUID(deviceUUID)

	// If device uuid empty return error message
	if device.UUID == "" {
		DB.Rollback().Close()
		c.JSON(http.StatusNotFound, Error.ResourceNotFoundError("Device", "uuid", deviceUUID))
		return
	}

	// Soft delete device
	err := dh.DeviceFactory.Delete("uuid", deviceUUID)

	if err != nil {
		DB.Rollback().Close()
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	// Response data
	result := make(map[string]string)
	result["message"] = "Successfully deleted device with uuid " + device.UUID

	DB.Commit().Close()
	c.JSON(http.StatusOK, gin.H{"data": result})
}
