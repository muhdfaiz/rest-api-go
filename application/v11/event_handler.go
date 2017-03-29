package v11

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// EventHandler will handle all request related to Event resource.
type EventHandler struct {
	EventService EventServiceInterface
}

func (eh *EventHandler) ViewAll(c *gin.Context) {
	tokenData := c.MustGet("Token").(map[string]string)

	// Retrieve user guid in url
	userGUID := c.Param("guid")

	// If user GUID not match user GUID inside the token return error message
	if tokenData["user_guid"] != userGUID {
		c.JSON(http.StatusUnauthorized, Error.TokenIdentityNotMatchError("view deals"))
		return
	}

	events := eh.EventService.GetAllIncludingDeals(userGUID)

	c.JSON(http.StatusOK, gin.H{"data": events})
}
