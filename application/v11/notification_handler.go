package v11

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// NotificationHandler will handle all request related to Notification resource.
type NotificationHandler struct {
	NotificationService NotificationServiceInterface
	DeviceService       DeviceServiceInterface
}

func (nh *NotificationHandler) ViewNotificationForGuest(context *gin.Context) {
	deviceUUID := context.Param("device_uuid")

	_, error := nh.DeviceService.CheckDeviceExistOrNot(deviceUUID)

	if error != nil {
		errorCode, _ := strconv.Atoi(error.Error.Status)
		context.JSON(errorCode, error)
		return
	}

	notifications := nh.NotificationService.GetNotificationsForGuest(deviceUUID)

	context.JSON(http.StatusOK, gin.H{"data": notifications})
}

func (nh *NotificationHandler) ViewNotificationForRegisteredUser(context *gin.Context) {
	deviceUUID := context.Param("device_uuid")

	userGUID := context.Param("user_guid")

	userToken := context.MustGet("Token").(map[string]string)

	if userToken["device_uuid"] != deviceUUID || userToken["user_guid"] != userGUID {
		context.JSON(http.StatusUnauthorized, Error.TokenIdentityNotMatchError("View Notifications"))
		return
	}

	notifications := nh.NotificationService.GetNotificationsForLoggedInUser(deviceUUID, userGUID)

	context.JSON(http.StatusOK, gin.H{"data": notifications})
}
