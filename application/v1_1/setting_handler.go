package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// SettingHandler will handle all request related to setting resources.
type SettingHandler struct {
	SettingService SettingServiceInterface
}

// ViewAll function used to retrieve all settings from database through Setting Service.
func (sh *SettingHandler) ViewAll(context *gin.Context) {
	settings := sh.SettingService.GetAllSettings()

	context.JSON(http.StatusOK, gin.H{"data": settings})
}
