package v1

import (
	"net/http"
	"strconv"

	validator "gopkg.in/go-playground/validator.v8"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// DeviceHandler will handle all request related to device endpoint
type DeviceHandler struct{}

// Create function used to create new device and store inside database
func (dh DeviceHandler) Create(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "application/json; charset=utf-8")

	db := c.MustGet("DB").(*gorm.DB)
	tx := db.Begin()

	// Bind request data based on header content type
	deviceData := CreateDevice{}
	if err := c.Bind(&deviceData); err != nil {
		c.JSON(http.StatusBadRequest, Error.ValidationErrors(err.(validator.ValidationErrors)))
		return
	}

	// Retrieve device by UUID
	deviceRepository := &DeviceRepository{DB: tx}
	device := deviceRepository.GetByUUID(deviceData.UUID)

	// If device UUID empty return error message
	if device.UUID != "" {
		c.JSON(http.StatusConflict, Error.DuplicateValueErrors("Device", "uuid", device.UUID))
		return
	}

	// If user GUID exist in the request
	if deviceData.UserGUID != "" {
		// Retrieve user by GUID
		userRepository := &UserRepository{DB: tx}
		user := userRepository.GetByGUID(deviceData.UserGUID)

		// If user GUID empty return error message
		if user.GUID == "" {
			c.JSON(http.StatusBadRequest, Error.ResourceNotFoundError("User", "guid", deviceData.UserGUID))
			return
		}
	}

	// Create new device
	deviceFactory := DeviceFactory{DB: tx}
	result, err := deviceFactory.Create(deviceData)

	// Output error if failed to create new device
	if err != nil {
		tx.Rollback()
		errorCode, _ := strconv.Atoi(err.Error.Status)
		c.JSON(errorCode, err)
		return
	}

	tx.Commit()
	c.JSON(http.StatusOK, gin.H{"data": result})
}

// Update function used to update device with new data.
func (dh DeviceHandler) Update(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "application/json; charset=utf-8")

	db := c.MustGet("DB").(*gorm.DB)
	tx := db.Begin()

	// Retrieve device UUID in url
	deviceUUID := c.Param("uuid")

	// Retrieve device by UUID
	deviceRepository := &DeviceRepository{DB: tx}
	device := deviceRepository.GetByUUID(deviceUUID)

	// If device UUID empty return error message
	if device.UUID == "" {
		c.JSON(http.StatusBadRequest, Error.ResourceNotFoundError("Device", "uuid", deviceUUID))
		return
	}

	// Bind Device data
	deviceData := UpdateDevice{}
	if err := c.Bind(&deviceData); err != nil {
		c.JSON(http.StatusBadRequest, Error.ValidationErrors(err.(validator.ValidationErrors)))
		return
	}

	// If user GUID exist in the request
	if deviceData.UserGUID != "" {
		// Retrieve user by GUID
		userRepository := &UserRepository{DB: tx}
		user := userRepository.GetByGUID(deviceData.UserGUID)

		// If user GUID empty return error message
		if user.GUID == "" {
			c.JSON(http.StatusBadRequest, Error.ResourceNotFoundError("User", "guid", deviceData.UserGUID))
			return
		}
	}

	// Update Device data
	deviceFactory := &DeviceFactory{DB: tx}
	err := deviceFactory.Update(deviceUUID, deviceData)

	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	// Retrieve device latest data
	device = deviceRepository.GetByUUID(deviceUUID)

	tx.Commit()
	c.JSON(http.StatusOK, gin.H{"data": device})
}

// Delete function used to soft delete device by setting current timeo the deleted_at column
func (dh DeviceHandler) Delete(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "application/json; charset=utf-8")

	db := c.MustGet("DB").(*gorm.DB)
	tx := db.Begin()

	// Retrieve device uuid in url
	deviceUUID := c.Param("uuid")

	// Retrieve device by UUID
	deviceRepository := &DeviceRepository{DB: tx}
	device := deviceRepository.GetByUUID(deviceUUID)

	// If device uuid empty return error message
	if device.UUID == "" {
		tx.Rollback()
		c.JSON(http.StatusNotFound, Error.ResourceNotFoundError("Device", "uuid", deviceUUID))
		return
	}

	// Soft delete device
	deviceFactory := &DeviceFactory{DB: tx}
	err := deviceFactory.Delete("uuid", deviceUUID)

	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	tx.Commit()
	// Response data
	result := make(map[string]string)
	result["message"] = "Successfully deleted device with uuid " + device.UUID

	tx.Commit()
	c.JSON(http.StatusOK, gin.H{"data": result})
}
