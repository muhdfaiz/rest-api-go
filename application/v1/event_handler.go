package v1

import "github.com/gin-gonic/gin"

type EventHandler struct {
	EventService EventServiceInterface
}

func (fdh *EventHandler) ViewAll(c *gin.Context) {

}
