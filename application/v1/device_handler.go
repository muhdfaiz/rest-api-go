package v1

import (
	"fmt"
	"net/http"
	"strconv"

	"bitbucket.org/shoppermate-api/systems"

	validator "gopkg.in/go-playground/validator.v8"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// DeviceHandler will handle all request related to device endpoint
type DeviceHandler struct{}

// Create function used to create new device and store inside database
func (dh DeviceHandler) Create(c *gin.Context) {
	db := c.MustGet("DB").(*gorm.DB)
	tx := db.Begin()

	c.Writer.Header().Set("Content-Type", "application/json; charset=utf-8")

	// Bind request data based on header content type
	deviceData := CreateDevice{}
	if err := c.Bind(&deviceData); err != nil {
		c.JSON(http.StatusBadRequest, ErrorMesg.ValidationErrors(err.(validator.ValidationErrors)))
		return
	}

	// Check Device already exist or not
	deviceRepository := &DeviceRepository{DB: tx}
	device := deviceRepository.GetByUUID(deviceData.UUID)
	if device.UUID != "" {
		c.JSON(http.StatusConflict, ErrorMesg.DuplicateValueErrors("Device", "uuid", device.UUID))
		return
	}

	// Check User GUID valid or not.
	// Return error if not valid
	if deviceData.UserGUID != "" {
		userRepository := &UserRepository{DB: tx}
		user := userRepository.GetByGUID(deviceData.UserGUID)

		if user.GUID == "" {
			c.JSON(http.StatusBadRequest, ErrorMesg.GenericError(strconv.Itoa(http.StatusBadRequest), systems.ResourceNotFound,
				fmt.Sprintf(systems.TitleResourceNotFoundError, "User"), "message",
				fmt.Sprintf(systems.ErrorResourceNotFound, "User", "guid", deviceData.UserGUID)))
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
	db := c.MustGet("DB").(*gorm.DB)
	tx := db.Begin()

	c.Writer.Header().Set("Content-Type", "application/json; charset=utf-8")

	// Retrieve device uuid in url
	deviceUUID := c.Param("uuid")

	// Check device uuid exist or not
	// If not exist display error message
	deviceRepository := &DeviceRepository{DB: tx}
	device := deviceRepository.GetByUUID(deviceUUID)

	// Return error message if device uuid not exist
	if device.UUID == "" {
		c.JSON(http.StatusBadRequest, ErrorMesg.GenericError(strconv.Itoa(http.StatusBadRequest), systems.ResourceNotFound,
			fmt.Sprintf(systems.TitleResourceNotFoundError, "Device"), "message",
			fmt.Sprintf(systems.ErrorResourceNotFound, "Device", "uuid", deviceUUID)))
		return
	}

	// Bind Device data
	deviceData := UpdateDevice{}
	if err := c.Bind(&deviceData); err != nil {
		c.JSON(http.StatusBadRequest, ErrorMesg.ValidationErrors(err.(validator.ValidationErrors)))
		return
	}

	// Check User GUID valid or not.
	// Return error if not valid
	if deviceData.UserGUID != "" {
		userRepository := &UserRepository{DB: tx}
		user := userRepository.GetByGUID(deviceData.UserGUID)

		if user.GUID == "" {
			c.JSON(http.StatusBadRequest, ErrorMesg.GenericError(strconv.Itoa(http.StatusBadRequest), systems.ResourceNotFound,
				fmt.Sprintf(systems.TitleResourceNotFoundError, "User"), "message",
				fmt.Sprintf(systems.ErrorResourceNotFound, "User", "guid", deviceData.UserGUID)))
			return
		}
	}

	// Update Device
	deviceFactory := &DeviceFactory{DB: tx}
	err := deviceFactory.Update(deviceUUID, deviceData)

	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	// Retrieve device latest update
	device = deviceRepository.GetByUUID(deviceUUID)
	tx.Commit()
	c.JSON(http.StatusOK, gin.H{"data": device})
}

func (dh DeviceHandler) Delete(c *gin.Context) {
	db := c.MustGet("DB").(*gorm.DB)
	tx := db.Begin()

	c.Writer.Header().Set("Content-Type", "application/json; charset=utf-8")

	// Retrieve device uuid in url
	deviceUUID := c.Param("uuid")

	// Check device uuid exist or not
	// If not exist display error message
	deviceRepository := &DeviceRepository{DB: tx}
	device := deviceRepository.GetByUUID(deviceUUID)

	// Return error message if device uuid not exist
	if device.UUID == "" {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, ErrorMesg.GenericError(strconv.Itoa(http.StatusBadRequest), systems.ResourceNotFound,
			fmt.Sprintf(systems.TitleResourceNotFoundError, "Device"), "message",
			fmt.Sprintf(systems.ErrorResourceNotFound, "Device", "uuid", deviceUUID)))
		return
	}

	deviceFactory := &DeviceFactory{DB: tx}
	err := deviceFactory.Delete("uuid", deviceUUID)

	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	tx.Commit()
	c.JSON(http.StatusNoContent, gin.H{})
}
